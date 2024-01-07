package datasource

import (
	"arknights_bot/utils"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strings"
)

// UpdateDataSource 更新数据源
func UpdateDataSource() func() {
	updateDataSource := func() {
		go UpdateDataSourceRunner()
	}
	return updateDataSource
}

// UpdateDataSourceRunner 更新数据源
func UpdateDataSourceRunner() {
	log.Println("开始更新数据源...")
	var operatorJson []Verify
	var operatorList []string
	api := viper.GetString("api.wiki")
	response, _ := http.Get(api + "干员一览")
	doc, _ := goquery.NewDocumentFromReader(response.Body)
	doc.Find("#filter-data div").Each(func(i int, selection *goquery.Selection) {
		operatorList = append(operatorList, selection.Nodes[0].Attr[0].Val)
	})
	for _, name := range operatorList {
		var operator Verify
		operator.Name = name
		pw, err := playwright.Run()
		if err != nil {
			log.Println("未检测到playwright，开始自动安装...")
			playwright.Install()
			pw, _ = playwright.Run()
		}
		browser, _ := pw.Chromium.Launch()
		page, _ := browser.NewPage()
		page.Goto(api+name, playwright.PageGotoOptions{
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
	utils.RedisSet("data_source", string(operatorsJson), 0)
	log.Println("数据源更新完毕")
}
