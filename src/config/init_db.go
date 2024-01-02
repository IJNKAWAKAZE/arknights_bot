package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var DBEngine *gorm.DB

func DB() {
	dsn := GetString("mysql.dsn")
	engine, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{Logger: logger.New(nil, logger.Config{LogLevel: logger.Silent})},
	)
	if err != nil {
		log.Println(err)
		return
	}
	engine.Logger.LogMode(logger.Silent)
	DBEngine = engine
	log.Println("数据库连接成功")
}
