package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
)

var GoRedis *redis.Client

func Redis() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.pwd"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool_size"),
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return err
	}
	GoRedis = rdb
	log.Println("redis连接成功")
	return nil
}
