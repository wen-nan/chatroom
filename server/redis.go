package main

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

// 定义一个全局的pool
var pool *redis.Pool

func initPool(address string, maxIdle, maxActivate int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle: maxIdle, // 最大空闲连接数
		MaxActive: maxActivate, // 表示和数据库最大连接数
		IdleTimeout: idleTimeout, // 最大空闲时间
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
}