package utils

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"net/http"
	"strings"
)

type Operator struct {
	Name       string   `json:"name"`
	Profession string   `json:"profession"`
	Rarity     int      `json:"rarity"`
	ThumbURL   string   `json:"thumbURL"`
	Skins      []string `json:"skins"`
}

func GetOperators() []gjson.Result {
	operatorsJson := RedisGet("data_source")
	return gjson.Parse(operatorsJson).Array()
}

func GetOperatorList() []Operator {
	var operatorList []Operator
	response, _ := http.Get(viper.GetString("api.wiki_bili"))
	doc, _ := goquery.NewDocumentFromReader(response.Body)
	doc.Find(".floatnone").Each(func(i int, selection *goquery.Selection) {
		if selection.Nodes[0].FirstChild.FirstChild.Attr != nil {
			var operator Operator
			operator.Name = selection.Nodes[0].FirstChild.Attr[1].Val
			operator.ThumbURL = selection.Nodes[0].FirstChild.FirstChild.Attr[1].Val
			operatorList = append(operatorList, operator)
		}
	})
	defer response.Body.Close()
	return operatorList
}

func GetOperatorByName(name string) Operator {
	var operator Operator
	for _, op := range GetOperators() {
		if op.Get("name").String() == name {
			json.Unmarshal([]byte(op.String()), &operator)
		}
	}
	return operator
}

func GetOperatorsByName(name string) []Operator {
	var operatorList []Operator
	for _, operator := range GetOperatorList() {
		if strings.Contains(strings.ToLower(operator.Name), strings.ToLower(name)) {
			operatorList = append(operatorList, operator)
		}
	}
	return operatorList
}
