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
	templates := []string{
		"./template/Help.tmpl",
		"./template/State.tmpl",
		"./template/Box.tmpl",
		"./template/Missing.tmpl",
		"./template/Gacha.tmpl",
		"./template/Card.tmpl",
		"./template/Base.tmpl",
		"./template/Headhunt.tmpl",
	}
	r.LoadHTMLFiles(templates...)
	port := viper.GetString("http.port")
	err := r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
