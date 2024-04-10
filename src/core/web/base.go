package web

import (
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
)

type PlayerBase struct {
	Labor        Labor         `json:"labor"`        // 无人机
	Control      Control       `json:"control"`      // 控制中枢
	Tradings     []Trading     `json:"tradings"`     // 贸易站
	Manufactures []Manufacture `json:"manufactures"` // 制造站
	Powers       []Power       `json:"powers"`       // 发电站
	Meeting      Meeting       `json:"meeting"`      // 会客室
	Hire         Hire          `json:"hire"`         // 办公室
	Training     Training      `json:"training"`     // 训练室
	Dormitories  []Dormitory   `json:"dormitories"`  // 宿舍
}

type Labor struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

type Control struct {
	Level int        `json:"level"`
	Chars []BaseChar `json:"chars"`
}

type Trading struct {
	Level    int        `json:"level"`
	Chars    []BaseChar `json:"chars"`
	Current  int        `json:"current"`
	Total    int        `json:"total"`
	Strategy string     `json:"strategy"`
}

type Manufacture struct {
	Level   int        `json:"level"`
	Chars   []BaseChar `json:"chars"`
	Current int        `json:"current"`
	Total   int        `json:"total"`
	Item    string     `json:"item"`
	Speed   string     `json:"speed"`
}

type Power struct {
	Level int        `json:"level"`
	Chars []BaseChar `json:"chars"`
	Power int        `json:"power"`
}

type Meeting struct {
	Level   int        `json:"level"`
	Chars   []BaseChar `json:"chars"`
	Board   []int      `json:"board"`
	Sharing bool       `json:"sharing"`
}

type Hire struct {
	Level        int        `json:"level"`
	Chars        []BaseChar `json:"chars"`
	RefreshCount int        `json:"refreshCount"`
}

type Training struct {
	Level           int        `json:"level"`
	Chars           []BaseChar `json:"chars"`
	Skill           string     `json:"skill"`
	SpecializeLevel int        `json:"specializeLevel"`
}

type Dormitory struct {
	Level   int        `json:"level"`
	Chars   []BaseChar `json:"chars"`
	Comfort int        `json:"comfort"`
}

type BaseChar struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	AP     int    `json:"AP"`
}

var Weight map[string]int

var Board map[string]int

var Item map[string]string

func init() {
	weight := make(map[string]int)
	weight["1"] = 2
	weight["2"] = 3
	weight["3"] = 5
	weight["4"] = 2
	weight["5"] = 5
	weight["6"] = 5
	weight["7"] = 5
	weight["8"] = 5
	weight["9"] = 5
	weight["10"] = 5
	weight["11"] = 5
	weight["12"] = 5
	weight["13"] = 3
	weight["14"] = 3
	Weight = weight

	item := make(map[string]string)
	item["1"] = "基础作战记录"
	item["2"] = "初级作战记录"
	item["3"] = "中级作战记录"
	item["4"] = "赤金"
	item["5"] = "先锋双芯片"
	item["6"] = "近卫双芯片"
	item["7"] = "重装双芯片"
	item["8"] = "狙击双芯片"
	item["9"] = "术师双芯片"
	item["10"] = "医疗双芯片"
	item["11"] = "辅助双芯片"
	item["12"] = "特种双芯片"
	item["13"] = "源石碎片"
	item["14"] = "源石碎片"
	Item = item

	board := make(map[string]int)
	board["RHINE"] = 1
	board["PENGUIN"] = 2
	board["BLACKSTEEL"] = 3
	board["URSUS"] = 4
	board["GLASGOW"] = 5
	board["KJERAG"] = 6
	board["RHODES"] = 7
	Board = board
}

