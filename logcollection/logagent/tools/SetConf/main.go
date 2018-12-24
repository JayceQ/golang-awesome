package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

const (
	EtcdKey = "/net/badme/logagent/config/182.61.137.53"
)

func SetLogConfToEtcd() {
	client, e := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if e != nil {
		fmt.Println("connect etcd failed, err:", e)
	}

	fmt.Println("connect success")
	defer client.Close()

	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	response, e := client.Get(timeout, EtcdKey)
	cancelFunc()
	if e != nil {
		fmt.Println("get failed, err:", e)
		return
	}

	for _, ev := range response.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}

func main() {
	SetLogConfToEtcd()
}
