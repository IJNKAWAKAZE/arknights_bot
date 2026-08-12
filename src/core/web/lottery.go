package web

import (
	"arknights_bot/utils/model"
	"arknights_bot/utils/repo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Lottery(r *gin.Engine) {
	r.GET("/lotteryDetail", func(c *gin.Context) {
		r.LoadHTMLFiles("./template/Lottery.tmpl")
		lotteryId := c.Query("lotteryId")
		var details []model.GroupLotteryDetail
		repo.GetLotteryDetails(lotteryId).Scan(&details)
		c.HTML(http.StatusOK, "Lottery.tmpl", details)
	})
}
