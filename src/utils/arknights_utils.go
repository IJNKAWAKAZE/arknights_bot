package utils

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"strings"
)

type Operator struct {
	Name         string   `json:"name"`         // 名字
	NameEn       string   `json:"nameEn"`       // 英文名
	NameJa       string   `json:"nameJp"`       // 日文名
	Code         string   `json:"code"`         // 编号
	Race         string   `json:"race"`         // 种族
	Profession   string   `json:"profession"`   // 职业
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
	Nation       string   `json:"nation"`       // 阵营
	ObtainMethod string   `json:"obtainMethod"` // 获取方式
	Tags         string   `json:"tags"`         // 标签
}

var operatorMap = make(map[string]Operator)

func GetOperators() []gjson.Result {
	operatorsJson := RedisGet("data_source")
	return gjson.Parse(operatorsJson).Array()
}

func GetOperatorByName(name string) Operator {
	if len(operatorMap) == 0 {
		for _, op := range GetOperators() {
			var operator Operator
			opName := op.Get("name").String()
			json.Unmarshal([]byte(op.String()), &operator)
			operatorMap[opName] = operator
		}
	}
	return operatorMap[name]
}

func GetOperatorsByName(name string) []Operator {
	var operatorList []Operator
	for _, op := range GetOperators() {
		if strings.Contains(strings.ToLower(op.Get("name").String()), strings.ToLower(name)) {
			var operator Operator
			json.Unmarshal([]byte(op.String()), &operator)
			operatorList = append(operatorList, operator)
		}
	}
	return operatorList
}
