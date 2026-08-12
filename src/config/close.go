package config

import "log"

// Close 关闭数据库与 Redis 连接
func Close() {
	if DBEngine != nil {
		if sqlDB, err := DBEngine.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Println("数据库关闭失败:", err)
			}
		}
	}
	if GoRedis != nil {
		if err := GoRedis.Close(); err != nil {
			log.Println("Redis关闭失败:", err)
		}
	}
}
