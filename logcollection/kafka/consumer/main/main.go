package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"strings"
	"time"
)

func main(){

	consumer, e := sarama.NewConsumer(strings.Split("127.0.0.1:9092", ","), nil)
	if e != nil{
		fmt.Println("failed to start consumer: ",e)
		return
	}
	partitionList, e := consumer.Partitions("test")
	if e != nil {
		fmt.Println("failed to get the list of partitions: ",e)
		return
	}
	fmt.Println(partitionList)
	for partition := range partitionList{
		pc, e := consumer.ConsumePartition("test", int32(partition), sarama.OffsetNewest)
		if e != nil {
			fmt.Printf("failed to start consumer for partition %d,%s\n",partition,e)
			return
		}
		defer pc.AsyncClose()

		go func(sarama.PartitionConsumer){
			for msg := range pc.Messages(){
				fmt.Printf("partition:%d, offset:%d, key:%s, value:%s",msg.Partition,msg.Offset,string(msg.Key),string(msg.Value))
				fmt.Println()
			}
		}(pc)

		time.Sleep(time.Hour)
		consumer.Close()
	}
}