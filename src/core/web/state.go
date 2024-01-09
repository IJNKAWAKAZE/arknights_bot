package web

import (
	"arknights_bot/plugins/skland"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func State(r *gin.Engine) {
	r.GET("/state/:data/:uid", func(c *gin.Context) {
		var account skland.Account
		uid := c.Param("uid")
		json.Unmarshal([]byte(c.Param("data")), &account)
		playerData, err := skland.GetPlayerInfo(uid, account)
		if err != nil {
			log.Println(err)
			return
		}
		playStatistic, err := skland.GetPlayerStatistic(uid, account)
		if err != nil {
			log.Println(err)
			return
		}

		playStatistic.Avatar = playerData.Status.Secretary.SkinID

		c.HTML(http.StatusOK, "State.tmpl", playStatistic)
	})
}
