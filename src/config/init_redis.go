package config

import (
	"github.com/go-redis/redis/v8"
	"log"
)

var GoRedis *redis.Client

func Redis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     GetString("redis.addr"),
		Password: GetString("redis.pwd"),
		DB:       int(GetInt64("redis.db")),
		PoolSize: int(GetInt64("redis.pool_size")),
	})
	GoRedis = rdb
	log.Println("redis连接成功")
}
