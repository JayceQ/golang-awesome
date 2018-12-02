package main

import "github.com/garyburd/redigo/redis"
import "time"

var pool *redis.Pool

func initRedis(addr string, idleConn, maxConn int, idleTimeout time.Duration) {

	pool = &redis.Pool{
		MaxIdle:     idleConn,
		MaxActive:   maxConn,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
	}
	return
}

func GetConn() redis.Conn {
	return pool.Get()
}

func PutConn(conn redis.Conn) {
	conn.Close()
}
