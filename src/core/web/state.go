package web

import (
	"arknights_bot/plugins/skland"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func State(r *gin.Engine) {
	r.GET("/state", func(c *gin.Context) {
		var account skland.Account
		uid := c.Query("uid")
		json.Unmarshal([]byte(c.Query("data")), &account)
		playerData, account, err := skland.GetPlayerInfo(uid, account)
		if err != nil {
			log.Println(err)
			return
		}
		playStatistic, _, err := skland.GetPlayerStatistic(uid, account)
		if err != nil {
			log.Println(err)
			return
		}

		playStatistic.Avatar = playerData.Status.Secretary.SkinID

		c.HTML(http.StatusOK, "State.tmpl", playStatistic)
	})
}
