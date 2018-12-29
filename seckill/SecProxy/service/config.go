package service

import "sync"

type RedisConf struct {
	RedisAddr string
	RedisPwd string
	RedisMaxIdle int
	RedisMaxActive int
	RedisIdleTimeout int
}

type EtcdConf struct {
	EtcdAddr string
	Timeout int
	EtcdSecKeyPrefix string
	EtcdSecProductKey string
}

type SecProductInfoConf struct {
	ProductId int
	StartTime int
	EndTime int
	Status int
	Total int
	Left int
}

type SecKillConf struct {
	RedisConf RedisConf
	EtcdConf EtcdConf
	LogPath string
	LogLevel string
	SecProductInfoMap map[int]*SecProductInfoConf
	RWSecProductLock sync.RWMutex
}