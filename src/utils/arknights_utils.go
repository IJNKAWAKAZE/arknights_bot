package utils

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"strings"
)

type Operator struct {
	Name       string   `json:"name"`
	Profession string   `json:"profession"`
	Rarity     int      `json:"rarity"`
	Avatar     string   `json:"avatar"`
	ThumbURL   string   `json:"thumbURL"`
	Skins      []string `json:"skins"`
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
