package web

import (
	"arknights_bot/plugins/skland"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sort"
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
	r.GET("/box/:data/:uid", func(c *gin.Context) {
		var box BoxInfo
		var account skland.Account
		uid := c.Param("uid")
		json.Unmarshal([]byte(c.Param("data")), &account)
		playerData, err := skland.GetPlayerInfo(uid, account)
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

		sort.Slice(chars, func(i, j int) bool {
			return chars[i].Rarity > chars[j].Rarity
		})

		box.Name = playerData.Status.Name
		box.Chars = chars

		c.HTML(http.StatusOK, "Box.tmpl", box)
	})
}
