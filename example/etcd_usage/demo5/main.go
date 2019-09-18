package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

func main() {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		err     error
		kv      clientv3.KV
		delResp *clientv3.DeleteResponse
		kvPair  *mvccpb.KeyValue
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

	// 用于读写etcd的键值对
	kv = clientv3.NewKV(client)

	// 删除KV,删除以"/cron/jobs/job打头的key",并返回删除的kv信息
	if delResp, err = kv.Delete(context.TODO(), "/cron/jobs/job", clientv3.WithPrevKV(), clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
		return
	}

	// 被删除之前的value是什么
	if len(delResp.PrevKvs) != 0 {
		for _, kvPair = range delResp.PrevKvs {
			fmt.Println("删除了:", string(kvPair.Key), string(kvPair.Value))
		}
	}
}
