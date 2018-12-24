package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

func run() (err error) {
	fmt.Println(kafkaClient)

	partitionList, err := kafkaClient.client.Partitions(kafkaClient.topic)
	if err != nil {
		logs.Error("Failed to get the list of partitions: ", err)
		return
	}

	for partition := range partitionList {
		pc, errRet := kafkaClient.client.ConsumePartition(kafkaClient.topic, int32(partition), sarama.OffsetNewest)
		if errRet != nil {
			err = errRet
			logs.Error("Failed to start consumer for partition %d: %s\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		kafkaClient.wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				logs.Debug("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				fmt.Println("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				//fmt.Println()
				err = sendToES(kafkaClient.topic, msg.Value)
				if err != nil {
					logs.Warn("send to es failed, err:%v", err)
				}
			}
			fmt.Println(1)
			kafkaClient.wg.Done()
		}(pc)
	}

	//time.Sleep(time.Hour)
	kafkaClient.wg.Wait()
	return
}
