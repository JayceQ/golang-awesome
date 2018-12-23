package main

import (
	"fmt"

	"github.com/astaxie/beego/logs"
)

func main() {

	err := initConfig("ini", "./log_transfer.conf")
	if err != nil {
		panic(err)
		return
	}
	fmt.Println(logConfig)

	err = initLogger(logConfig.LogPath, logConfig.LogLevel)
	if err != nil {
		panic(err)
		return
	}

	logs.Debug("init logger succ")

	err = initKafka(logConfig.KafkaAddr, logConfig.KafkaTopic)
	if err != nil {
		logs.Error("init kafka failed, err:%v", err)
		return
	}

	logs.Debug("init kafka succ")

	err = initES(logConfig.ESAddr)
	if err != nil {
		logs.Error("init es failed, err:%v", err)
		return
	}

	logs.Debug("init es client succ")

	err = run()
	if err != nil {
		logs.Error("run  failed, err:%v", err)
		return
	}

	logs.Warn("warning, log_transfer is exited")

}

