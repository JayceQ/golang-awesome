package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"golang-awesome/logcollection/logagent/kafka"
	"golang-awesome/logcollection/logagent/tail"
)

func main(){
	filename := "./logagent.conf"
	err := loadConf("ini",filename)
	if err != nil {
		fmt.Printf("load conf failed,err%v\n",err)
		panic("load conf failed")
		return
	}

	err = initLogger()
	if err!= nil{
		fmt.Printf("load conf failed,err%v\n",err)
		panic("load log failed")
		return
	}

	logs.Debug("load conf succ ,config:",appConfig)

	collectConf, err := initEtcd(appConfig.etcdAddr, appConfig.etcdKey)
	if err != nil {
		logs.Error("init etcd failed, err:%v", err)
		return
	}

	err = tail.InitTail(collectConf,appConfig.chanSize)
	if err != nil {
		logs.Error("init tail failed ,err:%v\n",err)
		return
	}

	logs.Debug("initialize tail success!")

	err = kafka.InitKafka(appConfig.kafkaAddr)
	if err != nil {
		logs.Error("init kafka failed ,err:%v\n",err)
		return
	}

	logs.Debug("initialize kafka success!")


	logs.Debug("initialize all config success!")

	err = serverRun()
	if err != nil {
		logs.Error("serverRun failed ,err :%v\n",err)
		return
	}
	logs.Info("program exited")


}
