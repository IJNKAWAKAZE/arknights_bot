package web

import (
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/player"
	"arknights_bot/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type GachaLog struct {
	Name      string      `json:"name"`
	Total     int         `json:"total"`
	PoolCount []PoolCount `json:"poolCount"`
	Star6     int         `json:"star6"`
	Star5     int         `json:"star5"`
	Star4     int         `json:"star4"`
	Star3     int         `json:"star3"`
	Avg6      float64     `json:"avg6"`
	Avg5      float64     `json:"avg5"`
	Avg4      float64     `json:"avg4"`
	Avg3      float64     `json:"avg3"`
	Chars     []GachaChar `json:"chars"`
	Star6Info []Star6Info `json:"Star6Info"`
	BegTime   int64       `json:"begTime"`
	EndTime   int64       `json:"endTime"`
}

type GachaChar struct {
	PoolName string `json:"poolName"`
	CharName string `json:"charName"`
	Avatar   string `json:"avatar"`
	IsNew    bool   `json:"isNew"`
	Rarity   int64  `json:"rarity"`
	Ts       int64  `json:"ts"`
}
type Star6Info struct {
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	Ts        int64  `json:"ts"`
	Count     int    `json:"count"`
	IsNew     bool   `json:"isNew"`
	PoolName  string `json:"poolName"`
	PoolOrder int    `json:"poolOrder"`
}

type PoolCount struct {
	PoolName  string `json:"poolName"`
	PoolCount int    `json:"count"`
}

func Gacha(r *gin.Engine) {
	r.GET("/gacha", func(c *gin.Context) {
		r.LoadHTMLFiles("./template/Gacha.tmpl")
		var gachaLog GachaLog
		var userGacha []player.UserGacha
		var gachaChars []GachaChar
		var star6Info []Star6Info
		var poolCount []PoolCount
		var PoolMap = make(map[string][]player.UserGacha)
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

		for i := range userGacha {
			var gachaChar GachaChar
			gachaChar.PoolName = userGacha[i].PoolName
			gachaChar.CharName = userGacha[i].CharName
			gachaChar.IsNew = userGacha[i].IsNew
			gachaChar.Rarity = userGacha[i].Rarity
			gachaChar.Ts = userGacha[i].Ts

			c := userGacha[len(userGacha)-i-1]
			operatorList := PoolMap[c.PoolName]
			PoolMap[c.PoolName] = append(operatorList, c)
			switch c.Rarity {
			case 5:
				star6++
				star6Info = append(star6Info, Star6Info{
					Name:      c.CharName,
					Count:     0,
					Ts:        c.Ts,
					IsNew:     c.IsNew,
					PoolName:  c.PoolName,
					PoolOrder: c.PoolOrder,
				})
			case 4:
				star5++
			case 3:
				star4++
			case 2:
				star3++
			}
			gachaChars = append(gachaChars, gachaChar)
		}

		total := len(userGacha)
		gachaLog.Total = total
		gachaLog.Star6 = star6
		gachaLog.Star5 = star5
		gachaLog.Star4 = star4
		gachaLog.Star3 = star3
		gachaLog.Avg6, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(total)/float64(star6)), 64)
		gachaLog.Avg5, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(total)/float64(star5)), 64)
		gachaLog.Avg4, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(total)/float64(star4)), 64)
		gachaLog.Avg3, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(total)/float64(star3)), 64)

		gachaLog.Chars = gachaChars
		if len(gachaChars) > 20 {
			gachaLog.Chars = gachaChars[:20]
		}
		for i := range gachaLog.Chars {
			gachaLog.Chars[i].Avatar = utils.GetOperatorByName(gachaLog.Chars[i].CharName).Avatar
		}
		gachaLog.BegTime = userGacha[len(userGacha)-1].Ts
		gachaLog.EndTime = userGacha[0].Ts

		var userPlayer account.UserPlayer
		utils.GetPlayerByUserId(userId, uid).Scan(&userPlayer)
		gachaLog.Name = userPlayer.PlayerName

		count := 1
		for k, v := range PoolMap {
			count = 1
			for i, s := range v {
				if s.Rarity == 5 {
					PoolMap[k][i].Remark = strconv.Itoa(count)
					count = 1
					continue
				}
				count++
			}
		}

		for i, s6 := range star6Info {
			for _, m := range PoolMap[s6.PoolName] {
				if s6.Name == m.CharName && s6.Ts == m.Ts && s6.PoolOrder == m.PoolOrder {
					star6Info[i].Count, _ = strconv.Atoi(m.Remark)
				}
			}
		}
		utils.ReverseSlice(star6Info)
		gachaLog.Star6Info = star6Info
		if len(star6Info) > 20 {
			gachaLog.Star6Info = star6Info[:20]
		}
		for i := range gachaLog.Star6Info {
			gachaLog.Star6Info[i].Avatar = utils.GetOperatorByName(gachaLog.Star6Info[i].Name).Avatar
		}

		utils.GetUserPoolCount(userId, uid).Scan(&poolCount)
		gachaLog.PoolCount = poolCount
		if len(poolCount) > 10 {
			gachaLog.PoolCount = poolCount[len(poolCount)-10:]
		}

		c.HTML(http.StatusOK, "Gacha.tmpl", gachaLog)
	})
}
