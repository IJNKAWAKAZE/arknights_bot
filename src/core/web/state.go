package web

import (
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func State(r *gin.Engine) {
	r.GET("/state", func(c *gin.Context) {
		var userAccount account.UserAccount
		var skAccount skland.Account
		userId, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
		uid := c.Query("uid")
		utils.GetAccountByUserId(userId).Scan(&userAccount)
		skAccount.Hypergryph.Token = userAccount.HypergryphToken
		skAccount.Skland.Token = userAccount.SklandToken
		skAccount.Skland.Cred = userAccount.SklandCred
		playerData, skAccount, err := skland.GetPlayerInfo(uid, skAccount)
		if err != nil {
			log.Println(err)
			return
		}
		playStatistic, _, err := skland.GetPlayerStatistic(uid, skAccount)
		if err != nil {
			log.Println(err)
			return
		}

		playStatistic.Avatar = playerData.Status.Secretary.SkinID

		c.HTML(http.StatusOK, "State.tmpl", playStatistic)
	})
}
