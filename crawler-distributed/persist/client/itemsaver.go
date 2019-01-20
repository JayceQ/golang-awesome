package client

import (
	"golang-awesome/crawler-distributed/rpcsupport"
	"golang-awesome/crawler/engine"
	"log"
)

//json rpc client
//将数据通过rpc传送给rpc-server
func ItemSaver(host string)(chan engine.Item,error){
	ch := make (chan engine.Item,1024)

	rpc, err := rpcsupport.NewClient(host)
	go func() {
		itemCount := 0
		for item := range ch {
			itemCount ++
			log.Printf("item saver: got item #%d: %v",itemCount, item)

			result := ""
			rpc.Call("ItemSaverService.Save",item,&result)

			if err != nil {
				log.Printf("Item Saver: Save error: %v",err)
			}
		}
	}()

	return ch,err
}