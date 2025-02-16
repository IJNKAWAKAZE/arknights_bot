package skland

import (
	"fmt"
	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/tidwall/gjson"
	"log"
	"strconv"
)

type PoolInfo struct {
	PoolName string `json:"poolName"`
	PoolId   string `json:"poolId"`
}

type Char struct {
	PoolName  string `json:"poolName"`
	PoolOrder int    `json:"poolOrder"`
	Name      string `json:"name"`
	IsNew     bool   `json:"isNew"`
	Rarity    int64  `json:"rarity"`
	Ts        int64  `json:"ts"`
}

// GetPlayerGacha 抽卡记录
func GetPlayerGacha(token, channelId, uid string) ([]Char, error) {
	var chars []Char
	if channelId == "1" {
		_, err := CheckToken(token)
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

	u8Token, err := LoginHypergryph(token, uid)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("登录失败")
	}

	//  获取卡池信息
	pools, err := getPoolList(token, u8Token, uid)
	if err != nil {
		return nil, err
	}

	for _, pool := range pools {
		//  获取卡池抽卡记录
		res, err := getPlayerGacha(token, u8Token, uid, pool.PoolId, "0", "0")
		if err != nil {
			log.Println(err)
			return chars, err
		}

		// 最后一条记录时间戳
		lastTs := ""
		lastPos := ""
		for _, d := range res.Get("data.list").Array() {
			ts := d.Get("gachaTs").Int()
			pos := int(d.Get("pos").Int())
			char := Char{
				PoolName:  d.Get("poolName").String(),
				PoolOrder: pos,
				Name:      d.Get("charName").String(),
				IsNew:     d.Get("isNew").Bool(),
				Rarity:    d.Get("rarity").Int(),
				Ts:        ts,
			}
			lastTs = strconv.FormatInt(ts, 10)
			lastPos = strconv.Itoa(pos)
			chars = append(chars, char)
		}

		// 是否有更多记录
		hasMore := res.Get("data.hasMore").Bool()

		for hasMore {
			// 获取下一页抽卡记录
			res, err = getPlayerGacha(token, u8Token, uid, pool.PoolId, lastTs, lastPos)
			if err != nil {
				log.Println(err)
				return chars, err
			}
			for _, d := range res.Get("data.list").Array() {
				ts := d.Get("gachaTs").Int()
				pos := int(d.Get("pos").Int())
				char := Char{
					PoolName:  d.Get("poolName").String(),
					PoolOrder: pos,
					Name:      d.Get("charName").String(),
					IsNew:     d.Get("isNew").Bool(),
					Rarity:    d.Get("rarity").Int(),
					Ts:        ts,
				}
				lastTs = strconv.FormatInt(ts, 10)
				lastPos = strconv.Itoa(pos)
				chars = append(chars, char)
			}
			hasMore = res.Get("data.hasMore").Bool()
		}
	}
	return chars, err
}

func getPoolList(token, u8Token, uid string) ([]PoolInfo, error) {
	var pools []PoolInfo
	req := HR().SetQueryParams(gh.MS{"uid": uid})
	req.SetHeader("X-Account-Token", token).SetHeader("X-Role-Token", u8Token)
	res, err := HypergryphAKRequest(req, "GET", "/user/api/inquiry/gacha/cate")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for _, p := range gjson.Parse(res).Get("data").Array() {
		pool := PoolInfo{
			PoolId:   p.Get("id").String(),
			PoolName: p.Get("name").String(),
		}
		pools = append(pools, pool)
	}
	return pools, nil
}

func getPlayerGacha(token, u8Token, uid, category, lastTs, pos string) (gjson.Result, error) {
	req := HR().SetQueryParams(gh.MS{"uid": uid, "category": category, "size": "199", "gachaTs": lastTs, "pos": pos})
	req.SetHeader("X-Account-Token", token).SetHeader("X-Role-Token", u8Token)
	res, err := HypergryphAKRequest(req, "GET", "/user/api/inquiry/gacha/history")
	if err != nil {
		log.Println(err)
		return gjson.Result{}, err
	}
	return gjson.Parse(res), nil
}
