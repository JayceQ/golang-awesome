package main

import (
	"context"
	"fmt"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

type Tweet struct {
	User    string
	Message string
}

func main() {
	client, err := elastic.NewClient(elastic.SetSniff(false),
		elastic.SetURL("http://localhost:9200/"),
		elastic.SetBasicAuth("elastic","changeme"))
	if err != nil {
		fmt.Println("connect es error:", err)
		return
	}

	fmt.Println("conn es succ")

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	tweet := Tweet{User: "olivere", Message: "Take Five"}
	_, err = client.Index().
		Index("twitter").
		Type("tweet").
		Id("1").
		BodyJson(tweet).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
		return
	}
	cancelFunc()
	fmt.Println("insert succ")
}
