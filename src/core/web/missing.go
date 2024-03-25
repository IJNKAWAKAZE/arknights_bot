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

type MissingInfo struct {
	Name  string        `json:"name"`
	Chars []MissingChar `json:"chars"`
}

type MissingChar struct {
	SkinId     string `json:"skinId"`     // 皮肤id
	Name       string `json:"name"`       // 名字
	Rarity     int    `json:"rarity"`     // 稀有度
	Profession string `json:"profession"` // 职业
}

func Missing(r *gin.Engine) {
	r.GET("/missing", func(c *gin.Context) {
		r.LoadHTMLFiles("./template/Missing.tmpl")
		var missingInfo MissingInfo
		param := c.Query("param")
		var userAccount account.UserAccount
		var skAccount skland.Account
		userId, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
		uid := c.Query("uid")
		sklandId := c.Query("sklandId")
		utils.GetAccountByUserIdAndSklandId(userId, sklandId).Scan(&userAccount)
		skAccount.Hypergryph.Token = userAccount.HypergryphToken
		skAccount.Skland.Token = userAccount.SklandToken
		skAccount.Skland.Cred = userAccount.SklandCred
		playerData, _, err := skland.GetPlayerInfo(uid, skAccount)
		if err != nil {
			log.Println(err)
			return
		}

		var chars []MissingChar
		myOperators := make(map[string]Char)
		operatorList := utils.GetOperators()

		for _, c := range playerData.Chars {
			rarity := playerData.CharInfoMap[c.CharID].Rarity
			if filter(param, rarity) {
				name := playerData.CharInfoMap[c.CharID].Name
				char := Char{
					CharId:        c.CharID,
					SkinId:        c.SkinID,
					Name:          playerData.CharInfoMap[c.CharID].Name,
					Level:         c.Level,
					EvolvePhase:   c.EvolvePhase,
					PotentialRank: c.PotentialRank,
					FavorPercent:  c.FavorPercent,
					Rarity:        rarity,
					Profession:    playerData.CharInfoMap[c.CharID].Profession,
				}
				myOperators[name] = char
			}
		}

		for _, operator := range operatorList {
			name := operator.Name
			if name == "阿米娅(近卫)" {
				continue
			}
			rarity := operator.Rarity
			if filter(param, rarity) {
				if _, has := myOperators[name]; !has {
					char := MissingChar{
						SkinId:     operator.ThumbURL,
						Name:       name,
						Rarity:     rarity,
						Profession: operator.Profession,
					}
					chars = append(chars, char)
				}
			}
		}

		// 按稀有度、职业排序
		sort.Slice(chars, func(i, j int) bool {
			if chars[i].Rarity > chars[j].Rarity {
				return true
			}
			if chars[i].Rarity < chars[j].Rarity {
				return false
			}
			return chars[i].Profession > chars[j].Profession
		})

		missingInfo.Name = playerData.Status.Name
		missingInfo.Chars = chars

		c.HTML(http.StatusOK, "Missing.tmpl", missingInfo)
	})
}
