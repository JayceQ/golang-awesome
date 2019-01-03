package service

import (
	"github.com/astaxie/beego/logs"
	etcd "go.etcd.io/etcd/clientv3"
	"time"
)

func initEtcd(conf *SecLayerConf)(err error){
	cli,err :=etcd.New(etcd.Config{
		Endpoints:[]string{conf.EtcdConfig.EtcdAddr},
		DialTimeout:time.Duration(conf.EtcdConfig.Timeout) * time.Second,
	})

	if err != nil{
		logs.Error("connect etcd failed, err: ",err)
		return
	}

	secLayerContext.etcdClient = cli
	logs.Debug("init etcd success")
	return
}

func InitSecLayer(conf *SecLayerConf)(err error){

	err = initRedis(conf)
	if err!= nil{
		logs.Error("init redis failed, err: %v",err)
		return
	}
	logs.Debug("init redis success")

	err = initEtcd(conf)
	if err != nil {
		logs.Error("init etcd failed, err: %v",err)
		return
	}
	logs.Debug("init etcd success")

	err = loadProductFromEtcd(conf)
	if err!= nil {
		logs.Error("load product from etcd failed, err: %v",err)
		return
	}
	logs.Debug("load product info success")

	secLayerContext.secLayerConf = conf
	secLayerContext.Read2HandleChan = make(chan *SecRequest,
		secLayerContext.secLayerConf.Read2HandleChanSize)
	secLayerContext.Handle2WriteChan = make(chan *SecResponse,
		secLayerContext.secLayerConf.Handle2WriteChanSize)

	secLayerContext.HistoryMap = make(map[int]*UserBuyHistory,1000000)

	secLayerContext.productCountMgr = NewProductCountMgr()

	logs.Debug("init all success")
	return
}