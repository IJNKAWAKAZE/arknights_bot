package config

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var DBEngine *gorm.DB

func DB() error {
	dsn := viper.GetString("mysql.dsn")
	engine, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{Logger: logger.New(nil, logger.Config{LogLevel: logger.Silent})},
	)
	if err != nil {
		log.Println(err)
		return err
	}
	engine.Logger.LogMode(logger.Silent)
	DBEngine = engine
	log.Println("数据库连接成功")
	return nil
}
