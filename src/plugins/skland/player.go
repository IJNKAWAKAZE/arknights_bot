package skland

import (
	"encoding/json"
	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/tidwall/gjson"
	"log"
)

// GetPlayerInfo 获取玩家信息
func GetPlayerInfo(uid string, account Account, serverName string) (*PlayerData, Account, error) {
	var playerData *PlayerData
	account, err := RefreshToken(account, serverName)
	if err != nil {
		log.Println(err.Error())
		return playerData, account, err
	}
	res, err := getPlayerDataStr(uid, account.Skland, "/api/v1/game/player/info", serverName == serverGlobal)
	if err != nil {
		return playerData, account, err
	}
	err = json.Unmarshal([]byte(gjson.Get(res, "data").String()), &playerData)
	return playerData, account, err
}

// GetPlayerStatistic 获取玩家概况
func GetPlayerStatistic(uid string, account Account, serverName string) (*PlayerStatistic, Account, error) {
	var playerStatistic *PlayerStatistic
	account, err := RefreshToken(account, serverName)
	if err != nil {
		log.Println(err.Error())
		return playerStatistic, account, err
	}
	req := SKR().SetQueryParams(gh.MS{"uid": uid})
	playerStatistic, err = skRequest[*PlayerStatistic](req, "GET", "/api/v1/game/player/statistic", serverName == serverGlobal, account.Skland)
	return playerStatistic, account, err
}

// GetPlayerCards 获取玩家名片
func GetPlayerCards(account Account, serverName string) (*PlayerCards, error) {
	var playerCards *PlayerCards
	account, err := RefreshToken(account, serverName)
	if err != nil {
		log.Println(err.Error())
		return playerCards, err
	}
	res, err := skRequestData(SKR(), "GET", "/api/v1/game/cards", serverName == serverGlobal, account.Skland)
	if err != nil {
		return playerCards, err
	}
	err = json.Unmarshal([]byte(gjson.Get(res, "data").String()), &playerCards)
	return playerCards, err
}

// GetPlayerCultivate 获取玩家养成数据
func GetPlayerCultivate(uid string, account Account, serverName string) (*PlayerCultivate, error) {
	var playerCultivate *PlayerCultivate
	account, err := RefreshToken(account, serverName)
	if err != nil {
		log.Println(err.Error())
		return playerCultivate, err
	}
	res, err := getPlayerDataStr(uid, account.Skland, "/api/v1/game/cultivate/player", serverName == serverGlobal)
	if err != nil {
		return playerCultivate, err
	}
	err = json.Unmarshal([]byte(gjson.Get(res, "data").String()), &playerCultivate)
	return playerCultivate, err
}

// getPlayerDataStr 获取森空岛玩家数据接口的原始响应
func getPlayerDataStr(uid string, skland AccountSkland, path string, isGlobal bool) (string, error) {
	req := SKR().SetQueryParams(gh.MS{"uid": uid})
	return skRequestData(req, "GET", path, isGlobal, skland)
}
