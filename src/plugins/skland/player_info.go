package skland

import (
	bot "arknights_bot/config"
	"encoding/json"
	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/tidwall/gjson"
	"log"
)

func GetPlayerInfo(uid string, account Account, serverName string) (*PlayerData, Account, error) {
	var playerData *PlayerData
	account, err := RefreshToken(account, serverName)
	if err != nil {
		log.Println(err.Error())
		return playerData, account, err
	}
	var playerDatastr string
	if serverName == "国服" {
		playerDatastr, err = getPlayerInfoStr(uid, account.Skland)
	} else if serverName == "国际服" {
		playerDatastr, err = iGetPlayerInfoStr(uid, account.Skland)
	}
	if err != nil {
		return playerData, account, err
	}

	json.Unmarshal([]byte(gjson.Get(playerDatastr, "data").String()), &playerData)
	bot.DBEngine.Exec("update user_player set player_name = ? where uid = ?", playerData.Status.Name, uid)
	return playerData, account, nil
}

func getPlayerInfoStr(uid string, skland AccountSkland) (string, error) {
	req := SKR().SetQueryParams(gh.MS{"uid": uid})
	return SklandRequestPlayerData(req, "GET", "/api/v1/game/player/info", skland)
}

func iGetPlayerInfoStr(uid string, skland AccountSkland) (string, error) {
	req := SKR().SetQueryParams(gh.MS{"uid": uid})
	return SkportRequestPlayerData(req, "GET", "/api/v1/game/player/info", skland)
}

func GetPlayerStatistic(uid string, account Account, serverName string) (*PlayerStatistic, Account, error) {
	var playerStatistic *PlayerStatistic
	account, err := RefreshToken(account, serverName)
	if err != nil {
		log.Println(err.Error())
		return playerStatistic, account, err
	}
	if serverName == "国服" {
		playerStatistic, err = getPlayerStatistic(uid, account.Skland)
	} else if serverName == "国际服" {
		playerStatistic, err = iGetPlayerStatistic(uid, account.Skland)
	}
	if err != nil {
		return playerStatistic, account, err
	}

	return playerStatistic, account, nil
}

func getPlayerStatistic(uid string, skland AccountSkland) (*PlayerStatistic, error) {
	req := SKR().SetQueryParams(gh.MS{"uid": uid})
	return SklandRequest[*PlayerStatistic](req, "GET", "/api/v1/game/player/statistic", skland)
}

func iGetPlayerStatistic(uid string, skland AccountSkland) (*PlayerStatistic, error) {
	req := SKR().SetQueryParams(gh.MS{"uid": uid})
	return SkportRequest[*PlayerStatistic](req, "GET", "/api/v1/game/player/statistic", skland)
}
