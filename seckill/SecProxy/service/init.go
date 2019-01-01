package service

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

func InitService(serviceConf *SecKillConf)(err error){
	secKillConf = serviceConf

	err = loadBalckList()
	if err != nil{
		logs.Error("load black list err: %v",err)
		return
	}
	logs.Debug("init service success, config:%v",serviceConf)

	err = initProxy2LayerRedis()
	if err != nil{
		logs.Error("load proxy2layer redis pool failed, err: %v",err)
		return
	}

	secKillConf.secLimitMgr = &SecLimitMgr{
		UserLimitMap:make(map[int]*Limit,10000),
		IpLimitMap: make(map[string]*Limit,10000),
	}
	secKillConf.SecReqChan = make(chan *SecRequest,secKillConf.SecReqChanSize)
	secKillConf.UserConnMap = make (map[string]chan *SecResult,10000)

	initRedisProcessFunc()
	return
}

func loadBalckList()(err error){
	secKillConf.ipBlackMap = make(map[string]bool,10000)
	secKillConf.idBlackMap = make(map[int]bool,10000)

	err = initBlackRedis()

	if err != nil{
		logs.Error("init black redis failed, err:%v",err)
		return
	}

	conn := secKillConf.blackRedisPool.Get()
	defer conn.Close()

	reply, err := conn.Do("hgetall", "id_black_list")
	idlist,err :=redis.Strings(reply,err)
	if err != nil {
		logs.Warn("hget all black id failed, err %v",err)
		return
	}

	for _,v :=range idlist{
		id, err := strconv.Atoi(v)
		if err != nil {
			logs.Warn("invalid user id [%v]",id)
			continue
		}
		secKillConf.idBlackMap[id] = true
	}

	reply, err = conn.Do("hgetall", "ip_black_list")
	iplist,err :=redis.Strings(reply,err)
	if err != nil {
		logs.Warn("hget all black ip failed, err %v",err)
		return
	}

	for _,v :=range iplist{
		secKillConf.ipBlackMap[v] = true
	}

	go syncIpBLackList()
	go syncIdBLackList()
	return
}

//用redis的BLPOP阻塞弹出,达到watch效果
func syncIpBLackList(){
	var ipList []string
	lastTime := time.Now().Unix()
	for{
		conn := secKillConf.blackRedisPool.Get()
		defer conn.Close()
		replay,err := conn.Do("BLPOP","blackiplist",time.Second)
		ip,err :=redis.String(replay,err)
		if err != nil {
			continue
		}

		curTime := time.Now().Unix()
		ipList = append(ipList,ip)

		if len(ipList) >100 || curTime -lastTime >5{
			secKillConf.RWBlackLock.Lock()
			for _,v :=range ipList{
				secKillConf.ipBlackMap[v] = true
			}
			secKillConf.RWBlackLock.Unlock()

			lastTime = curTime
			logs.Info("sync ip list from redis success, ip[%v]",ipList)
		}
	}
}

func syncIdBLackList(){
	for{
		conn := secKillConf.blackRedisPool.Get()
		defer conn.Close()
		reply,err := conn.Do("BLPOP","blackidlist",time.Second)
		id,err := redis.Int(reply,err)
		if err != nil{
			continue
		}
		secKillConf.RWBlackLock.Lock()
		secKillConf.idBlackMap[id] = true
		secKillConf.RWBlackLock.Unlock()
		logs.Info("sync id list from redis success, ip[%v]",id)
	}
}

func initRedisProcessFunc(){
	for i :=0 ;i<secKillConf.WriteProxy2LayerGoroutineNum;i++{
		go WriteHandle()
	}

	for i :=0 ;i<secKillConf.ReadProxy2LayerGoroutineNum;i++{
		go ReadHandle()
	}
}

func initProxy2LayerRedis()(err error){
	secKillConf.proxy2LayerRedisPool = &redis.Pool{
		MaxIdle:secKillConf.RedisProxy2LayerConf.RedisMaxIdle,
		MaxActive: secKillConf.RedisProxy2LayerConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillConf.RedisProxy2LayerConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", secKillConf.RedisProxy2LayerConf.RedisAddr)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", secKillConf.RedisProxy2LayerConf.RedisPwd); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},

	}

	conn := secKillConf.proxy2LayerRedisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err:%v", err)
		return
	}

	return
}

func initLayer2ProxyRedis()(err error){
	secKillConf.layer2ProxyRedisPool = &redis.Pool{
		MaxIdle:secKillConf.RedisLayer2ProxyConf.RedisMaxIdle,
		MaxActive: secKillConf.RedisLayer2ProxyConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillConf.RedisLayer2ProxyConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", secKillConf.RedisLayer2ProxyConf.RedisAddr)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", secKillConf.RedisLayer2ProxyConf.RedisPwd); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},

	}

	conn := secKillConf.layer2ProxyRedisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err:%v", err)
		return
	}

	return
}

func initBlackRedis()(err error){
	secKillConf.blackRedisPool = &redis.Pool{
		MaxIdle:secKillConf.RedisBlackConf.RedisMaxIdle,
		MaxActive: secKillConf.RedisBlackConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillConf.RedisBlackConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", secKillConf.RedisBlackConf.RedisAddr)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", secKillConf.RedisBlackConf.RedisPwd); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},

	}

	conn := secKillConf.blackRedisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err:%v", err)
		return
	}

	return
}