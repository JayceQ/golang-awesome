package client

import (
	"golang-awesome/crawler-distributed/config"
	"golang-awesome/crawler-distributed/worker"
	"golang-awesome/crawler/engine"
	"net/rpc"
)

func CreateProcessor(clientChan chan*rpc.Client)engine.Processor{
	return func(request engine.Request) (result engine.ParserResult, e error) {
		var sReq = worker.SerializeRequest(request)
		var sResult worker.ParseResult
		c := <-clientChan
		err := c.Call(config.CrawlServiceRpc,sReq, &sResult)
		if err != nil{
			return engine.ParserResult{},err
		}
		return worker.DeserializeResult(sResult), nil
	}
}
