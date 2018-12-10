package main

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/config"
	"golang-awesome/logcollection/logagent/tail"
)

var (
	appConfig *Config
)

type Config struct {
	logLevel    string
	logPath     string
	chanSize int
	kafkaAddr string
	collectConf []tail.CollectConf
}


func loadCollectConf(conf config.Configer) (err error) {
	var cc tail.CollectConf
	cc.LogPath = conf.String("collect::log_path")
	if len(cc.LogPath) == 0 {
		err = errors.New("invalid collect::log_path")
		return
	}

	cc.Topic = conf.String("collect::topic")
	if len(cc.Topic) == 0 {
		err = errors.New("invalid collect::topic")
		return
	}

	appConfig.collectConf = append(appConfig.collectConf, cc)
	return
}

func loadConf(confType string, filename string) (err error) {
	conf, err := config.NewConfig(confType, filename)
	if err != nil {
		fmt.Println("new config failed ,err :", err)
		return
	}

	appConfig = &Config{}
	appConfig.logLevel = conf.String("logs::log_level")
	if len(appConfig.logLevel) == 0 {
		appConfig.logLevel = "debug"
	}

	appConfig.logPath = conf.String("logs::log_path")
	if len(appConfig.logPath) == 0 {
		appConfig.logPath = "./logs"
	}

	appConfig.chanSize, err = conf.Int("collect::chan_size")
	if err != nil {
		appConfig.chanSize = 100
	}

	appConfig.kafkaAddr= conf.String("kafka::server_addr")
	if len(appConfig.kafkaAddr) == 0 {
		err = fmt.Errorf("invalid kafka address")
		return
	}

	err = loadCollectConf(conf)
	if err != nil {
		fmt.Printf("load collect conf failed, err:%v\n", err)
		return
	}
	return
}
