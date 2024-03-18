package utils

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Operator struct {
	Name         string   `json:"name"`         // 名字
	NameEn       string   `json:"nameEn"`       // 英文名
	NameJa       string   `json:"nameJp"`       // 日文名
	Code         string   `json:"code"`         // 编号
	Race         string   `json:"race"`         // 种族
	Profession   string   `json:"profession"`   // 职业
	ProfessionZH string   `json:"professionZH"` // 职业
	Rarity       int      `json:"rarity"`       // 稀有度
	Avatar       string   `json:"avatar"`       // 头像
	ThumbURL     string   `json:"thumbURL"`     // 半身像
	Skins        []string `json:"skins"`        // 皮肤
	HP           string   `json:"hp"`           // 生命值
	ATK          string   `json:"atk"`          // 攻击
	DEF          string   `json:"def"`          // 防御
	Res          string   `json:"res"`          // 法抗
	ReDeploy     string   `json:"reDeploy"`     // 再部署时间
	Cost         string   `json:"cost"`         // 费用
	Block        string   `json:"block"`        // 阻挡数
	Interval     string   `json:"interval"`     // 攻击间隔
	Sex          string   `json:"sex"`          // 性别
	Position     string   `json:"position"`     // 部署位
	Logo         string   `json:"logo"`         // 所属
	ObtainMethod string   `json:"obtainMethod"` // 获取方式
	Tags         string   `json:"tags"`         // 标签
}

type Stage struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Path string `json:"path"`
}

var OperatorMap = make(map[string]Operator)
var RecruitOperatorList []Operator

func GetOperators() []Operator {
	var operators []Operator
	operatorsJson := RedisGet("data_source")
	json.Unmarshal([]byte(operatorsJson), &operators)
	return operators
}

func GetOperatorByName(name string) Operator {
	if len(OperatorMap) == 0 {
		for _, op := range GetOperators() {
			OperatorMap[op.Name] = op
		}
	}
	return OperatorMap[name]
}

func GetOperatorsByName(name string) []Operator {
	var operatorList []Operator
	for _, op := range GetOperators() {
		if strings.Contains(strings.ToLower(op.Name), strings.ToLower(name)) {
			operatorList = append(operatorList, op)
		}
	}
	return operatorList
}

func GetRecruitOperatorList() []Operator {
	if len(RecruitOperatorList) == 0 {
		for _, op := range GetOperators() {
			if strings.Contains(op.ObtainMethod, "公开招募") {
				RecruitOperatorList = append(RecruitOperatorList, op)
			}
		}
	}
	return RecruitOperatorList
}

func GetEnemiesByName(name string) map[string]string {
	var enemyList = make(map[string]string)
	api := viper.GetString("api.enemy")
	response, _ := http.Get(api)
	e, _ := io.ReadAll(response.Body)
	defer response.Body.Close()
	enemyJson := gjson.ParseBytes(e)
	for _, en := range enemyJson.Array() {
		n := en.Get("name").String()
		if strings.Contains(strings.ToLower(n), strings.ToLower(name)) {
			paintingName := fmt.Sprintf("头像_敌人_%s.png", n)
			m := Md5(paintingName)
			path := "https://prts.wiki" + fmt.Sprintf("/images/%s/%s/", m[:1], m[:2])
			pic := path + url.PathEscape(paintingName)
			enemyList[n] = pic
		}
	}
	return enemyList
}
