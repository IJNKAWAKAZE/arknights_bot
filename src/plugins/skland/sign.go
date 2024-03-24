package skland

import (
	"fmt"
	"github.com/starudream/go-lib/core/v2/gh"
	"strconv"
	"strings"
	"time"
)

var signGameCodeByAppCode = map[string]string{
	GameAppCodeArknights: GameCodeArknights,
}

type SignGameData struct {
	Ts     string         `json:"ts"`
	Awards SignGameAwards `json:"awards"`
}

type SignGameRecord struct {
	GameId        string
	GameName      string
	PlayerName    string
	PlayerUid     string
	PlayerChannel string
	HasSigned     bool
	Award         string
}

type ListAttendanceData struct {
	CurrentTs       string                  `json:"currentTs"`
	Calendar        []*Calendar             `json:"calendar"`
	Records         CalendarRecords         `json:"records"`
	ResourceInfoMap map[string]*SignGameRes `json:"resourceInfoMap"`
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

type CalendarRecord struct {
	Ts         string `json:"ts"`
	ResourceId string `json:"resourceId"`
	Type       string `json:"type"`
	Count      int    `json:"count"`
}

type Calendar struct {
	ResourceId string `json:"resourceId"`
	Type       string `json:"type"`
	Count      int    `json:"count"`
	Available  bool   `json:"available"`
	Done       bool   `json:"done"`
}

type CalendarRecords []*CalendarRecord
type SignGameAwards []*SignGameAward

func SignGamePlayer(player *Player, account Account) (record SignGameRecord, err error) {
	account, err = RefreshToken(account)
	if err != nil {
		return
	}
	record.GameName = "明日方舟"
	record.PlayerName = player.NickName
	record.PlayerUid = player.Uid
	record.PlayerChannel = player.ChannelName

	gameId := signGameCodeByAppCode["arknights"]

	record.GameId = gameId

	/*list, err := listSignGame(gameId, player.Uid, account.Skland)
	if err != nil {
		err = fmt.Errorf("list sign game error: %w", err)
		return
	}

	today := list.Records.today()
	if len(today) > 0 {
		record.HasSigned = true
		record.Award = today.shortString(list.ResourceInfoMap)
		return
	}*/

	signGameData, err := signGame(gameId, player.Uid, account.Skland)
	if err != nil {
		if err.Error() == "[skland] response status: 403 Forbidden, error: code: 10001, message: 请勿重复签到！" {
			err = nil
			record.HasSigned = true
		} else {
			err = fmt.Errorf("sign game error: %w", err)
			return
		}
	} else {
		record.Award = signGameData.Awards.shortString()
	}

	return
}

// 签到信息
func listSignGame(gid, uid string, skland AccountSkland) (*ListAttendanceData, error) {
	req := SKR().SetQueryParams(gh.MS{"gameId": gid, "uid": uid})
	return SklandRequest[*ListAttendanceData](req, "GET", "/api/v1/game/attendance", skland)
}

// 签到
func signGame(gid, uid string, skland AccountSkland) (*SignGameData, error) {
	req := SKR().SetBody(gh.M{"gameId": gid, "uid": uid})
	return SklandRequest[*SignGameData](req, "POST", "/api/v1/game/attendance", skland)
}

func (v1 CalendarRecords) today() (v2 CalendarRecords) {
	now := time.Now()
	zero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	zeroTs := strconv.FormatInt(zero.Unix(), 10)
	for _, r := range v1 {
		if r.Ts == zeroTs {
			v2 = append(v2, r)
		}
	}
	return
}

func (v1 CalendarRecords) shortString(m map[string]*SignGameRes) string {
	v2 := make([]string, len(v1))
	for i, v := range v1 {
		v2[i] = m[v.ResourceId].Name + "*" + strconv.Itoa(v.Count)
	}
	return strings.Join(v2, ", ")
}

func (t SignGameAwards) shortString() string {
	v := make([]string, len(t))
	for i, a := range t {
		v[i] = a.Resource.Name + "*" + strconv.Itoa(a.Count)
	}
	return strings.Join(v, ", ")
}
