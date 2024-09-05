package web

import (
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
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
		r.LoadHTMLFiles("./template/Box.tmpl")
		var box BoxInfo
		var userAccount account.UserAccount
		var skAccount skland.Account
		userId, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
		uid := c.Query("uid")
		param := c.Query("param")
		sklandId := c.Query("sklandId")
		utils.GetAccountByUserIdAndSklandId(userId, sklandId).Scan(&userAccount)
		skAccount.Hypergryph.Token = userAccount.HypergryphToken
		skAccount.Skland.Token = userAccount.SklandToken
		skAccount.Skland.Cred = userAccount.SklandCred
		playerData, _, err := skland.GetPlayerInfo(uid, skAccount)
		if err != nil {
			log.Println(err)
			utils.WebC <- err
			return
		}

		var chars []Char

		for _, c := range playerData.Chars {
			rarity := playerData.CharInfoMap[c.CharID].Rarity
			if filter(param, rarity) {
				name := playerData.CharInfoMap[c.CharID].Name
				if c.CharID == "char_1001_amiya2" {
					name = "阿米娅(近卫)"
				}
				if c.CharID == "char_1037_amiya3" {
					name = "阿米娅(医疗)"
				}
				char := Char{
					CharId:        c.CharID,
					SkinId:        c.SkinID,
					Name:          name,
					Level:         c.Level,
					EvolvePhase:   c.EvolvePhase,
					PotentialRank: c.PotentialRank,
					FavorPercent:  c.FavorPercent,
					Rarity:        rarity,
					Profession:    playerData.CharInfoMap[c.CharID].Profession,
				}
				chars = append(chars, char)
			}
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
			if chars[i].Level > chars[j].Level {
				return true
			}
			if chars[i].Level < chars[j].Level {
				return false
			}
			return chars[i].Profession > chars[j].Profession
		})

		box.Name = playerData.Status.Name
		box.Chars = chars
		if len(box.Chars) == 0 {
			utils.WebC <- fmt.Errorf("无符合干员")
			return
		}

		c.HTML(http.StatusOK, "Box.tmpl", box)
	})
}

func filter(param string, rarity int) bool {
	switch param {
	case "":
		if rarity == 5 {
			return true
		}
	case "all":
		return true
	default:
		matched, _ := regexp.MatchString("^[1-6](,[1-6])*$", param)
		if matched {
			nums := strings.Split(param, ",")
			for _, num := range nums {
				r, _ := strconv.Atoi(num)
				if r == rarity+1 {
					return true
				}
			}
		}
	}
	return false
}
