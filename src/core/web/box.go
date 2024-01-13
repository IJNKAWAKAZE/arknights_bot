package web

import (
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type BoxInfo struct {
	Name  string `json:"name"`
	Chars []Char `json:"chars"`
}

type Char struct {
	CharId        string `json:"charId"`        // 角色id
	SkinId        string `json:"skinId"`        // 皮肤id
	Name          string `json:"name"`          // 名字
	Level         int    `json:"level"`         // 等级
	EvolvePhase   int    `json:"evolvePhase"`   // 精英等级
	PotentialRank int    `json:"potentialRank"` // 潜能等级
	FavorPercent  int    `json:"favorPercent"`  // 亲密度
	Rarity        int    `json:"rarity"`        // 稀有度
	Profession    string `json:"profession"`    // 职业
}

func Box(r *gin.Engine) {
	r.GET("/box", func(c *gin.Context) {
		var box BoxInfo
		var userAccount account.UserAccount
		var skAccount skland.Account
		userId, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
		uid := c.Query("uid")
		utils.GetAccountByUserId(userId).Scan(&userAccount)
		skAccount.Hypergryph.Token = userAccount.HypergryphToken
		skAccount.Skland.Token = userAccount.SklandToken
		skAccount.Skland.Cred = userAccount.SklandCred
		playerData, _, err := skland.GetPlayerInfo(uid, skAccount)
		if err != nil {
			log.Println(err)
			return
		}

		var chars []Char

		for _, c := range playerData.Chars {
			char := Char{
				CharId:        c.CharID,
				SkinId:        c.SkinID,
				Name:          playerData.CharInfoMap[c.CharID].Name,
				Level:         c.Level,
				EvolvePhase:   c.EvolvePhase,
				PotentialRank: c.PotentialRank,
				FavorPercent:  c.FavorPercent,
				Rarity:        playerData.CharInfoMap[c.CharID].Rarity,
				Profession:    playerData.CharInfoMap[c.CharID].Profession,
			}
			chars = append(chars, char)
		}

		// 按稀有度、精英等级、级别排序
		sort.Slice(chars, func(i, j int) bool {
			if chars[i].Rarity > chars[j].Rarity {
				return true
			}
			if chars[i].Rarity < chars[j].Rarity {
				return false
			}
			if chars[i].EvolvePhase > chars[j].EvolvePhase {
				return true
			}
			if chars[i].EvolvePhase < chars[j].EvolvePhase {
				return false
			}
			return chars[i].Level > chars[j].Level
		})

		box.Name = playerData.Status.Name
		box.Chars = chars

		c.HTML(http.StatusOK, "Box.tmpl", box)
	})
}
