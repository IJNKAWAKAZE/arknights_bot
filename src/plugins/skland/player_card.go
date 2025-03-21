package skland

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"log"
)

type PlayerCards struct {
	List []struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Icon    string `json:"icon"`
		BgURL   string `json:"bgUrl"`
		Privacy struct {
			CardOn         bool `json:"cardOn"`
			DetailOn       bool `json:"detailOn"`
			GameRelationOn bool `json:"gameRelationOn"`
		} `json:"privacy"`
		Link      string `json:"link"`
		Arknights struct {
			UID    string `json:"uid"`
			Name   string `json:"name"`
			Level  int    `json:"level"`
			Avatar struct {
				Type string `json:"type"`
				ID   string `json:"id"`
				URL  string `json:"url"`
			} `json:"avatar"`
			RegisterTs        int    `json:"registerTs"`
			MainStageProgress string `json:"mainStageProgress"`
			Secretary         struct {
				CharID string `json:"charId"`
				SkinID string `json:"skinId"`
			} `json:"secretary"`
			Resume          string `json:"resume"`
			SubscriptionEnd int    `json:"subscriptionEnd"`
			Ap              struct {
				Current              int `json:"current"`
				Max                  int `json:"max"`
				LastApAddTime        int `json:"lastApAddTime"`
				CompleteRecoveryTime int `json:"completeRecoveryTime"`
			} `json:"ap"`
			StoreTs      int `json:"storeTs"`
			LastOnlineTs int `json:"lastOnlineTs"`
			CharCnt      int `json:"charCnt"`
			FurnitureCnt int `json:"furnitureCnt"`
			SkinCnt      int `json:"skinCnt"`
			Exp          struct {
				Current int `json:"current"`
				Max     int `json:"max"`
			} `json:"exp"`
		} `json:"arknights"`
		IconBorderColor string `json:"iconBorderColor"`
		GameChar        string `json:"gameChar"`
		Decoration      struct {
			ID           int    `json:"id"`
			URL          string `json:"url"`
			Kind         int    `json:"kind"`
			ResourceKind int    `json:"resourceKind"`
			TopColor     string `json:"topColor"`
			TextColor    string `json:"textColor"`
		} `json:"decoration"`
	} `json:"list"`
}

func GetPlayerCards(account Account) (*PlayerCards, error) {
	var playerCards *PlayerCards
	account, err := RefreshToken(account)
	if err != nil {
		log.Println(err.Error())
		return playerCards, err
	}
	playerCardsStr, err := getPlayerCardsStr(account.Skland)
	if err != nil {
		return playerCards, err
	}

	json.Unmarshal([]byte(gjson.Get(playerCardsStr, "data").String()), &playerCards)
	return playerCards, nil
}

func getPlayerCardsStr(skland AccountSkland) (string, error) {
	return SklandRequestPlayerData(SKR(), "GET", "/api/v1/game/cards", skland)
}
