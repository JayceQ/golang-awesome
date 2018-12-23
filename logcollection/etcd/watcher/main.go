package main

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"golang-awesome/logcollection/logagent/tail"
	"time"
)

const (
	EtcdKey = "mykey"
)

func initEtcdWatch(){
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	watchKey(EtcdKey)
	//wg.Wait()
}
func watchKey(key string) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err:", err)
		return
	}

	logs.Debug("begin watch key:%s", key)
	for {
		rch := cli.Watch(context.Background(), key)
		var collectConf []tail.CollectConf
		var getConfSucc = true
		for wresp := range rch {
			logs.Debug("watch -->",wresp)
			for _, ev := range wresp.Events {
				fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}

			if getConfSucc {
				logs.Debug("get config from etcd succ, %v", collectConf)
				tail.UpdateConfig(collectConf)
			}
		}

	}
	//wg.Done()
}

func main(){
	//initEtcdWatch()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}

	fmt.Println("connect succ")
	defer cli.Close()

	response, err := cli.Put(context.Background(), "/logagent/conf/", "8888888")
	fmt.Println(response)
	for {
		rch := cli.Watch(context.Background(), "/logagent/conf/")
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}
}
