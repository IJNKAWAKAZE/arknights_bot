package web

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Start() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Static("/assets", "./assets")
	r.Static("/template/js", "./template/js")
	Help(r)
	State(r)
	Box(r)
	Missing(r)
	Gacha(r)
	Card(r)
	Base(r)
	Headhunt(r)
	Operator(r)
	Enemy(r)
	Material(r)
	Recruit(r)
	Calendar(r)
	Depot(r)
	BoxDetail(r)
	Summary(r)
	Lottery(r)
	port := viper.GetString("http.port")
	err := r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
