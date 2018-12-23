package main

import (
	"context"
	"fmt"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

type LogMessage struct {
	App     string
	Topic   string
	Message string
}

var (
	esClient *elastic.Client
	count int
)

func initES(addr string) (err error) {

	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(addr))
	if err != nil {
		fmt.Println("connect es error", err)
		return
	}
	esClient = client
	return
	/*
		fmt.Println("conn es succ")
		for i := 0; i < 10000; i++ {
			tweet := Tweet{User: "olivere", Message: "Take Five"}
			_, err = client.Index().
				Index("twitter").
				Type("tweet").
				Id(fmt.Sprintf("%d", i)).
				BodyJson(tweet).
				Do()
			if err != nil {
				// Handle error
				panic(err)
				return
			}
		}

		fmt.Println("insert succ")
	*/
}

func sendToES(topic string, data []byte) (err error) {

	msg := &LogMessage{}
	msg.Topic = topic
	msg.Message = string(data)
	count++
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	_, err = esClient.Index().
		Index(topic).
		Type(topic).
		Id(fmt.Sprintf("%d", count)).
		BodyJson(msg).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
		return
	}
	return
}
