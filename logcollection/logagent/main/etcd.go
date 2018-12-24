package main

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"golang-awesome/logcollection/logagent/tail"
	"time"
)

type EtcdClient struct {
	client *clientv3.Client
	keys   []string
}

var (
	etcdClient *EtcdClient
)

func initEtcd(addr,key string)(collectConf tail.CollectConf,err error){
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		logs.Error("connect etcd failed, err:",err)
		
	}

	etcdClient =&EtcdClient{
		client:client,
	}

	//if strings.HasSuffix(key,"/") == false {
	//	key = key + "/"
	//}

	//for _,ip := range localIPArray{
		//etcdKey := fmt.Sprintf("%s%s", key, ip)
		//etcdKey := fmt.Sprintf("%s", key)
		etcdClient.keys = append(etcdClient.keys, key)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := client.Get(ctx, key)
		if err != nil {
			logs.Error("client get from etcd failed, err:%v", err)
		}
		cancel()
		logs.Debug("resps from etcd:%v", resp.Kvs)
		for _, v := range resp.Kvs {
			if string(v.Key) == key {
				logs.Debug("v.Key:",v.Key)
				logs.Debug("v.Value:",v.Value)
				err = json.Unmarshal(v.Value, &collectConf)
				if err != nil {
					logs.Error("unmarshal failed, err:%v", err)
					continue
				}

				logs.Debug("log config is %v", collectConf)
			}
		}
	//}

	initEtcdWatcher()
	return
}

func initEtcdWatcher() {

	for _, key := range etcdClient.keys {
		logs.Debug("etcd key :" ,key)
		go watchKey(key)
	}
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
			for _, ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s] 's config deleted", key)
					continue
				}

				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err = json.Unmarshal(ev.Kv.Value, &collectConf)
					if err != nil {
						logs.Error("key [%s], Unmarshal[%s], err:%v ", err)
						getConfSucc = false
						continue
					}
				}
				logs.Debug("get config from etcd, %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}

			if getConfSucc {
				logs.Debug("get config from etcd succ, %v", collectConf)
				tail.UpdateConfig(collectConf)
			}
		}

	}
}
