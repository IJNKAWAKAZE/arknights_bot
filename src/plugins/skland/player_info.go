package skland

import (
	"encoding/json"
	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/tidwall/gjson"
	"log"
)

func GetPlayerInfo(uid string, account Account) (*PlayerData, Account, error) {
	var playerData *PlayerData
	account, err := RefreshToken(account)
	if err != nil {
		log.Println(err.Error())
		return playerData, account, err
	}
	playerDatastr, err := getPlayerInfoStr(uid, account.Skland)
	if err != nil {
		return playerData, account, err
	}

	json.Unmarshal([]byte(gjson.Get(playerDatastr, "data").String()), &playerData)

	return playerData, account, nil
}

func getPlayerInfo(uid string, skland AccountSkland) (*PlayerData, error) {
	req := SKR().SetQueryParams(gh.MS{"uid": uid})
	return SklandRequest[*PlayerData](req, "GET", "/api/v1/game/player/info", skland)
}

func getPlayerInfoStr(uid string, skland AccountSkland) (string, error) {
	req := SKR().SetQueryParams(gh.MS{"uid": uid})
	return SklandRequestPlayerData(req, "GET", "/api/v1/game/player/info", skland)
}

func GetPlayerStatistic(uid string, account Account) (*PlayerStatistic, Account, error) {
	var playerStatistic *PlayerStatistic
	account, err := RefreshToken(account)
	if err != nil {
		log.Println(err.Error())
		return playerStatistic, account, err
	}
	playerStatistic, err = getPlayerStatistic(uid, account.Skland)
	if err != nil {
		return playerStatistic, account, err
	}

	return playerStatistic, account, nil
}

func getPlayerStatistic(uid string, skland AccountSkland) (*PlayerStatistic, error) {
	req := SKR().SetQueryParams(gh.MS{"uid": uid})
	return SklandRequest[*PlayerStatistic](req, "GET", "/api/v1/game/player/statistic", skland)
}
