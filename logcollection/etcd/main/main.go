package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"golang-awesome/logcollection/logagent/tail"
	"time"
)

var (
	etcdKey = "/net/badme/logagent/config/182.61.137.53"
)

func main() {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}

	fmt.Println("connect success")
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	logconf := &tail.CollectConf{
		LogPath: "/Users/qinxun/dev/dev-go/go/test.log",
		Topic:   "test",
	}
	var confs []*tail.CollectConf
	confs = append(confs, logconf)
	bytes, err := json.Marshal(confs)
	_, err = cli.Put(ctx, etcdKey, string(bytes))
	cancel()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, etcdKey)
	cancel()
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		//var conf tail.CollectConf
		//err := json.Unmarshal(ev.Value, &conf)
		//if err != nil{
		//	panic(err)
		//}
		//fmt.Println("conf:",conf)
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)

	}
}
