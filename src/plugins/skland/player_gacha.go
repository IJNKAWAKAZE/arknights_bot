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
func GetPlayerGacha(token string) ([]Char, error) {
	var chars []Char
	err := checkToken(token)
	if err != nil {
		log.Println(err)
		return chars, err
	}
	res, err := getPlayerGacha(token, "1")
	if err != nil {
		log.Println(err)
		return chars, err
	}

	totalPage := int(math.Ceil(gjson.Get(res, "data.pagination.total").Float() / 10))

	for i := 1; i <= totalPage; i++ {
		res, err = getPlayerGacha(token, strconv.Itoa(i))
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

func getPlayerGacha(token, page string) (string, error) {
	req := SKR().SetQueryParams(gh.MS{"token": token, "page": page})
	res, err := HypergryphGacheRequest(req, "GET", "https://ak.hypergryph.com/user/api/inquiry/gacha")
	if err != nil {
		log.Println(err)
		return "", err
	}
	return res, nil
}
