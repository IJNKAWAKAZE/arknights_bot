package skland

import (
	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/tidwall/gjson"
	"log"
	"math"
	"strconv"
)

type Char struct {
	PoolName  string `json:"poolName"`
	PoolOrder int    `json:"poolOrder"`
	Name      string `json:"name"`
	IsNew     bool   `json:"isNew"`
	Rarity    int64  `json:"rarity"`
	Ts        int64  `json:"ts"`
}

// GetPlayerGacha 抽卡记录
func GetPlayerGacha(token, channelId string) ([]Char, error) {
	var chars []Char
	if channelId == "1" {
		err := CheckToken(token)
		if err != nil {
			log.Println(err)
			return chars, err
		}
	} else if channelId == "2" {
		err := CheckBToken(token)
		if err != nil {
			log.Println(err)
			return chars, err
		}
	}
	res, err := getPlayerGacha(token, "1", channelId)
	if err != nil {
		log.Println(err)
		return chars, err
	}

	totalPage := int(math.Ceil(gjson.Get(res, "data.pagination.total").Float() / 10))

	for i := 1; i <= totalPage; i++ {
		res, err = getPlayerGacha(token, strconv.Itoa(i), channelId)
		if err != nil {
			break
		}
		for _, d := range gjson.Get(res, "data.list").Array() {
			poolName := d.Get("pool").String()
			ts := d.Get("ts").Int()
			order := 1
			for _, c := range d.Get("chars").Array() {
				char := Char{
					PoolName:  poolName,
					PoolOrder: order,
					Name:      c.Get("name").String(),
					IsNew:     c.Get("isNew").Bool(),
					Rarity:    c.Get("rarity").Int(),
					Ts:        ts,
				}
				order++
				chars = append(chars, char)
			}
		}
	}
	return chars, err
}

func getPlayerGacha(token, page, channelId string) (string, error) {
	req := SKR().SetQueryParams(gh.MS{"token": token, "page": page, "channelId": channelId})
	res, err := HypergryphAKRequest(req, "GET", "/user/api/inquiry/gacha")
	if err != nil {
		log.Println(err)
		return "", err
	}
	return res, nil
}
