package web

import (
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

type PlayerCard struct {
	Name              string `json:"name"`
	Uid               string `json:"uid"`
	Level             int    `json:"level"`
	RegTime           int    `json:"regTime"`
	MainStageProgress string `json:"mainStageProgress"`
	Avatar            string `json:"avatar"`
	Assistant         string `json:"assistant"`
	CharCnt           int    `json:"charCnt"`
	FurnitureCnt      int    `json:"furnitureCnt"`
	SkinCnt           int    `json:"skinCnt"`
	AssistChars       []struct {
		CharID          string `json:"charId"`
		SkinID          string `json:"skinId"`
		Level           int    `json:"level"`
		EvolvePhase     int    `json:"evolvePhase"`
		PotentialRank   int    `json:"potentialRank"`
		SkillID         string `json:"skillId"`
		MainSkillLvl    int    `json:"mainSkillLvl"`
		SpecializeLevel int    `json:"specializeLevel"`
		Equip           struct {
			ID    string `json:"id"`
			Level int    `json:"level"`
		} `json:"equip"`
	} `json:"assistChars"`
}

func Card(r *gin.Engine) {
	r.GET("/card", func(c *gin.Context) {
		var playerCard PlayerCard
		var playerData skland.PlayerData
		file, _ := os.Open("aaa.txt")
		readAll, _ := io.ReadAll(file)
		json.Unmarshal(readAll, &playerData)

		playerCard.Name = playerData.Status.Name
		playerCard.Uid = playerData.Status.UID
		playerCard.Level = playerData.Status.Level
		playerCard.RegTime = playerData.Status.RegisterTs
		playerCard.MainStageProgress = playerData.Status.MainStageProgress
		playerCard.Avatar = playerData.Status.Secretary.SkinID
		operatorName := playerData.CharInfoMap[playerData.Status.Secretary.CharID].Name
		playerCard.Assistant = utils.GetOperatorByName(operatorName).Painting
		playerCard.CharCnt = len(playerData.Chars)
		playerCard.SkinCnt = len(playerData.Skins)
		playerCard.FurnitureCnt = playerData.Building.Furniture.Total
		playerCard.AssistChars = playerData.AssistChars
		c.HTML(http.StatusOK, "Card.tmpl", playerCard)
	})
}
