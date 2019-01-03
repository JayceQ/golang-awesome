package service

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	etcd "go.etcd.io/etcd/clientv3"
	"strings"
	"sync"
	"time"
)

var (
	AppConfig *SecLayerConf
	secLayerContext = &SecLayerContext{}
)
type SecProductInfoConf = struct {
	ProductId int
	StartTime int64
	EndTime int64
	Status int
	Total int
	Left int
	OnePersonBuyLimit int
	BuyRate float64
	//每秒组多能卖多少个
	SoldMaxLimit int
	//限速控制
	secLimit *SecLimit
}

type RedisConf struct {
	RedisAddr string
	RedisPwd string
	RedisMaxIdle int
	RedisMaxActive int
	RedisIdleTimeout int
	RedisQueueName string
}

type EtcdConf struct {
	EtcdAddr string
	Timeout int
	EtcdSecKeyPrefix string
	EtcdSecProductKey string
}

type SecLayerConf struct{
	Proxy2LayerRedis RedisConf
	Layer2ProxyRedis RedisConf
	EtcdConfig EtcdConf
	LogPath string
	LogLevel string

	WriteGoroutineNum int
	ReadGoroutineNum int
	HandleUserGoroutineNum int
	Read2HandleChanSize int
	Handle2WriteChanSize int
	MaxRequestWaitTimeout int

	SendToWriteChanTimeout int
	SendToHandleChanTimeout int

	SecProductInfoMap map[int]*SecProductInfoConf
	TokenPassword string
}

type SecLayerContext struct{
	proxy2LayerRedisPool *redis.Pool
	layer2ProxyRedisPool *redis.Pool
	etcdClient *etcd.Client
	RWSecProductLock sync.RWMutex

	secLayerConf *SecLayerConf
	waitGroup sync.WaitGroup
	Read2HandleChan chan *SecRequest
	Handle2WriteChan chan *SecResponse

	HistoryMap map[int]*UserBuyHistory
	HistoryMaoLock sync.Mutex

	//商品的计数
	productCountMgr *ProductCountMgr
}
type SecRequest struct {
	ProductId     int
	Source        string
	AuthCode      string
	SecTime       string
	Nance         string
	UserId        int
	UserAuthSign  string
	AccessTime    time.Time
	ClientAddr    string
	ClientReferer string
}

type SecResponse struct {
	ProductId int
	UserId    int
	Token     string
	TokenTime int64
	Code      int
}



