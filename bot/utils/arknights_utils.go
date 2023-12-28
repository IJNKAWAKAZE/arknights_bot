package utils

import (
	"arknights_bot/bot/config"
	"arknights_bot/bot/modules"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
	"github.com/tidwall/gjson"
	"log"
	"net/http"
	"strings"
)

func GetOperators() []gjson.Result {
	shipsJson := RedisGet("data_source")
	return gjson.Parse(shipsJson).Array()
}

// UpdateDataSource 更新数据源
func UpdateDataSource() {
	log.Println("开始更新数据源...")
	var operatorJson []modules.Verify
	var operatorList []string
	response, _ := http.Get(config.GetString("api.wiki") + "干员一览")
	doc, _ := goquery.NewDocumentFromReader(response.Body)
	doc.Find("#filter-data div").Each(func(i int, selection *goquery.Selection) {
		operatorList = append(operatorList, selection.Nodes[0].Attr[0].Val)
	})
	for _, name := range operatorList {
		var operator modules.Verify
		operator.Name = name
		pw, _ := playwright.Run()
		browser, _ := pw.Chromium.Launch()
		page, _ := browser.NewPage()
		page.Goto(config.GetString("api.wiki")+name, playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
		})
		locator, _ := page.Locator("#charimg")
		imgs, _ := locator.InnerHTML()
		imgHtml, _ := goquery.NewDocumentFromReader(strings.NewReader(imgs))
		imgHtml.Find("img").Each(func(i int, selection *goquery.Selection) {
			if i == 0 {
				operator.Painting = "http:" + selection.Nodes[0].Attr[1].Val
			}
		})
		page.Close()
		browser.Close()
		pw.Stop()
		operatorJson = append(operatorJson, operator)
	}
	operatorsJson, err := json.Marshal(operatorJson)
	if err != nil {
		log.Println(err)
		return
	}
	RedisSet("data_source", string(operatorsJson), 0)
	log.Println("数据源更新完毕")
}
