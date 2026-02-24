package web

import (
	"arknights_bot/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Lottery(r *gin.Engine) {
	r.GET("/lotteryDetail", func(c *gin.Context) {
		r.LoadHTMLFiles("./template/Lottery.tmpl")
		lotteryId := c.Query("lotteryId")
		var details []utils.GroupLotteryDetail
		utils.GetLotteryDetails(lotteryId).Scan(&details)
		c.HTML(http.StatusOK, "Lottery.tmpl", details)
	})
}
