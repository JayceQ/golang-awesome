package main

import (
	"flag"
	"golang-awesome/crawler-distributed/config"
	"golang-awesome/crawler-distributed/persist/client"
	"golang-awesome/crawler-distributed/rpcsupport"
	worker "golang-awesome/crawler-distributed/worker/client"
	"golang-awesome/crawler/engine"
	"golang-awesome/crawler/scheduler"
	"golang-awesome/crawler/zhenai/parser"
	"log"
	"net/rpc"
	"strings"
)

func createClientPool(hosts []string) chan *rpc.Client{
	var clients []*rpc.Client
	for _, h := range hosts{
		client, err := rpcsupport.NewClient(h)
		if err != nil {
			log.Printf("error connection to %s : %s",h,err)
		}else {
			clients = append(clients, client)
			log.Printf("connected to %s",h)
		}
	}
	out := make(chan *rpc.Client)
	//持续分发客户端
	go func() {
		for{
			for _, c := range clients {
				out <- c
			}
		}
	}()
	return out
}

var hosts = flag.String("hosts", ":9002","多个工作节点的端口，以都好分隔，如:9002,:9003")

func main(){
	flag.Parse()
	itemChan, err := client.ItemSaver(config.ItemSaverPort)
	if err != nil{
		panic(err)
	}
	log.Printf("connected item saver at: %v",config.ItemSaverPort)

	pool := createClientPool(strings.Split(*hosts,","))

	processor:= worker.CreateProcessor(pool)

	var seed []engine.Request
	seed = []engine.Request{
		{
			Url:"http://www.zhenai.com/zhenghun",
			Parse: engine.NewFuncParser(parser.ParseCityList,"ParseCityList"),
		},
		//{
		//	Url:       "http://www.zhenai.com/zhenghun/henan",
		//	Parse: engine.NewFuncParser(parser.ParseCity, "ParseCity"),
		//},
		//{
		//	Url:   "http://www.zhenai.com/zhenghun/beijing",
		//	Parse: engine.NewFuncParser(parser.ParseCity, "ParseCity"),
		//},
	}

	e := engine.ConcurrentEngine{
		MaxWorkerCount:100,
		Scheduler:&scheduler.QueuedScheduler{},
		ItemChan:itemChan,
		RequestWorker:processor,
	}
	e.Run(seed...)
}
