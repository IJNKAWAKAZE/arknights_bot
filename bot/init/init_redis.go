package init

import (
	"arknights_bot/bot/config"
	"github.com/go-redis/redis/v8"
	"log"
)

var GoRedis *redis.Client

func Redis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.GetString("redis.addr"),
		Password: config.GetString("redis.pwd"),
		DB:       int(config.GetInt64("redis.db")),
		PoolSize: int(config.GetInt64("redis.pool_size")),
	})
	GoRedis = rdb
	log.Println("redis连接成功")
}
