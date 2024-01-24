package utils

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"net/http"
	"strings"
)

type Operator struct {
	Name     string `json:"name"`
	Painting string `json:"painting"`
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
		var operator Operator
		operator.Name = selection.Nodes[0].FirstChild.Attr[1].Val
		operator.Painting = selection.Nodes[0].FirstChild.FirstChild.Attr[1].Val
		operatorList = append(operatorList, operator)
	})
	defer response.Body.Close()
	return operatorList
}

func GetOperatorByName(name string) Operator {
	var operator Operator
	for _, op := range GetOperators() {
		if op.Get("name").String() == name {
			operator.Name = op.Get("name").String()
			operator.Painting = op.Get("skins").Array()[0].String()
		}
	}
	return operator
}

func GetOperatorsByName(name string) []Operator {
	var operatorList []Operator
	for _, operator := range GetOperatorList() {
		if strings.Contains(operator.Name, name) {
			operatorList = append(operatorList, operator)
		}
	}
	return operatorList
}
