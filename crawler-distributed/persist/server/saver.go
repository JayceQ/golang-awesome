package main

import (
	"fmt"
	"golang-awesome/crawler-distributed/config"
	"golang-awesome/crawler-distributed/persist"
	"golang-awesome/crawler-distributed/rpcsupport"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

func main(){
	fmt.Printf("item saver start at: %v", config.ItemSaverPort)
	log.Fatal(serveRpc(config.ItemSaverPort,config.ElasticIndex))
}


func serveRpc(host, index string)error {
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetBasicAuth("elastic", "changeme"))

	if err != nil {
		return err
	}

	return rpcsupport.ServeRpc(host,
		&persist.ItemSaverService{
			Client:client,
			Index:index,
		},
	)
}