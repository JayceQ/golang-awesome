package main

import (
	"github.com/Shopify/sarama"
	"fmt"
	"log"
	"os"
	"time"
)

func main(){

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	var logger = log.New(os.Stderr, "[kafka]", log.LstdFlags)
	sarama.Logger = logger
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil{
		logger.Println("producer close,err:",err)
		return
	}

	defer client.Close()

	for {
		msg := &sarama.ProducerMessage{}
		msg.Topic = "test"
		msg.Value = sarama.StringEncoder("this is a test massage, my massage is good")

		partition, offset, e := client.SendMessage(msg)
		if e != nil {
			logger.Println("send massage failed," , err)
			return
		}

		fmt.Printf("partition :%v ,offset:%v\n" ,partition,offset)
		time.Sleep(time.Second *10)
	}




}