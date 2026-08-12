package skland

import (
	"fmt"
	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/tidwall/gjson"
	"log"
	"strconv"
)

// GetPlayerGacha 抽卡记录
func GetPlayerGacha(token, uid string) ([]Char, error) {
	var chars []Char
	if _, err := checkToken(token, false); err != nil {
		log.Println(err)
		return chars, err
	}

	u8Token, err := loginHypergryph(token, uid)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("登录失败")
	}

	// 获取卡池信息
	pools, err := getPoolList(token, u8Token, uid)
	if err != nil {
		return nil, err
	}

	for _, pool := range pools {
		lastTs, lastPos := "0", "0"
		for {
			// 获取卡池抽卡记录
			res, err := getPlayerGacha(token, u8Token, uid, pool.PoolId, lastTs, lastPos)
			if err != nil {
				log.Println(err)
				return chars, err
			}
			for _, d := range res.Get("data.list").Array() {
				ts := d.Get("gachaTs").Int()
				pos := int(d.Get("pos").Int())
				chars = append(chars, Char{
					PoolName:  d.Get("poolName").String(),
					PoolOrder: pos,
					Name:      d.Get("charName").String(),
					IsNew:     d.Get("isNew").Bool(),
					Rarity:    d.Get("rarity").Int(),
					Ts:        ts,
				})
				lastTs = strconv.FormatInt(ts, 10)
				lastPos = strconv.Itoa(pos)
			}
			// 是否有更多记录
			if !res.Get("data.hasMore").Bool() {
				break
			}
		}
	}
	return chars, nil
}

// getPoolList 获取卡池信息
func getPoolList(token, u8Token, uid string) ([]PoolInfo, error) {
	var pools []PoolInfo
	req := HR().SetQueryParams(gh.MS{"uid": uid})
	req.SetHeader("X-Account-Token", token).SetHeader("X-Role-Token", u8Token)
	res, err := hgRawRequest(req, "GET", "/user/api/inquiry/gacha/cate", hypergryphAKAddr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for _, p := range gjson.Parse(res).Get("data").Array() {
		pools = append(pools, PoolInfo{
			PoolId:   p.Get("id").String(),
			PoolName: p.Get("name").String(),
		})
	}
	return pools, nil
}

// getPlayerGacha 获取卡池抽卡记录
func getPlayerGacha(token, u8Token, uid, category, lastTs, pos string) (gjson.Result, error) {
	req := HR().SetQueryParams(gh.MS{"uid": uid, "category": category, "size": "199", "gachaTs": lastTs, "pos": pos})
	req.SetHeader("X-Account-Token", token).SetHeader("X-Role-Token", u8Token)
	res, err := hgRawRequest(req, "GET", "/user/api/inquiry/gacha/history", hypergryphAKAddr)
	if err != nil {
		return gjson.Result{}, err
	}
	return gjson.Parse(res), nil
}
