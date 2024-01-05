package config

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
)

var GoRedis *redis.Client

func Redis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.pwd"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool_size"),
	})
	GoRedis = rdb
	log.Println("redis连接成功")
}
