package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

var (
	client sarama.SyncProducer
)

func InitKafka(addr string)(e error){

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, e = sarama.NewSyncProducer([]string{addr}, config)
	if e != nil {
		logs.Error("init kafka producer failed, err :",e)
		return
	}

	logs.Debug("init kafka success")
	return
}

func SendToKafka(data,topic string)(err error){
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)

	partition, offset, err := client.SendMessage(msg)
	if err != nil {
		logs.Error("send message failed, err :%v data:%v topic :%v",err,data,topic)
		return
	}
	logs.Debug("send success, partition:%v offset:%v topic:%v",partition,offset,topic)
	return
}