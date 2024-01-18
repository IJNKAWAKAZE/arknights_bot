package utils

import (
	"github.com/PuerkitoBio/goquery"
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
	response, _ := http.Get("https://wiki.biligame.com/arknights/%E5%B9%B2%E5%91%98%E6%95%B0%E6%8D%AE%E8%A1%A8")
	doc, _ := goquery.NewDocumentFromReader(response.Body)
	doc.Find(".floatnone").Each(func(i int, selection *goquery.Selection) {
		var operator Operator
		operator.Name = selection.Nodes[0].FirstChild.Attr[1].Val
		operator.Painting = selection.Nodes[0].FirstChild.FirstChild.Attr[1].Val
		operatorList = append(operatorList, operator)
	})
	return operatorList
}

func GetOperatorByName(name string) Operator {
	var operator Operator
	for _, op := range GetOperators() {
		if op.Get("name").String() == name {
			operator.Name = op.Get("name").String()
			operator.Painting = op.Get("painting").String()
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
