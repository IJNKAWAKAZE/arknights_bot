package skland

import (
	"github.com/starudream/go-lib/core/v2/gh"
	"log"
)

func GetPlayerInfo(uid string, account Account) (*PlayerData, error) {
	var playerData *PlayerData
	account, err := RefreshToken(uid, account)
	if err != nil {
		log.Println(err.Error())
		return playerData, err
	}
	playerData, err = getPlayerInfo(uid, account.Skland)
	if err != nil {
		return playerData, err
	}
	return playerData, nil
}

func getPlayerInfo(uid string, skland AccountSkland) (*PlayerData, error) {
	req := SKR().SetQueryParams(gh.MS{"uid": uid})
	return SklandRequest[*PlayerData](req, "GET", "/api/v1/game/player/info", skland)
}
