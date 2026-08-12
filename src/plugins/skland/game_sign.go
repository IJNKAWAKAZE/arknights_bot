package skland

import (
	"fmt"
	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/starudream/go-lib/resty/v2"
	"strconv"
	"strings"
)

// SignGamePlayer 森空岛签到
func SignGamePlayer(uid string, account Account, serverName string) (award string, hasSigned bool, err error) {
	account, err = RefreshToken(account, serverName)
	if err != nil {
		return
	}
	signGameData, err := signGame("1", uid, serverName == serverGlobal, account.Skland)
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

// signGame 签到
func signGame(gid, uid string, isGlobal bool, skland AccountSkland) (*SignGameData, error) {
	req := SKR()
	if isGlobal {
		req.SetHeader("sk-language", "zh_Hans")
	}
	req.SetBody(gh.M{"gameId": gid, "uid": uid})
	return skRequest[*SignGameData](req, "POST", "/api/v1/game/attendance", isGlobal, skland)
}

func (t SignGameAwards) shortString() string {
	v := make([]string, len(t))
	for i, a := range t {
		if a.Resource != nil {
			v[i] = a.Resource.Name + "*" + strconv.Itoa(a.Count)
		}
	}
	return strings.Join(v, ", ")
}
