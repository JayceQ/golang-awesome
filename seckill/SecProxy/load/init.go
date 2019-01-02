package load

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"

	etcd "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"golang-awesome/seckill/SecProxy/service"
	"time"
)

var (
	//redisClient *redis.Client
	//redisPool *redis.Pool
	etcdClient *etcd.Client
)

//func initRedis()(err error){
//	redisClient = redis.NewClient(&redis.Options{
//		Addr: secKillConf.RedisConf.RedisAddr,
//		Password: "123456", // no password set
//		DB:       0,  // use default DB
//	})
//
//	pong, err := redisClient.Ping().Result()
//
//	fmt.Println(pong, err)
//
//	if err != nil {
//		logs.Error("ping redis failed,err :%v",err)
//		return
//	}
//
//	return
//}

//func initRedis() (err error) {
//	redisPool = &redis.Pool{
//		MaxIdle:     secKillConf.RedisConf.RedisMaxIdle,
//		MaxActive:   secKillConf.RedisConf.RedisMaxActive,
//		IdleTimeout: time.Duration(secKillConf.RedisConf.RedisIdleTimeout) * time.Second,
//		Dial: func() (redis.Conn, error) {
//			redis.Dial("tcp", secKillConf.RedisConf.RedisAddr)
//			c, err := redis.Dial("tcp", secKillConf.RedisConf.RedisAddr)
//			if err != nil {
//				return nil, err
//			}
//			if _, err := c.Do("AUTH", secKillConf.RedisConf.RedisPwd); err != nil {
//				c.Close()
//				return nil, err
//			}
//			return c, err
//		},
//	}
//
//	conn := redisPool.Get()
//	defer conn.Close()
//
//	_, err = conn.Do("ping")
//	if err != nil {
//		logs.Error("ping redis failed, err:%v", err)
//		return
//	}
//
//	return
//}


func initEtcd()(err error){
	logs.Debug("etcd timeout: %v ",time.Duration(secKillConf.EtcdConf.Timeout) * time.Second)
	client, err := etcd.New(etcd.Config{
		Endpoints:   []string{secKillConf.EtcdConf.EtcdAddr},
		DialTimeout: time.Duration(secKillConf.EtcdConf.Timeout) * time.Millisecond,
	})

	if err != nil {
		logs.Error("connect etcd failed ,err:", err)
		return
	}

	etcdClient = client
	return
}

func convertLogLevel (level string) int {
	switch level{
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	default:
		return logs.LevelDebug
	}
}


func initLogger() (err error){
	config := make(map[string]interface{})
	config["filename"] = secKillConf.LogPath
	config["level"] = convertLogLevel(secKillConf.LogLevel)

	bytes, err := json.Marshal(config)
	if err != nil {
		fmt.Println(" marshal failed,err: ",err)
		return
	}

	logs.SetLogger(logs.AdapterFile,string(bytes))
	logs.SetLogger(logs.AdapterConsole,string(bytes))

	return
}


func loadSecConf() (err error){
	response, err := etcdClient.Get(context.Background(), secKillConf.EtcdConf.EtcdSecProductKey)
	if err != nil {
		logs.Debug("get [%s] from etcd failed, err:%v ",secKillConf.EtcdConf.EtcdSecProductKey,err)
		return
	}
	var secProductInfo []service.SecProductInfoConf
	for k,v := range response.Kvs{
		logs.Debug("key[%s] value[%s]",k,v)
		err = json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
			logs.Debug("Unmarshal sec producr info failed, err:%v ",err)
			return
		}
		logs.Debug("sec info conf is [%v]",secProductInfo)
	}
	updateSecProductInfo(secProductInfo)
	return
}

func InitSecKill()(err error){
	err = initLogger()
	if err != nil{
		fmt.Println("load   logger failed,err: %v",err)
		return
	}

	//err = initRedis()
	//if err != nil{
	//	fmt.Println("load redis failed,err: %v",err)
	//	return
	//}

	err = initEtcd()
	if err != nil{
		fmt.Println("load etcd failed,err: %v",err)
		return
	}

	err = loadSecConf()
	if err != nil {
		logs.Error("load sec conf failed, err:%v", err)
		return
	}
	service.InitService(secKillConf)

	initSecProductWatcher()

	logs.Info("load sec success")
	return 
}

func initSecProductWatcher(){
	go watchSecProductKey(secKillConf.EtcdConf.EtcdSecProductKey)
}

func watchSecProductKey(key string) {
	logs.Debug("begin watch key: %s",key)

	for {
		watchChan := etcdClient.Watch(context.Background(), key)
		var secProductInfo []service.SecProductInfoConf
		var getConfSucc = true

		for wresp := range watchChan {
			for _,ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE{
					logs.Warn("key[%s]'s config deleted",key)
					continue
				}
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err := json.Unmarshal(ev.Kv.Value,&secProductInfo)
					if err != nil {
						logs.Error("key[%s],Unmarshal[%s] failed,err: %v",ev.Kv.Key,ev.Kv.Value,err)
						getConfSucc = false
						continue
					}
				}
				logs.Debug("get config from etcd,%s %q: %q\n",ev.Type,ev.Kv.Key,ev.Kv.Value)
			}

			if getConfSucc {
				logs.Debug("get config from etcd success, %v",secProductInfo)
				updateSecProductInfo(secProductInfo)
			}
		}
	}
}

func updateSecProductInfo(confs []service.SecProductInfoConf) {
	var tmp  = make(map[int]*service.SecProductInfoConf,1024)

	for _,v :=range confs{
		product := v
		tmp[v.ProductId] = &product
	}

	secKillConf.RWSecProductLock.Lock()
	secKillConf.SecProductInfoMap = tmp
	secKillConf.RWSecProductLock.Unlock()
}