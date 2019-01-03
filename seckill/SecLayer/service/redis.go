package service

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"time"
)

func initRedisPool(redisConf RedisConf) (pool *redis.Pool, err error) {
	pool = &redis.Pool{
		MaxIdle:     redisConf.RedisMaxIdle,
		MaxActive:   redisConf.RedisMaxActive,
		IdleTimeout: time.Duration(redisConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisConf.RedisAddr)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", redisConf.RedisPwd); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
	}

	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err:%v", err)
		return
	}

	return
}


func initRedis(conf *SecLayerConf) (err error) {

	secLayerContext.proxy2LayerRedisPool, err = initRedisPool(conf.Proxy2LayerRedis)
	if err != nil {
		logs.Error("init proxy2layer redis pool failed, err:%v", err)
		return
	}

	secLayerContext.layer2ProxyRedisPool, err = initRedisPool(conf.Layer2ProxyRedis)
	if err != nil {
		logs.Error("init layer2proxy redis pool failed, err:%v", err)
		return
	}

	return
}

func RunProcess() (err error) {

	for i := 0; i < secLayerContext.secLayerConf.ReadGoroutineNum; i++ {
		secLayerContext.waitGroup.Add(1)
		go HandleReader()
	}

	for i := 0; i < secLayerContext.secLayerConf.WriteGoroutineNum; i++ {
		secLayerContext.waitGroup.Add(1)
		go HandleWrite()
	}

	for i := 0; i < secLayerContext.secLayerConf.HandleUserGoroutineNum; i++ {
		secLayerContext.waitGroup.Add(1)
		go HandleUser()
	}

	logs.Debug("all process goroutine started")
	secLayerContext.waitGroup.Wait()
	logs.Debug("wait all goroutine exited")
	return
}

func HandleReader() {

	logs.Debug("read goroutine running")
	for {
		conn := secLayerContext.proxy2LayerRedisPool.Get()
		for {
			ret, err := conn.Do("brpop", secLayerContext.secLayerConf.Proxy2LayerRedis.RedisQueueName, 0)
			if err != nil {
				logs.Error("pop from queue failed, err:%v", err)
				break
			}

			tmp, ok := ret.([]interface{})
			if !ok || len(tmp) != 2{
				logs.Error("pop from queue failed, err:%v", err)
				continue
			}

			data, ok := tmp[1].([]byte)
			if !ok {
				logs.Error("pop from queue failed, err:%v", err)
				continue
			}

			logs.Debug("pop from queue, data:%s", string(data))

			var req SecRequest
			err = json.Unmarshal([]byte(data), &req)
			if err != nil {
				logs.Error("unmarshal to secrequest failed, err:%v", err)
				continue
			}

			now := time.Now().Unix()
			if now-req.AccessTime.Unix() >= int64(secLayerContext.secLayerConf.MaxRequestWaitTimeout) {
				logs.Warn("req[%v] is expire", req)
				continue
			}

			timer := time.NewTicker(time.Millisecond * time.Duration(secLayerContext.secLayerConf.SendToHandleChanTimeout))
			select {
			case secLayerContext.Read2HandleChan <- &req:
			case <-timer.C:
				logs.Warn("send to handle chan timeout, req:%v", req)
				break
			}
		}

		conn.Close()
	}
}

func HandleWrite() {
	logs.Debug("handle write running")

	for res := range secLayerContext.Handle2WriteChan {
		err := sendToRedis(res)
		if err != nil {
			logs.Error("send to redis, err:%v, res:%v", err, res)
			continue
		}
	}
}

func sendToRedis(res *SecResponse) (err error) {

	data, err := json.Marshal(res)
	if err != nil {
		logs.Error("marshal failed, err:%v", err)
		return
	}

	conn := secLayerContext.layer2ProxyRedisPool.Get()
	_, err = conn.Do("lpush", secLayerContext.secLayerConf.Layer2ProxyRedis.RedisQueueName, string(data))
	if err != nil {
		logs.Warn("rpush to redis failed, err:%v", err)
		return
	}

	return
}

func HandleUser() {

	logs.Debug("handle user running")
	for req := range secLayerContext.Read2HandleChan {
		logs.Debug("begin process request:%v", req)
		res, err := HandleSecKill(req)
		if err != nil {
			logs.Warn("process request %v failed, err:%v", err)
			res = &SecResponse{
				Code: ErrServiceBusy,
			}
		}

		timer := time.NewTicker(time.Millisecond * time.Duration(secLayerContext.secLayerConf.SendToWriteChanTimeout))
		select {
		case secLayerContext.Handle2WriteChan <- res:
		case <-timer.C:
			logs.Warn("send to response chan timeout, res:%v", res)
			break
		}

	}
	return
}

func HandleSecKill(req *SecRequest) (res *SecResponse, err error) {

	secLayerContext.RWSecProductLock.RLock()
	defer secLayerContext.RWSecProductLock.RUnlock()

	res = &SecResponse{}
	res.UserId = req.UserId
	res.ProductId = req.ProductId
	product, ok := secLayerContext.secLayerConf.SecProductInfoMap[req.ProductId]
	if !ok {
		logs.Error("not found product:%v", req.ProductId)
		res.Code = ErrNotFoundProduct
		return
	}

	if product.Status == ProductStatusSoldout {
		res.Code = ErrSoldout
		return
	}

	now := time.Now().Unix()
	alreadySoldCount := product.secLimit.Check(now)
	if alreadySoldCount >= product.SoldMaxLimit {
		res.Code = ErrRetry
		return
	}

	secLayerContext.HistoryMaoLock.Lock()
	userHistory, ok := secLayerContext.HistoryMap[req.UserId]
	if !ok {
		userHistory = &UserBuyHistory{
			history: make(map[int]int, 16),
		}

		secLayerContext.HistoryMap[req.UserId] = userHistory
	}

	histryCount := userHistory.GetProductBuyCount(req.ProductId)
	secLayerContext.HistoryMaoLock.Unlock()

	if histryCount >= product.OnePersonBuyLimit {
		res.Code = ErrAlreadyBuy
		return
	}

	curSoldCount := secLayerContext.productCountMgr.Count(req.ProductId)
	if curSoldCount >= product.Total {
		res.Code = ErrSoldout
		product.Status = ProductStatusSoldout
		return
	}


	curRate := rand.Float64()
	fmt.Printf("curRate:%v product:%v count:%v total:%v\n", curRate, product.BuyRate, curSoldCount, product.Total)
	if curRate > product.BuyRate {
		res.Code = ErrRetry
		return
	}

	userHistory.Add(req.ProductId, 1)
	secLayerContext.productCountMgr.Add(req.ProductId, 1)

	//用户id&商品id&当前时间&密钥
	res.Code = ErrSecKillSucc
	tokenData := fmt.Sprintf("userId=%d&productId=%d&timestamp=%d&security=%s",
		req.UserId, req.ProductId, now, secLayerContext.secLayerConf.TokenPassword)

	res.Token = fmt.Sprintf("%x", md5.Sum([]byte(tokenData)))
	logs.Warn(res.Token)
	res.TokenTime = now

	return
}
