package skland

import (
	"encoding/json"
	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/tidwall/gjson"
	"log"
)

type PlayerCultivate struct {
	Characters []struct {
		ID             string `json:"id"`
		Level          int    `json:"level"`
		EvolvePhase    int    `json:"evolvePhase"`
		MainSkillLevel int    `json:"mainSkillLevel"`
		Skills         []struct {
			ID    string `json:"id"`
			Level int    `json:"level"`
		} `json:"skills"`
		Equips []struct {
			ID    string `json:"id"`
			Level int    `json:"level"`
		} `json:"equips"`
		PotentialRank int `json:"potentialRank"`
	} `json:"characters"`
	Items []struct {
		ID    string `json:"id"`
		Count string `json:"count"`
	} `json:"items"`
}

func GetPlayerCultivate(uid string, account Account, serverName string) (*PlayerCultivate, error) {
	var playerCultivate *PlayerCultivate
	account, err := RefreshToken(account, serverName)
	if err != nil {
		log.Println(err.Error())
		return playerCultivate, err
	}
	var playerCultivateStr string
	if serverName == "国服" {
		playerCultivateStr, err = getPlayerCultivateStr(uid, account.Skland)
	} else if serverName == "国际服" {
		playerCultivateStr, err = iGetPlayerCultivateStr(uid, account.Skland)
	}
	if err != nil {
		return playerCultivate, err
	}

	json.Unmarshal([]byte(gjson.Get(playerCultivateStr, "data").String()), &playerCultivate)
	return playerCultivate, nil
}

func getPlayerCultivateStr(uid string, skland AccountSkland) (string, error) {
	req := SKR().SetQueryParams(gh.MS{"uid": uid})
	return SklandRequestPlayerData(req, "GET", "/api/v1/game/cultivate/player", skland)
}

func iGetPlayerCultivateStr(uid string, skland AccountSkland) (string, error) {
	req := SKR().SetQueryParams(gh.MS{"uid": uid})
	return SkportRequestPlayerData(req, "GET", "/api/v1/game/cultivate/player", skland)
}
