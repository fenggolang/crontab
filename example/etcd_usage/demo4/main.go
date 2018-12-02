package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		err     error
		kv      clientv3.KV
		getResp *clientv3.GetResponse
	)

	config = clientv3.Config{
		Endpoints:   []string{"172.17.10.210:2379"},
		DialTimeout: 5 * time.Second,
	}

	// 建立一个客户端
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	// 用于读写etcd的键值对
	kv = clientv3.NewKV(client)

	// 协议另外一个job
	kv.Put(context.TODO(), "/cron/jobs/job2", "{...}")

	// 读取/cron/jobs为前缀的所有key
	//if getResp, err = kv.Get(context.TODO(), "/cron/jobs/", clientv3.WithPrefix(), clientv3.GetLogger()); err != nil {
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/", clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
		return
	} else { // 获取成功，遍历所有的kvc
		fmt.Println(getResp.Header.RaftTerm)
		fmt.Println(getResp.Header.GetRaftTerm())
	}
}
