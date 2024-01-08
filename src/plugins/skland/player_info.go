package skland

import (
	"encoding/json"
	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/tidwall/gjson"
	"log"
	"time"
)

var Weight map[string]int

func init() {
	weight := make(map[string]int)
	weight["1"] = 2
	weight["2"] = 3
	weight["3"] = 5
	weight["4"] = 2
	weight["5"] = 5
	weight["6"] = 5
	weight["7"] = 5
	weight["8"] = 5
	weight["9"] = 5
	weight["10"] = 5
	weight["11"] = 5
	weight["12"] = 5
	weight["13"] = 3
	weight["14"] = 3
	Weight = weight
}
func GetPlayerInfo(uid string, account Account) (*PlayerData, error) {
	var playerData *PlayerData
	account, err := RefreshToken(uid, account)
	if err != nil {
		log.Println(err.Error())
		return playerData, err
	}
	playerDatastr, err := getPlayerInfoStr(uid, account.Skland)
	if err != nil {
		return playerData, err
	}

	json.Unmarshal([]byte(gjson.Get(playerDatastr, "data").String()), &playerData)
	finished := 0
	for _, r := range playerData.Recruit {
		if r.State == 2 && time.Now().Unix() > int64(r.FinishTs) {
			finished += 1
		}
	}
	playerData.RecruitTotal = len(playerData.Recruit)
	playerData.RecruitFinished = finished

	playerData.Building.TiredCharsCount = len(playerData.Building.TiredChars)

	tradingsCurrent := 0
	tradingsTotal := 0
	for _, t := range playerData.Building.Tradings {
		tradingsTotal += t.StockLimit
		tradingsCurrent += t.StockLimit - len(t.Stock)
	}
	playerData.Building.TradingsCurrent = tradingsCurrent
	playerData.Building.TradingsTotal = tradingsTotal

	manufacturesCurrent := 0
	manufacturesTotal := 0
	for _, m := range playerData.Building.Manufactures {
		manufacturesCurrent += m.Complete
		manufacturesTotal += m.Capacity / Weight[m.FormulaID]
	}
	playerData.Building.ManufacturesCurrent = manufacturesCurrent
	playerData.Building.ManufacturesTotal = manufacturesTotal

	return playerData, nil
}

func getPlayerInfo(uid string, skland AccountSkland) (*PlayerData, error) {
	req := SKR().SetQueryParams(gh.MS{"uid": uid})
	return SklandRequest[*PlayerData](req, "GET", "/api/v1/game/player/info", skland)
}

func getPlayerInfoStr(uid string, skland AccountSkland) (string, error) {
	req := SKR().SetQueryParams(gh.MS{"uid": uid})
	return SklandRequestPlayerData(req, "GET", "/api/v1/game/player/info", skland)
}

func GetPlayerStatistic(uid string, account Account) (*PlayerStatistic, error) {
	var playerStatistic *PlayerStatistic
	account, err := RefreshToken(uid, account)
	if err != nil {
		log.Println(err.Error())
		return playerStatistic, err
	}
	playerStatistic, err = getPlayerStatistic(uid, account.Skland)
	if err != nil {
		return playerStatistic, err
	}

	return playerStatistic, nil
}

func getPlayerStatistic(uid string, skland AccountSkland) (*PlayerStatistic, error) {
	req := SKR().SetQueryParams(gh.MS{"uid": uid})
	return SklandRequest[*PlayerStatistic](req, "GET", "/api/v1/game/player/statistic", skland)
}
