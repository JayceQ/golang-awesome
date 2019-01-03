package service

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"time"
)

func WriteHandle(){
	for{
		req := <-secKillConf.SecReqChan
		conn := secKillConf.proxy2LayerRedisPool.Get()

		data ,err := json.Marshal(req)
		if err !=nil{
			logs.Error("json marshal failed, err: %v, req %v",err,req)
		}

		_, err = conn.Do("LPUSH", "sec_queue", string(data))
		if err!= nil{
			logs.Error("lpush failed, err: %v",err)
			conn.Close()
		}
		conn.Close()
	}
}

func ReadHandle(){
	for {
		conn := secKillConf.proxy2LayerRedisPool.Get()
		reply, err := conn.Do("RPOP", "recv_queue")
		data, err := redis.String(reply, err)
		if err!= nil {
			time.Sleep(time.Second)
			conn.Close()
			continue
		}

		logs.Debug("rpop from recv_queue success, data: %v",string(data))
		if err != nil{
			logs.Error("rpop failed, err %v",err)
			continue
		}

		var result SecResult
		err = json.Unmarshal([]byte(data),&result)

		if err != nil {
			logs.Error("json Unmarshal failed, err :%v",err)
			continue
		}
		userKey := fmt.Sprintf("%s_%s", result.UserId, result.ProductId)

		secKillConf.UserConnMapLock.Lock()
		resultChan, ok := secKillConf.UserConnMap[userKey]
		secKillConf.UserConnMapLock.Unlock()
		if !ok {
			conn.Close()
			logs.Warn("user not found: %s",userKey)
			continue
		}
		resultChan <- &result
		conn.Close()
	}
}