func InitConfig(confType,fileName string )(err error){


	conf, err := config.NewConfig(confType, fileName)
	if err != nil {
		fmt.Println("new config failed, err:", err)
		return
	}

	//读取日志库配置
	AppConfig = &SecLayerConf{}
	AppConfig.LogLevel = conf.String("logs::log_level")
	if len(AppConfig.LogLevel) == 0 {
		AppConfig.LogLevel = "debug"
	}

	AppConfig.LogPath = conf.String("logs::log_path")
	if len(AppConfig.LogPath) == 0 {
		AppConfig.LogPath = "./logs"
	}

	//读取redis相关的配置
	AppConfig.Proxy2LayerRedis.RedisAddr = conf.String("redis::redis_proxy2layer_addr")
	if len(AppConfig.Proxy2LayerRedis.RedisAddr) == 0 {
		logs.Error("read redis::redis_proxy2layer_addr failed")
		err = fmt.Errorf("read redis::redis_proxy2layer_addr failed")
		return
	}

	AppConfig.Proxy2LayerRedis.RedisPwd = conf.String("redis::redis_proxy2layer_pwd")
	if len(AppConfig.Proxy2LayerRedis.RedisPwd) == 0 {
		logs.Error("read redis::redis_proxy2layer_pwd failed")
		err = fmt.Errorf("read redis::redis_proxy2layer_pwd failed")
		return
	}

	AppConfig.Proxy2LayerRedis.RedisQueueName = conf.String("redis::redis_proxy2layer_queue_name")
	if len(AppConfig.Proxy2LayerRedis.RedisQueueName) == 0 {
		logs.Error("read redis::redis_proxy2layer_queue_name failed")
		err = fmt.Errorf("read redis::redis_proxy2layer_queue_name failed")
		return
	}

	AppConfig.Proxy2LayerRedis.RedisMaxIdle, err = conf.Int("redis::redis_proxy2layer_idle")
	if err != nil {
		logs.Error("read redis::redis_proxy2layer_idle failed, err:%v", err)
		return
	}

	AppConfig.Proxy2LayerRedis.RedisIdleTimeout, err = conf.Int("redis::redis_proxy2layer_idle_timeout")
	if err != nil {
		logs.Error("read redis::redis_proxy2layer_idle_timeout failed, err:%v", err)
		return
	}

	AppConfig.Proxy2LayerRedis.RedisMaxActive, err = conf.Int("redis::redis_proxy2layer_active")
	if err != nil {
		logs.Error("read redis::redis_proxy2layer_active failed, err:%v", err)
		return
	}

	//读取redis layer2proxy相关的配置
	AppConfig.Layer2ProxyRedis.RedisAddr = conf.String("redis::redis_layer2proxy_addr")
	if len(AppConfig.Proxy2LayerRedis.RedisAddr) == 0 {
		logs.Error("read redis::redis_layer2proxy_addr failed")
		err = fmt.Errorf("read redis::redis_layer2proxy_addr failed")
		return
	}

	AppConfig.Layer2ProxyRedis.RedisPwd = conf.String("redis::redis_layer2proxy_pwd")
	if len(AppConfig.Layer2ProxyRedis.RedisPwd) == 0 {
		logs.Error("read redis::redis_layer2proxy_pwd failed")
		err = fmt.Errorf("read redis::redis_layer2proxy_pwd failed")
		return
	}

	AppConfig.Layer2ProxyRedis.RedisQueueName = conf.String("redis::redis_layer2proxy_queue_name")
	if len(AppConfig.Layer2ProxyRedis.RedisQueueName) == 0 {
		logs.Error("read redis::redis_layer2proxy_queue_name failed")
		err = fmt.Errorf("read redis::redis_layer2proxy_queue_name failed")
		return
	}

	AppConfig.Layer2ProxyRedis.RedisMaxIdle, err = conf.Int("redis::redis_layer2proxy_idle")
	if err != nil {
		logs.Error("read redis::redis_layer2proxy_idle failed, err:%v", err)
		return
	}

	AppConfig.Layer2ProxyRedis.RedisIdleTimeout, err = conf.Int("redis::redis_layer2proxy_idle_timeout")
	if err != nil {
		logs.Error("read redis::redis_layer2proxy_idle_timeout failed, err:%v", err)
		return
	}

	AppConfig.Layer2ProxyRedis.RedisMaxActive, err = conf.Int("redis::redis_layer2proxy_active")
	if err != nil {
		logs.Error("read redis::redis_layer2proxy_active failed, err:%v", err)
		return
	}

	//读取各类goroutine线程数量
	AppConfig.ReadGoroutineNum, err = conf.Int("service::read_layer2proxy_goroutine_num")
	if err != nil {
		logs.Error("read service::read_layer2proxy_goroutine_num failed, err:%v", err)
		return
	}

	AppConfig.WriteGoroutineNum, err = conf.Int("service::write_proxy2layer_goroutine_num")
	if err != nil {
		logs.Error("read service::write_proxy2layer_goroutine_num failed, err:%v", err)
		return
	}

	AppConfig.HandleUserGoroutineNum, err = conf.Int("service::handle_user_goroutine_num")
	if err != nil {
		logs.Error("read service::handle_user_goroutine_num failed, err:%v", err)
		return
	}

	AppConfig.Read2HandleChanSize, err = conf.Int("service::read2handle_chan_size")
	if err != nil {
		logs.Error("read service::read2handle_chan_size failed, err:%v", err)
		return
	}

	AppConfig.MaxRequestWaitTimeout, err = conf.Int("service::max_request_wait_timeout")
	if err != nil {
		logs.Error("read service::max_request_wait_timeout failed, err:%v", err)
		return
	}

	AppConfig.Handle2WriteChanSize, err = conf.Int("service::handle2write_chan_size")
	if err != nil {
		logs.Error("read service::handle2write_chan_size failed, err:%v", err)
		return
	}

	AppConfig.SendToWriteChanTimeout, err = conf.Int("service::send_to_write_chan_timeout")
	if err != nil {
		logs.Error("read service::send_to_write_chan_timeout failed, err:%v", err)
		return
	}

	AppConfig.SendToHandleChanTimeout, err = conf.Int("service::send_to_handle_chan_timeout")
	if err != nil {
		logs.Error("read service::send_to_handle_chan_timeout failed, err:%v", err)
		return
	}

	//读取token秘钥
	AppConfig.TokenPassword = conf.String("service::seckill_token_passwd")
	if len(AppConfig.TokenPassword) == 0 {
		logs.Error("read service::seckill_token_passwd failed")
		err = fmt.Errorf("read service::seckill_token_passwd failed")
		return
	}

	//读取etcd相关的配置

	AppConfig.EtcdConfig.EtcdAddr = conf.String("etcd::server_addr")
	if len(AppConfig.TokenPassword) == 0 {
		logs.Error("read service::seckill_token_passwd failed")
		err = fmt.Errorf("read service::seckill_token_passwd failed")
		return
	}

	etcdTimeout, err := conf.Int("etcd::etcd_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read etcd_timeout error:%v", err)
		return
	}

	AppConfig.EtcdConfig.Timeout = etcdTimeout
	AppConfig.EtcdConfig.EtcdSecKeyPrefix = conf.String("etcd::etcd_sec_key_prefix")
	if len(AppConfig.EtcdConfig.EtcdSecKeyPrefix) == 0 {
		err = fmt.Errorf("init config failed, read etcd_sec_key error:%v", err)
		return
	}

	productKey := conf.String("etcd::etcd_product_key")
	if len(productKey) == 0 {
		err = fmt.Errorf("init config failed, read etcd_product_key error:%v", err)
		return
	}

	if strings.HasSuffix(AppConfig.EtcdConfig.EtcdSecKeyPrefix, "/") == false {
		AppConfig.EtcdConfig.EtcdSecKeyPrefix = AppConfig.EtcdConfig.EtcdSecKeyPrefix + "/"
	}

	AppConfig.EtcdConfig.EtcdSecProductKey = fmt.Sprintf("%s%s", AppConfig.EtcdConfig.EtcdSecKeyPrefix, productKey)
	return
}
