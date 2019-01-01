package load

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"golang-awesome/seckill/SecProxy/service"
	"strings"
)

var (
	secKillConf = &service.SecKillConf{
		SecProductInfoMap: make(map[int]*service.SecProductInfoConf,1024),
	}
)
func InitConfig()(err error){
	//黑名单redis配置
	redisBlackAddr := beego.AppConfig.String("redis_black_addr")
	redisBlackPwd := beego.AppConfig.String("redis_black_pwd")
	etcdAddr := beego.AppConfig.String("etcd_addr")

	logs.Debug("read redis config success, redis addr: %v",redisBlackAddr)
	logs.Debug("read etcd config success, etcd addr: %v",etcdAddr)

	secKillConf.RedisBlackConf.RedisAddr = redisBlackAddr
	secKillConf.RedisBlackConf.RedisPwd = redisBlackPwd
	secKillConf.EtcdConf.EtcdAddr = etcdAddr

	if len(redisBlackAddr) ==0 || len(etcdAddr) == 0{
		err = fmt.Errorf("load config failed, redis[%s] or etcd[%s] is nil",redisBlackAddr,etcdAddr)
		return
	}

	redisMaxIdle, err := beego.AppConfig.Int("redis_black_idle")
	if err != nil {
		err = fmt.Errorf("load config failed, read redis_black_idle error: %v",err)
	}

	redisMaxActive, err := beego.AppConfig.Int("redis_black_active")
	if err != nil {
		err = fmt.Errorf("load config failed, read redis_black_active error: %v",err)
	}

	redisIdleTimeout, err := beego.AppConfig.Int("redis_black_idle_timeout")
	if err != nil {
		err = fmt.Errorf("load config failed, read redis_black_idle_timeout error: %v",err)
	}

	secKillConf.RedisBlackConf.RedisMaxIdle = redisMaxIdle
	secKillConf.RedisBlackConf.RedisMaxActive = redisMaxActive
	secKillConf.RedisBlackConf.RedisIdleTimeout = redisIdleTimeout

	etcdTimeout, err := beego.AppConfig.Int("etcd_timeout")
	if err != nil {
		err = fmt.Errorf("load config failed, read etcd_timeout error:%v", err)
		return
	}

	secKillConf.EtcdConf.Timeout = etcdTimeout
	secKillConf.EtcdConf.EtcdSecKeyPrefix = beego.AppConfig.String("etcd_sec_key_prefix")
	if len(secKillConf.EtcdConf.EtcdSecKeyPrefix) == 0 {
		err = fmt.Errorf("load config failed, read etcd_sec_key error:%v", err)
		return
	}

	productKey := beego.AppConfig.String("etcd_product_key")
	if len(productKey) == 0 {
		err = fmt.Errorf("load config failed, read etcd_product_key error:%v", err)
		return
	}

	if strings.HasSuffix(secKillConf.EtcdConf.EtcdSecKeyPrefix, "/") == false {
		secKillConf.EtcdConf.EtcdSecKeyPrefix = secKillConf.EtcdConf.EtcdSecKeyPrefix + "/"
	}

	secKillConf.EtcdConf.EtcdSecProductKey = fmt.Sprintf("%s%s", secKillConf.EtcdConf.EtcdSecKeyPrefix, productKey)
	secKillConf.LogPath = beego.AppConfig.String("log_path")
	secKillConf.LogLevel = beego.AppConfig.String("log_level")

	secKillConf.CookieSecretKey = beego.AppConfig.String("cookie_secretkey")
	secLimit, err := beego.AppConfig.Int("user_sec_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read user_sec_access_limit error:%v", err)
		return
	}

	secKillConf.AccessLimitConf.UserSecAccessLimit = secLimit
	referList := beego.AppConfig.String("refer_whitelist")
	if len(referList) > 0 {
		secKillConf.ReferWhiteList = strings.Split(referList, ",")
	}

	ipLimit, err := beego.AppConfig.Int("ip_sec_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read ip_sec_access_limit error:%v", err)
		return
	}

	secKillConf.AccessLimitConf.IPSecAccessLimit = ipLimit

	minIdLimit, err := beego.AppConfig.Int("user_min_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read user_min_access_limit error:%v", err)
		return
	}

	secKillConf.AccessLimitConf.UserMinAccessLimit = minIdLimit
	minIpLimit, err := beego.AppConfig.Int("ip_min_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read ip_min_access_limit error:%v", err)
		return
	}

	secKillConf.AccessLimitConf.IPMinAccessLimit = minIpLimit

	redisProxy2LayerAddr := beego.AppConfig.String("redis_proxy2layer_addr")
	logs.Debug("read config succ, redis addr:%v", redisProxy2LayerAddr)

	secKillConf.RedisProxy2LayerConf.RedisAddr = redisProxy2LayerAddr

	if len(redisProxy2LayerAddr) == 0 {
		err = fmt.Errorf("init config failed, redis[%s]  config is null", redisProxy2LayerAddr)
		return
	}

	redisPwd := beego.AppConfig.String("redis_proxy2layer_pwd")
	if len(redisPwd) == 0 {
		err = fmt.Errorf("init config failed, read redis_layer2proxy_pwd error:%v", err)
		return
	}

	redisMaxIdle, err = beego.AppConfig.Int("redis_proxy2layer_idle")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_proxy2layer_idle error:%v", err)
		return
	}

	redisMaxActive, err = beego.AppConfig.Int("redis_proxy2layer_active")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_proxy2layer_active error:%v", err)
		return
	}

	redisIdleTimeout, err = beego.AppConfig.Int("redis_proxy2layer_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_proxy2layer_idle_timeout error:%v", err)
		return
	}

	secKillConf.RedisProxy2LayerConf.RedisPwd= redisPwd
	secKillConf.RedisProxy2LayerConf.RedisMaxIdle = redisMaxIdle
	secKillConf.RedisProxy2LayerConf.RedisMaxActive = redisMaxActive
	secKillConf.RedisProxy2LayerConf.RedisIdleTimeout = redisIdleTimeout

	writeGoNums, err := beego.AppConfig.Int("write_proxy2layer_goroutine_num")
	if err != nil {
		err = fmt.Errorf("init config failed, read write_proxy2layer_goroutine_num error:%v", err)
		return
	}

	secKillConf.WriteProxy2LayerGoroutineNum = writeGoNums

	readGoNums, err := beego.AppConfig.Int("read_layer2proxy_goroutine_num")
	if err != nil {
		err = fmt.Errorf("init config failed, read read_layer2proxy_goroutine_num error:%v", err)
		return
	}

	secKillConf.ReadProxy2LayerGoroutineNum = readGoNums

	//读取业务逻辑层到proxy的redis配置
	redisLayer2ProxyAddr := beego.AppConfig.String("redis_layer2proxy_addr")
	logs.Debug("read config succ, redis addr:%v", redisLayer2ProxyAddr)

	secKillConf.RedisProxy2LayerConf.RedisAddr = redisLayer2ProxyAddr

	if len(redisLayer2ProxyAddr) == 0 {
		err = fmt.Errorf("init config failed, redis[%s]  config is null", redisProxy2LayerAddr)
		return
	}

	redisPwd = beego.AppConfig.String("redis_layer2proxy_pwd")
	if len(redisPwd) == 0 {
		err = fmt.Errorf("init config failed, read redis_layer2proxy_pwd error:%v", err)
		return
	}

	redisMaxIdle, err = beego.AppConfig.Int("redis_layer2proxy_idle")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_layer2proxy_idle error:%v", err)
		return
	}

	redisMaxActive, err = beego.AppConfig.Int("redis_layer2proxy_active")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_layer2proxy_active error:%v", err)
		return
	}

	redisIdleTimeout, err = beego.AppConfig.Int("redis_layer2proxy_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_layer2proxy_idle_timeout error:%v", err)
		return
	}

	secKillConf.RedisLayer2ProxyConf.RedisPwd= redisPwd
	secKillConf.RedisLayer2ProxyConf.RedisMaxIdle = redisMaxIdle
	secKillConf.RedisLayer2ProxyConf.RedisMaxActive = redisMaxActive
	secKillConf.RedisLayer2ProxyConf.RedisIdleTimeout = redisIdleTimeout
	return
}