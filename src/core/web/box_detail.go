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

type Detail struct {
	Name          string  `json:"name"`
	Id            string  `json:"id"`
	Rarity        int     `json:"rarity"`
	Level         int     `json:"level"`
	EvolvePhase   int     `json:"evolvePhase"`
	PotentialRank int     `json:"potentialRank"`
	Skills        []Skill `json:"skills"`
	Equips        []Equip `json:"equips"`
}

type Skill struct {
	Id    string `json:"id"`
	Level int    `json:"level"`
}

type Equip struct {
	Id    string `json:"id"`
	Level int    `json:"level"`
}

func BoxDetail(r *gin.Engine) {
	r.GET("/boxDetail", func(c *gin.Context) {
		r.LoadHTMLFiles("./template/BoxDetail.tmpl")
		var detailList []Detail
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
		playerCultivate, err := skland.GetPlayerCultivate(uid, skAccount)
		if err != nil {
			log.Println(err)
			return
		}
		playerData, _, err := skland.GetPlayerInfo(uid, skAccount)
		if err != nil {
			log.Println(err)
			return
		}

		for _, char := range playerCultivate.Characters {
			rarity := playerData.CharInfoMap[char.ID].Rarity
			if filter(param, rarity) {
				var detail Detail
				name := playerData.CharInfoMap[char.ID].Name
				if char.ID == "char_1001_amiya2" {
					name = "阿米娅(近卫)"
				}
				if char.ID == "char_1037_amiya3" {
					name = "阿米娅(医疗)"
				}
				detail.Name = name
				for _, c := range playerData.Chars {
					if char.ID == c.CharID {
						detail.Id = c.SkinID
					}
				}
				detail.Rarity = rarity
				detail.Level = char.Level
				detail.EvolvePhase = char.EvolvePhase
				detail.PotentialRank = char.PotentialRank
				var skills []Skill
				for _, sk := range char.Skills {
					var skill Skill
					if sk.Level == 0 {
						skill.Level = char.MainSkillLevel
					} else {
						skill.Level = sk.Level + 7
					}
					skill.Id = sk.ID
					skills = append(skills, skill)
				}
				detail.Skills = skills
				var equips []Equip
				for _, eq := range char.Equips {
					var equip Equip
					if eq.Level == 0 {
						continue
					}
					equip.Id = playerData.EquipmentInfoMap[eq.ID].TypeIcon
					equip.Level = eq.Level
					equips = append(equips, equip)
				}
				detail.Equips = equips
				detailList = append(detailList, detail)
			}
		}
		sort.Slice(detailList, func(i, j int) bool {
			if detailList[i].EvolvePhase > detailList[j].EvolvePhase {
				return true
			}
			if detailList[i].EvolvePhase < detailList[j].EvolvePhase {
				return false
			}
			if detailList[i].Level > detailList[j].Level {
				return true
			}
			if detailList[i].Level < detailList[j].Level {
				return false
			}
			return detailList[i].Rarity > detailList[j].Rarity
		})
		c.HTML(http.StatusOK, "BoxDetail.tmpl", detailList)
	})
}
