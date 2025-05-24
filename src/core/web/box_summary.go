package web

import (
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type BoxSummary struct {
	Name               string `json:"name"`               // 角色名
	AllCharCnt         string `json:"allCharCnt"`         //干员招募数
	AllEvolvePhase2Cnt int    `json:"allEvolvePhase2Cnt"` // 精二数
	AllSkill10Cnt      int    `json:"allSkill10Cnt"`      // 专3技能数
	AllSkill9Cnt       int    `json:"allSkill9Cnt"`       // 专2技能数
	AllSkill8Cnt       int    `json:"allSkill8Cnt"`       // 专1技能数
	AllEquipStage3Cnt  int    `json:"allEquipStage3Cnt"`  // 3级模组数
	AllEquipStage2Cnt  int    `json:"allEquipStage2Cnt"`  // 2级模组数
	AllEquipStage1Cnt  int    `json:"allEquipStage1Cnt"`  // 1级模组数
	// 六星
	Star6CharCnt         string `json:"star6CharCnt"`         //干员招募数
	Star6EvolvePhase2Cnt int    `json:"star6EvolvePhase2Cnt"` // 精二数
	Star6Skill10Cnt      int    `json:"star6Skill10Cnt"`      // 专3技能数
	Star6Skill9Cnt       int    `json:"star6Skill9Cnt"`       // 专2技能数
	Star6Skill8Cnt       int    `json:"star6Skill8Cnt"`       // 专1技能数
	Star6EquipStage3Cnt  int    `json:"star6EquipStage3Cnt"`  // 3级模组数
	Star6EquipStage2Cnt  int    `json:"star6EquipStage2Cnt"`  // 2级模组数
	Star6EquipStage1Cnt  int    `json:"star6EquipStage1Cnt"`  // 1级模组数
	// 五星
	Star5CharCnt         string `json:"star5CharCnt"`         //干员招募数
	Star5EvolvePhase2Cnt int    `json:"star5EvolvePhase2Cnt"` // 精二数
	Star5Skill10Cnt      int    `json:"star5Skill10Cnt"`      // 专3技能数
	Star5Skill9Cnt       int    `json:"star5Skill9Cnt"`       // 专2技能数
	Star5Skill8Cnt       int    `json:"star5Skill8Cnt"`       // 专1技能数
	Star5EquipStage3Cnt  int    `json:"star5EquipStage3Cnt"`  // 3级模组数
	Star5EquipStage2Cnt  int    `json:"star5EquipStage2Cnt"`  // 2级模组数
	Star5EquipStage1Cnt  int    `json:"star5EquipStage1Cnt"`  // 1级模组数
	// 四星
	Star4CharCnt         string        `json:"star4CharCnt"`         //干员招募数
	Star4EvolvePhase2Cnt int           `json:"star4EvolvePhase2Cnt"` // 精二数
	Star4Skill10Cnt      int           `json:"star4Skill10Cnt"`      // 专3技能数
	Star4Skill9Cnt       int           `json:"star4Skill9Cnt"`       // 专2技能数
	Star4Skill8Cnt       int           `json:"star4Skill8Cnt"`       // 专1技能数
	Star4EquipStage3Cnt  int           `json:"star4EquipStage3Cnt"`  // 3级模组数
	Star4EquipStage2Cnt  int           `json:"star4EquipStage2Cnt"`  // 2级模组数
	Star4EquipStage1Cnt  int           `json:"star4EquipStage1Cnt"`  // 1级模组数
	MissingChars         []MissingChar `json:"missingChars"`         //未招募干员
}

func Summary(r *gin.Engine) {
	r.GET("/boxSummary", func(c *gin.Context) {
		r.LoadHTMLFiles("./template/BoxSummary.tmpl")
		userId, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
		uid := c.Query("uid")
		sklandId := c.Query("sklandId")
		var boxSummary BoxSummary
		var userAccount account.UserAccount
		var skAccount skland.Account
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
		playerCultivate, err := skland.GetPlayerCultivate(uid, skAccount)
		if err != nil {
			log.Println(err)
			utils.WebC <- err
			return
		}
		var missingChars []MissingChar
		operatorList := utils.GetOperators()
		myOperators := make(map[string]Char)
		charMap := playerData.CharInfoMap
		boxSummary.AllCharCnt = fmt.Sprintf("%d/%d", len(myOperators), len(operatorList))
		for _, c := range playerData.Chars {
			rarity := playerData.CharInfoMap[c.CharID].Rarity
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
			myOperators[name] = char
		}
		allCharCnt := len(operatorList) - 2 // 排除两个阿米娅形态
		myCharCnt := len(myOperators)
		star6Cnt := 0
		star5Cnt := 0
		star4Cnt := 0
		myStar6Cnt := 0
		myStar5Cnt := 0
		myStar4Cnt := 0
		for _, operator := range operatorList {
			name := operator.Name
			rarity := operator.Rarity
			myOperator := myOperators[name]
			switch rarity {
			case 5:
				star6Cnt++
				if check(myOperators, name) {
					myStar6Cnt++
					if myOperator.EvolvePhase == 2 {
						boxSummary.Star6EvolvePhase2Cnt++
					}
				} else {
					missingChars = append(missingChars, MissingChar{
						SkinId:     operator.Avatar,
						Name:       name,
						Rarity:     rarity,
						Profession: operator.Profession,
					})
				}
			case 4:
				star5Cnt++
				if check(myOperators, name) {
					myStar5Cnt++
					if myOperator.EvolvePhase == 2 {
						boxSummary.Star5EvolvePhase2Cnt++
					}
				} else {
					missingChars = append(missingChars, MissingChar{
						SkinId:     operator.Avatar,
						Name:       name,
						Rarity:     rarity,
						Profession: operator.Profession,
					})
				}
				break
			case 3:
				star4Cnt++
				if check(myOperators, name) {
					myStar4Cnt++
					if myOperator.EvolvePhase == 2 {
						boxSummary.Star4EvolvePhase2Cnt++
					}
				} else {
					missingChars = append(missingChars, MissingChar{
						SkinId:     operator.Avatar,
						Name:       name,
						Rarity:     rarity,
						Profession: operator.Profession,
					})
				}
				break
			default:
				if !check(myOperators, name) {
					missingChars = append(missingChars, MissingChar{
						SkinId:     operator.Avatar,
						Name:       name,
						Rarity:     rarity,
						Profession: operator.Profession,
					})
				}
			}
		}
		if _, has := charMap["char_1001_amiya2"]; has {
			myCharCnt--
			star5Cnt--
			myStar5Cnt--
		}
		if _, has := charMap["char_1037_amiya3"]; has {
			myCharCnt--
			star5Cnt--
			myStar5Cnt--
		}
		for _, char := range playerCultivate.Characters {
			rarity := charMap[char.ID].Rarity
			switch rarity {
			case 5:
				for _, skill := range char.Skills {
					if skill.Level == 3 {
						boxSummary.Star6Skill10Cnt++
					}
					if skill.Level == 2 {
						boxSummary.Star6Skill9Cnt++
					}
					if skill.Level == 1 {
						boxSummary.Star6Skill8Cnt++
					}
				}
				for _, equip := range char.Equips {
					if equip.Level == 3 {
						boxSummary.Star6EquipStage3Cnt++
					}
					if equip.Level == 2 {
						boxSummary.Star6EquipStage2Cnt++
					}
					if equip.Level == 1 {
						boxSummary.Star6EquipStage1Cnt++
					}
				}
				break
			case 4:
				for _, skill := range char.Skills {
					if skill.Level == 3 {
						boxSummary.Star5Skill10Cnt++
					}
					if skill.Level == 2 {
						boxSummary.Star5Skill9Cnt++
					}
					if skill.Level == 1 {
						boxSummary.Star5Skill8Cnt++
					}
				}
				for _, equip := range char.Equips {
					if equip.Level == 3 {
						boxSummary.Star5EquipStage3Cnt++
					}
					if equip.Level == 2 {
						boxSummary.Star5EquipStage2Cnt++
					}
					if equip.Level == 1 {
						boxSummary.Star5EquipStage1Cnt++
					}
				}
				break
			case 3:
				for _, skill := range char.Skills {
					if skill.Level == 3 {
						boxSummary.Star4Skill10Cnt++
					}
					if skill.Level == 2 {
						boxSummary.Star4Skill9Cnt++
					}
					if skill.Level == 1 {
						boxSummary.Star4Skill8Cnt++
					}
				}
				for _, equip := range char.Equips {
					if equip.Level == 3 {
						boxSummary.Star4EquipStage3Cnt++
					}
					if equip.Level == 2 {
						boxSummary.Star4EquipStage2Cnt++
					}
					if equip.Level == 1 {
						boxSummary.Star4EquipStage1Cnt++
					}
				}
				break
			}
		}
		boxSummary.Name = playerData.Status.Name
		boxSummary.AllCharCnt = fmt.Sprintf("%d/%d", myCharCnt, allCharCnt)
		boxSummary.Star6CharCnt = fmt.Sprintf("%d/%d", myStar6Cnt, star6Cnt)
		boxSummary.Star5CharCnt = fmt.Sprintf("%d/%d", myStar5Cnt, star5Cnt)
		boxSummary.Star4CharCnt = fmt.Sprintf("%d/%d", myStar4Cnt, star4Cnt)
		boxSummary.AllEvolvePhase2Cnt = boxSummary.Star6EvolvePhase2Cnt + boxSummary.Star5EvolvePhase2Cnt + boxSummary.Star4EvolvePhase2Cnt
		boxSummary.AllSkill10Cnt = boxSummary.Star6Skill10Cnt + boxSummary.Star5Skill10Cnt + boxSummary.Star4Skill10Cnt
		boxSummary.AllSkill9Cnt = boxSummary.Star6Skill9Cnt + boxSummary.Star5Skill9Cnt + boxSummary.Star4Skill9Cnt
		boxSummary.AllSkill8Cnt = boxSummary.Star6Skill8Cnt + boxSummary.Star5Skill8Cnt + boxSummary.Star4Skill8Cnt
		boxSummary.AllEquipStage3Cnt = boxSummary.Star6EquipStage3Cnt + boxSummary.Star5EquipStage3Cnt + boxSummary.Star4EquipStage3Cnt
		boxSummary.AllEquipStage2Cnt = boxSummary.Star6EquipStage2Cnt + boxSummary.Star5EquipStage2Cnt + boxSummary.Star4EquipStage2Cnt
		boxSummary.AllEquipStage1Cnt = boxSummary.Star6EquipStage1Cnt + boxSummary.Star5EquipStage1Cnt + boxSummary.Star4EquipStage1Cnt
		boxSummary.MissingChars = missingChars

		// 按稀有度、职业排序
		sort.Slice(boxSummary.MissingChars, func(i, j int) bool {
			if boxSummary.MissingChars[i].Rarity > boxSummary.MissingChars[j].Rarity {
				return true
			}
			if boxSummary.MissingChars[i].Rarity < boxSummary.MissingChars[j].Rarity {
				return false
			}
			return boxSummary.MissingChars[i].Profession > boxSummary.MissingChars[j].Profession
		})
		c.HTML(http.StatusOK, "BoxSummary.tmpl", boxSummary)
	})
}

func check(myOperators map[string]Char, name string) bool {
	if _, has := myOperators[name]; !has {
		return false
	}
	return true
}
