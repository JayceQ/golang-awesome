package main

import (
	"github.com/Shopify/sarama"
	"fmt"
)

func main(){

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer([]string{"182.61.137.53:9092"}, config)
	if err != nil{
		fmt.Println("producer close,err:",err)
		return
	}

	defer client.Close()


	msg := &sarama.ProducerMessage{}
	msg.Topic = "test"
	msg.Value = sarama.StringEncoder("this is a test massage, my massage is good")

	partition, offset, e := client.SendMessage(msg)
	if e != nil {
		fmt.Println("send massage failed," , err)
		return
	}

	fmt.Printf("partition :%v ,offset:%v" ,partition,offset)


}