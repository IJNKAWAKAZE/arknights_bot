package web

import (
	"arknights_bot/plugins/enemy"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Enemy(r *gin.Engine) {
	r.GET("/enemy", func(c *gin.Context) {
		r.LoadHTMLFiles("./template/Enemy.tmpl")
		name := c.Query("name")
		enemyInfo := enemy.ParseEnemy(name)
		c.HTML(http.StatusOK, "Enemy.tmpl", enemyInfo)
	})
}
