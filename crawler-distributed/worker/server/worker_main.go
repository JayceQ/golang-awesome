package main

import (
	"flag"
	"fmt"
	"golang-awesome/crawler-distributed/config"
	"golang-awesome/crawler-distributed/rpcsupport"
	"golang-awesome/crawler-distributed/worker"
	"log"
)

var port = flag.Int("port", config.WorkerPort0,"请输入默认工作端口号(默认10086)")

func main(){
	flag.Parse()
	port2 := fmt.Sprintf(":%d", *port)

	fmt.Println("Worker Server Start At:", port2)
	log.Fatal(rpcsupport.ServeRpc(port2,worker.CrawlService{}))
}
