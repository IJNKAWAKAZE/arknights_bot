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
		utils.WebC = make(chan error, 10)
		defer close(utils.WebC)
		r.LoadHTMLFiles("./template/State.tmpl")
		var userAccount account.UserAccount
		var skAccount skland.Account
		userId, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
		uid := c.Query("uid")
		sklandId := c.Query("sklandId")
		utils.GetAccountByUserIdAndSklandId(userId, sklandId).Scan(&userAccount)
		skAccount.Hypergryph.Token = userAccount.HypergryphToken
		skAccount.Skland.Token = userAccount.SklandToken
		skAccount.Skland.Cred = userAccount.SklandCred
		playerData, skAccount, err := skland.GetPlayerInfo(uid, skAccount)
		if err != nil {
			log.Println(err)
			utils.WebC <- err
			return
		}
		playStatistic, _, err := skland.GetPlayerStatistic(uid, skAccount)
		if err != nil {
			log.Println(err)
			utils.WebC <- err
			return
		}

		playStatistic.Avatar = playerData.Status.Secretary.SkinID

		c.HTML(http.StatusOK, "State.tmpl", playStatistic)
	})
}
