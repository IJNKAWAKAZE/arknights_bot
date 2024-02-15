package web

import (
	"arknights_bot/plugins/player"
	"arknights_bot/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type GachaLog struct {
	Total     int                `json:"total"`
	Current   int                `json:"current"`
	Star6     int                `json:"star6"`
	Star5     int                `json:"star5"`
	Star4     int                `json:"star4"`
	Star3     int                `json:"star3"`
	Chars     []player.UserGacha `json:"chars"`
	Star6Info []Star6Info        `json:"Star6Info"`
	BegTime   int64              `json:"begTime"`
	EndTime   int64              `json:"endTime"`
}
type Star6Info struct {
	Name     string `json:"name"`
	Ts       int64  `json:"ts"`
	Count    int    `json:"count"`
	IsNew    bool   `json:"isNew"`
	PoolName string `json:"poolName"`
}

func Gacha(r *gin.Engine) {
	r.GET("/gacha", func(c *gin.Context) {
		var gachaLog GachaLog
		var userGacha []player.UserGacha
		var star6Info []Star6Info
		userId, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
		uid := c.Query("uid")
		res := utils.GetUserGacha(userId, uid).Scan(&userGacha)
		if res.Error != nil {
			log.Println(res.Error)
			return
		}

		star6 := 0
		star5 := 0
		star4 := 0
		star3 := 0

		count := 1

		for i := range userGacha {
			c := userGacha[len(userGacha)-i-1]
			switch c.Rarity {
			case 5:
				star6++
				star6Info = append(star6Info, Star6Info{
					Name:     c.CharName,
					Count:    count,
					Ts:       c.Ts,
					IsNew:    c.IsNew,
					PoolName: c.PoolName,
				})
				count = 1
				continue
			case 4:
				star5++
			case 3:
				star4++
			case 2:
				star3++
			}
			count++
		}

		gachaLog.Total = len(userGacha)
		gachaLog.Current = count - 1
		gachaLog.Star6 = star6
		gachaLog.Star5 = star5
		gachaLog.Star4 = star4
		gachaLog.Star3 = star3
		gachaLog.Chars = userGacha
		gachaLog.BegTime = userGacha[len(userGacha)-1].Ts
		gachaLog.EndTime = userGacha[0].Ts

		utils.ReverseSlice(star6Info)
		gachaLog.Star6Info = star6Info

		c.HTML(http.StatusOK, "Gacha.tmpl", gachaLog)
	})
}
