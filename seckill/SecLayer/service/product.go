package service

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

func loadProductFromEtcd(conf *SecLayerConf)(err error){

	logs.Debug("start get product from etcd")
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	resp,err := secLayerContext.etcdClient.Get(ctx,conf.EtcdConfig.EtcdSecProductKey)
	if err!= nil{
		logs.Debug("get [%s] from etcd failed, err: %v",conf.EtcdConfig.EtcdSecProductKey,err)
		return
	}

	logs.Debug("get from etcd success, resp: %v",resp)
	var secProductInfo []SecProductInfoConf
	for k,v := range resp.Kvs{
		logs.Debug("key[%v] value[%v]",k,v)
		err = json.Unmarshal(v.Value,&secProductInfo)
		if err!= nil {
			logs.Error("Unmarshal sec product info failed, err: %v",err)
			return
		}

		logs.Debug("sec info conf is [%v]",secProductInfo)
	}

	updateSecProductInfo(conf,secProductInfo)
	logs.Debug("update product info success, data:%v",err)

	initSecProductWatcher(conf)
	logs.Debug("init etcd watcher success")
	return
}


func updateSecProductInfo(conf *SecLayerConf,secProductIndo []SecProductInfoConf){
	tmp := make(map[int]*SecProductInfoConf,1024)
	for _,v :=range secProductIndo{
		productInfo := v
		productInfo.secLimit = &SecLimit{}
		tmp[v.ProductId] = &productInfo
	}

	secLayerContext.RWSecProductLock.Lock()
	conf.SecProductInfoMap = tmp
	secLayerContext.RWSecProductLock.Unlock()
}

func initSecProductWatcher(conf *SecLayerConf){
	go watchSecProductKey(conf)
}

func watchSecProductKey(conf *SecLayerConf){

	key := conf.EtcdConfig.EtcdSecProductKey
	logs.Debug("begin watch key: %s",key)
	var err error
	for{
		rch := secLayerContext.etcdClient.Watch(context.Background(),key)
		var secProductInfo []SecProductInfoConf
		var getConfSucc = true

		for wresp := range rch{
			for _,ev := range wresp.Events{
				if ev.Type == mvccpb.DELETE{
					logs.Warn("key[%s] 's config deleted ",key)
					continue
				}

				if ev.Type ==mvccpb.PUT && string(ev.Kv.Key) == key{
					err = json.Unmarshal(ev.Kv.Value,&secProductInfo)
					if err !=nil{
						logs.Error("key [%s] ,Unmarshal [%s] ,err: %v",err)
						getConfSucc = false
						continue
					}
				}

				logs.Debug("get config from etcd, %s %q: %q\n",ev.Type,ev.Kv.Key,ev.Kv.Value)
			}

			if getConfSucc {
				logs.Debug("get config from etcd success, %v",secProductInfo)
				updateSecProductInfo(conf,secProductInfo)
			}
		}
	}
}