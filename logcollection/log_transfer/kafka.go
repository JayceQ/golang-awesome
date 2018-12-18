package main

import (
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

/*
var (
	wg sync.WaitGroup
)
*/

type KafkaClient struct {
	client sarama.Consumer
	addr   string
	topic  string
	wg     sync.WaitGroup
}

var (
	kafkaClient *KafkaClient
)

func initKafka(addr string, topic string) (err error) {

	kafkaClient = &KafkaClient{}

	consumer, err := sarama.NewConsumer(strings.Split(addr, ","), nil)
	if err != nil {
		logs.Error("init kafka failed, err:%v", err)
		return
	}

	kafkaClient.client = consumer
	kafkaClient.addr = addr
	kafkaClient.topic = topic
	return
	/*
		partitionList, err := consumer.Partitions(topic)
		if err != nil {
			logs.Error("Failed to get the list of partitions: ", err)
			return
		}

		for partition := range partitionList {
			pc, errRet := consumer.ConsumePartition("nginx_log", int32(partition), sarama.OffsetNewest)
			if errRet != nil {
				err = errRet
				logs.Error("Failed to start consumer for partition %d: %s\n", partition, err)
				return
			}
			defer pc.AsyncClose()
			go func(pc sarama.PartitionConsumer) {
				//wg.Add(1)
				for msg := range pc.Messages() {
					logs.Debug("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
					//fmt.Println()
				}
				//wg.Done()
			}(pc)
		}*/
	//time.Sleep(time.Hour)
	//wg.Wait()
	//consumer.Close()
	//return
}
