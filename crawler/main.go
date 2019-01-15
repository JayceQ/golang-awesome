package main

import (
	"golang-awesome/crawler/engine"
	"golang-awesome/crawler/persist"
	"golang-awesome/crawler/scheduler"
	"golang-awesome/crawler/zhenai/parser"
)

func main(){
	itemChan, err := persist.ItemSaver("profiles")
	if err != nil {
		panic(err)
	}

	var seed []engine.Request

	seed = []engine.Request{
		{
			Url:   "http://www.zhenai.com/zhenghun/beijing",
			Parse: engine.NewFuncParser(parser.ParseCity, "ParseCity"),
		},
	}
	e := engine.ConcurrentEngine{
		MaxWorkerCount: 100,
		Scheduler: &scheduler.QueuedScheduler{},
		ItemChan: itemChan,
		RequestWorker: engine.Worker,
	}

	e.Run(seed...)
}