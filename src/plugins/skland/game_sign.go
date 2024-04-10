package skland

import (
	"fmt"
	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/starudream/go-lib/resty/v2"
	"strconv"
	"strings"
)

type SignGameData struct {
	Ts     string         `json:"ts"`
	Awards SignGameAwards `json:"awards"`
}

type SignGameAward struct {
	Type     string       `json:"type"`
	Count    int          `json:"count"`
	Resource *SignGameRes `json:"resource"`
}

type SignGameRes struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Rarity int    `json:"rarity"`
}

type SignGameAwards []*SignGameAward

func SignGamePlayer(uid string, account Account) (award string, hasSigned bool, err error) {
	account, err = RefreshToken(account)
	if err != nil {
		return
	}
	signGameData, err := signGame("1", uid, account.Skland)
	if err != nil {
		e, ok1 := resty.AsRespErr(err)
		if ok1 {
			t, ok2 := e.Response.Error().(*SKBaseResp[interface{}])
			if ok2 && t.Message == "请勿重复签到！" {
				err = nil
				hasSigned = true
			}
		} else {
			err = fmt.Errorf("sign game error: %w", err)
			return
		}
	} else {
		award = signGameData.Awards.shortString()
	}
	return
}

// 签到
func signGame(gid, uid string, skland AccountSkland) (*SignGameData, error) {
	req := SKR().SetBody(gh.M{"gameId": gid, "uid": uid})
	return SklandRequest[*SignGameData](req, "POST", "/api/v1/game/attendance", skland)
}

func (t SignGameAwards) shortString() string {
	v := make([]string, len(t))
	for i, a := range t {
		v[i] = a.Resource.Name + "*" + strconv.Itoa(a.Count)
	}
	return strings.Join(v, ", ")
}
