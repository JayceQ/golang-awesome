package main

import (
	"github.com/astaxie/beego/logs"
	"golang-awesome/logcollection/logagent/kafka"
	"golang-awesome/logcollection/logagent/tail"
	"time"
)

func serverRun() (err error) {
	for {
		msg := tail.GetOneLine()
		err = sendToKafka(msg)
		if err != nil {
			logs.Error("send to kafka failed,err:%v", err)
			time.Sleep(time.Second)
			continue
		}
	}
	return
}

func sendToKafka(msg *tail.TextMsg) (err error) {
	err = kafka.SendToKafka(msg.Msg, msg.Topic)
	return
}