func Base(r *gin.Engine) {
	r.GET("/base", func(c *gin.Context) {
		r.LoadHTMLFiles("./template/Base.tmpl")
		var playerBase PlayerBase
		var userAccount account.UserAccount
		var skAccount skland.Account
		userId, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
		uid := c.Query("uid")
		sklandId := c.Query("sklandId")
		utils.GetAccountByUserIdAndSklandId(userId, sklandId).Scan(&userAccount)
		skAccount.Hypergryph.Token = userAccount.HypergryphToken
		skAccount.Skland.Token = userAccount.SklandToken
		skAccount.Skland.Cred = userAccount.SklandCred
		playerData, skAccount, err := skland.GetPlayerInfo(uid, skAccount)
		if err != nil {
			log.Println(err)
			return
		}

		building := playerData.Building
		charMap := playerData.CharInfoMap

		// 无人机
		var labor Labor
		labor.Current = building.Labor.Value
		labor.Total = building.Labor.MaxValue
		playerBase.Labor = labor

		// 控制中枢
		var control Control
		control.Level = building.Control.Level
		for _, char := range building.Control.Chars {
			baseChar := BaseChar{
				Name:   charMap[char.CharID].Name,
				Avatar: getCharSkinID(char.CharID, playerData),
				AP:     getAp(char.Ap),
			}
			control.Chars = append(control.Chars, baseChar)
		}
		playerBase.Control = control

		// 贸易站
		var Tradings []Trading
		for _, t := range building.Tradings {
			var trading Trading
			trading.Level = t.Level
			for _, char := range t.Chars {
				baseChar := BaseChar{
					Name:   charMap[char.CharID].Name,
					Avatar: getCharSkinID(char.CharID, playerData),
					AP:     getAp(char.Ap),
				}
				trading.Chars = append(trading.Chars, baseChar)
			}
			trading.Current = len(t.Stock)
			trading.Total = t.StockLimit
			if t.Strategy == "O_GOLD" {
				trading.Strategy = "贵金属订单"
			} else {
				trading.Strategy = "源石订单"
			}
			Tradings = append(Tradings, trading)
		}
		playerBase.Tradings = Tradings

		// 制造站
		var manufactures []Manufacture
		for _, m := range building.Manufactures {
			var manufacture Manufacture
			manufacture.Level = m.Level
			for _, char := range m.Chars {
				baseChar := BaseChar{
					Name:   charMap[char.CharID].Name,
					Avatar: getCharSkinID(char.CharID, playerData),
					AP:     getAp(char.Ap),
				}
				manufacture.Chars = append(manufacture.Chars, baseChar)
			}
			manufacture.Current = m.Complete
			manufacture.Total = m.Capacity / Weight[m.FormulaID]
			manufacture.Item = Item[m.FormulaID]
			manufacture.Speed = fmt.Sprintf("%d%s", m.Speed*100, "%")
			manufactures = append(manufactures, manufacture)
		}
		playerBase.Manufactures = manufactures

		// 发电站
		var powers []Power
		for _, p := range building.Powers {
			var power Power
			power.Level = p.Level
			for _, char := range p.Chars {
				baseChar := BaseChar{
					Name:   charMap[char.CharID].Name,
					Avatar: getCharSkinID(char.CharID, playerData),
					AP:     getAp(char.Ap),
				}
				power.Chars = append(power.Chars, baseChar)
			}
			power.Power = int(math.Pow(2, float64(p.Level-1))*60 + (math.Pow(2, float64(p.Level-1))-1)*10)
			powers = append(powers, power)
		}
		playerBase.Powers = powers

		// 会客室
		var meeting Meeting
		meeting.Level = building.Meeting.Level
		for _, char := range building.Meeting.Chars {
			baseChar := BaseChar{
				Name:   charMap[char.CharID].Name,
				Avatar: getCharSkinID(char.CharID, playerData),
				AP:     getAp(char.Ap),
			}
			meeting.Chars = append(meeting.Chars, baseChar)
		}
		for _, board := range building.Meeting.Clue.Board {
			meeting.Board = append(meeting.Board, Board[board])
		}

		sort.Slice(meeting.Board, func(i, j int) bool {
			return meeting.Board[i] < meeting.Board[j]
		})

		meeting.Sharing = building.Meeting.Clue.Sharing
		playerBase.Meeting = meeting

		// 办公室
		var hire Hire
		hire.Level = building.Hire.Level
		for _, char := range building.Hire.Chars {
			baseChar := BaseChar{
				Name:   charMap[char.CharID].Name,
				Avatar: getCharSkinID(char.CharID, playerData),
				AP:     getAp(char.Ap),
			}
			hire.Chars = append(hire.Chars, baseChar)
		}
		hire.RefreshCount = building.Hire.RefreshCount
		playerBase.Hire = hire

		// 训练室
		var training Training
		training.Level = building.Training.Level
		if building.Training.Trainee.CharID != "" {
			training.Chars = append(training.Chars, BaseChar{
				Name:   charMap[building.Training.Trainee.CharID].Name,
				Avatar: getCharSkinID(building.Training.Trainee.CharID, playerData),
				AP:     getAp(building.Training.Trainee.Ap),
			})
			training.Skill, training.SpecializeLevel = getCharSkillID(building.Training.Trainee.CharID, playerData, building.Training.Trainee.TargetSkill)
		}
		if building.Training.Trainer.CharID != "" {
			training.Chars = append(training.Chars, BaseChar{
				Name:   charMap[building.Training.Trainer.CharID].Name,
				Avatar: getCharSkinID(building.Training.Trainer.CharID, playerData),
				AP:     getAp(building.Training.Trainer.Ap),
			})
		}

		playerBase.Training = training

		// 宿舍
		var dormitories []Dormitory
		for _, d := range building.Dormitories {
			var dormitory Dormitory
			dormitory.Level = d.Level
			for _, char := range d.Chars {
				baseChar := BaseChar{
					Name:   charMap[char.CharID].Name,
					Avatar: getCharSkinID(char.CharID, playerData),
					AP:     getAp(char.Ap),
				}
				dormitory.Chars = append(dormitory.Chars, baseChar)
			}
			dormitory.Comfort = d.Comfort
			dormitories = append(dormitories, dormitory)
		}

		playerBase.Dormitories = dormitories
		c.HTML(http.StatusOK, "Base.tmpl", playerBase)
	})
}

func getAp(ap int) int {
	return int(math.Ceil(float64(ap) / float64(86400)))
}

func getCharSkinID(charId string, data *skland.PlayerData) string {
	for _, char := range data.Chars {
		if char.CharID == charId {
			return char.SkinID
		}
	}
	return ""
}

func getCharSkillID(charId string, data *skland.PlayerData, target int) (string, int) {
	for _, char := range data.Chars {
		if char.CharID == charId && target != -1 {
			return char.Skills[target].ID, char.Skills[target].SpecializeLevel + 7
		}
	}
	return "", 7
}
