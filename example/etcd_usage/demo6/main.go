package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	var (
		config         clientv3.Config
		client         *clientv3.Client
		err            error
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		putResp        *clientv3.PutResponse
		getResp        *clientv3.GetResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse // 只读channel
		kv             clientv3.KV
	)

	config = clientv3.Config{
		Endpoints:   []string{"172.17.10.210:2379"}, // 集群列表
		DialTimeout: 5 * time.Second,
	}

	// 建立一个客户端
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	// 申请一个lease(租约)
	lease = clientv3.NewLease(client)

	// 申请一个10秒的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}

	// 拿到租约的ID
	leaseId = leaseGrantResp.ID

	// 自动续租(保持租约不过期),每隔多长时间自动续租(好像是3秒或者4秒)
	if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
		fmt.Println(err)
		return
	}

	// 处理续约应答的协程
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				//if keepRespChan == nil { // 不能用keepRespChan判断
				if keepResp == nil {
					fmt.Println("租约已经失效了") // 可能是由于网络或者机器down机，或者程序用了cancelFunc取消函数
					goto END
				} else { // 每3-4秒会续租一次(旧的etcd 版本是每隔1秒续租一次)，所以就会收到一次应答
					fmt.Println("收到自动续租应答：", keepResp.ID)
				}
			}
		}
	END:
	}()

	// 获得kv API子集
	kv = clientv3.NewKV(client)

	// Put一个KV,让它与租约关联起来，从而实现10秒后自动过期
	if putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "job1", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("写入成功：", putResp.Header.Revision)

	// 定时的看一下key过期没有
	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("kv过期了")
			break
		}
		fmt.Println("还没过期：", getResp.Kvs)
		time.Sleep(2 * time.Second) // 每隔2秒检查一次kv是否过期
	}
}
