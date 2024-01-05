package skland

import (
	"github.com/starudream/go-lib/core/v2/gh"
	"log"
)

type PlayerInfo struct {
}

func GetPlayerInfo(uid string, account Account) (PlayerInfo, error) {
	playerInfo := PlayerInfo{}
	account, err := RefreshToken(account)
	if err != nil {
		log.Println(err.Error())
		return playerInfo, err
	}
	_, err = getPlayerInfo(uid, account.Skland)
	if err != nil {
		return playerInfo, err
	}
	return playerInfo, nil
}

func getPlayerInfo(uid string, skland AccountSkland) (*PlayerInfo, error) {
	req := SKR().SetQueryParams(gh.MS{"uid": uid})
	return SklandRequest[*PlayerInfo](req, "GET", "/api/v1/game/player/info", skland)
}
