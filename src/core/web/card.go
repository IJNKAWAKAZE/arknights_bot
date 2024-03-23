package web

import (
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type PlayerCard struct {
	Name              string `json:"name"`
	Uid               string `json:"uid"`
	ServerName        string `json:"serverName"`
	Resume            string `json:"resume"`
	Level             int    `json:"level"`
	RegTime           int    `json:"regTime"`
	MainStageProgress string `json:"mainStageProgress"`
	Avatar            string `json:"avatar"`
	CharCnt           int    `json:"charCnt"`
	FurnitureCnt      int    `json:"furnitureCnt"`
	SkinCnt           int    `json:"skinCnt"`
	AssistChars       []struct {
		Name            string `json:"name"`
		CharID          string `json:"charId"`
		SkinID          string `json:"skinId"`
		Level           int    `json:"level"`
		EvolvePhase     int    `json:"evolvePhase"`
		PotentialRank   int    `json:"potentialRank"`
		SkillID         string `json:"skillId"`
		MainSkillLvl    int    `json:"mainSkillLvl"`
		SpecializeLevel int    `json:"specializeLevel"`
		Equip           struct {
			ID           string `json:"id"`
			Level        int    `json:"level"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"equip"`
	} `json:"assistChars"`
}

func Card(r *gin.Engine) {
	r.GET("/card", func(c *gin.Context) {
		r.LoadHTMLFiles("./template/Card.tmpl")
		var playerCard PlayerCard
		var userAccount account.UserAccount
		var skAccount skland.Account
		userId, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
		uid := c.Query("uid")
		utils.GetAccountByUserId(userId).Scan(&userAccount)
		var userPlayer account.UserPlayer
		utils.GetPlayerByUserId(userAccount.UserNumber, uid).Scan(&userPlayer)
		playerCard.ServerName = userPlayer.ServerName
		playerCard.Resume = userPlayer.Resume
		skAccount.Hypergryph.Token = userAccount.HypergryphToken
		skAccount.Skland.Token = userAccount.SklandToken
		skAccount.Skland.Cred = userAccount.SklandCred
		playerData, skAccount, err := skland.GetPlayerInfo(uid, skAccount)
		if err != nil {
			log.Println(err)
			return
		}
		playerCard.Name = playerData.Status.Name
		playerCard.Uid = playerData.Status.UID
		playerCard.Level = playerData.Status.Level
		playerCard.RegTime = playerData.Status.RegisterTs
		playerCard.MainStageProgress = playerData.StageInfoMap[playerData.Status.MainStageProgress].Code
		playerCard.Avatar = playerData.Status.Secretary.SkinID
		playerCard.CharCnt = len(playerData.Chars)
		playerCard.SkinCnt = len(playerData.Skins)
		playerCard.FurnitureCnt = playerData.Building.Furniture.Total
		playerCard.AssistChars = playerData.AssistChars
		for i, char := range playerCard.AssistChars {
			name := playerData.CharInfoMap[char.CharID].Name
			equipId := playerCard.AssistChars[i].Equip.ID
			playerCard.AssistChars[i].Name = name
			playerCard.AssistChars[i].Equip.Name = strings.ToUpper(playerData.EquipmentInfoMap[equipId].TypeIcon)
			playerCard.AssistChars[i].Equip.TypeIcon = playerData.EquipmentInfoMap[equipId].TypeIcon
			playerCard.AssistChars[i].Equip.ShiningColor = playerData.EquipmentInfoMap[equipId].ShiningColor
			playerCard.AssistChars[i].MainSkillLvl = playerCard.AssistChars[i].MainSkillLvl + playerCard.AssistChars[i].SpecializeLevel
		}
		c.HTML(http.StatusOK, "Card.tmpl", playerCard)
	})
}
