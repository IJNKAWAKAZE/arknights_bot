package web

import (
	"arknights_bot/plugins/player"
	"arknights_bot/plugins/skland"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func State(r *gin.Engine) {
	r.GET("/state", func(c *gin.Context) {
		r.LoadHTMLFiles("./template/State.tmpl")
		userId, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
		uid := c.Query("uid")
		sklandId := c.Query("sklandId")
		playerData, userAccount, skAccount, err := player.GetPlayerData(userId, sklandId, uid)
		if err != nil {
			log.Println(err)
			renderError(c, err)
			return
		}
		playStatistic, _, err := skland.GetPlayerStatistic(uid, skAccount, userAccount.ServerName)
		if err != nil {
			log.Println(err)
			renderError(c, err)
			return
		}

		playStatistic.Avatar = playerData.Status.Secretary.SkinID

		c.HTML(http.StatusOK, "State.tmpl", playStatistic)
	})
}
