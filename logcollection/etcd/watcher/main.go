package main

import (
	"context"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"golang-awesome/logcollection/logagent/tail"
	"time"
)

const (
	EtcdKey = "/logagent/conf/"
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
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s] 's config deleted", key)
					continue
				}

				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					//err = json.Unmarshal(ev.Kv.Value, &collectConf)
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
	//wg.Done()
}

func main(){
	initEtcdWatch()
}
