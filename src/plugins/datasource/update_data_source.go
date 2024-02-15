package datasource

import (
	"arknights_bot/utils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
	"github.com/starudream/go-lib/core/v2/codec/json"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

var Profession = make(map[string]string)

func init() {
	profession := make(map[string]string)
	profession["术师"] = "CASTER"
	profession["医疗"] = "MEDIC"
	profession["先锋"] = "PIONEER"
	profession["狙击"] = "SNIPER"
	profession["特种"] = "SPECIAL"
	profession["辅助"] = "SUPPORT"
	profession["重装"] = "TANK"
	profession["近卫"] = "WARRIOR"
	Profession = profession
}

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
	var operators []utils.Operator
	api := viper.GetString("api.wiki")
	response, _ := http.Get(api + "干员一览")
	doc, _ := goquery.NewDocumentFromReader(response.Body)
	doc.Find("#filter-data div").Each(func(i int, selection *goquery.Selection) {
		var operator utils.Operator
		operator.Name = selection.Nodes[0].Attr[0].Val
		operator.Profession = Profession[selection.Nodes[0].Attr[1].Val]
		operator.Rarity, _ = strconv.Atoi(selection.Nodes[0].Attr[2].Val)

		// 半身像
		paintingName := fmt.Sprintf("半身像_%s_1.png", operator.Name)
		m := utils.Md5(paintingName)
		path := "https://prts.wiki" + fmt.Sprintf("/images/%s/%s/", m[:1], m[:2])
		operator.ThumbURL = path + paintingName + "?image_process=format,webp/quality,Q_90"
		operators = append(operators, operator)
	})

	skinCount := make(map[string]int)
	response, _ = http.Get(api + "时装回廊")
	doc, _ = goquery.NewDocumentFromReader(response.Body)
	doc.Find(".charimg").Each(func(i int, selection *goquery.Selection) {
		img, _ := url.QueryUnescape(selection.Nodes[0].FirstChild.Attr[0].Val)
		compileRegex := regexp.MustCompile("_(.*?)_")
		match := compileRegex.FindStringSubmatch(img)
		name := match[1]
		if skinCount[name] == 0 {
			skinCount[name] = 1
		} else {
			skinCount[name] = skinCount[name] + 1
		}
	})

	for i, operator := range operators {
		name := operator.Name
		if name == "阿米娅" {
			// 立绘
			for e := 0; e < 2; e++ {
				paintingName := fmt.Sprintf("立绘_%s_%d.png", name, e+1)
				m := utils.Md5(paintingName)
				path := "https://prts.wiki" + fmt.Sprintf("/images/%s/%s/", m[:1], m[:2])
				painting := path + paintingName + "?image_process=format,webp/quality,Q_90"
				operators[i].Skins = append(operators[i].Skins, painting)
			}
			// 精1立绘
			paintingName := fmt.Sprintf("立绘_%s_1+.png", name)
			m := utils.Md5(paintingName)
			path := "https://prts.wiki" + fmt.Sprintf("/images/%s/%s/", m[:1], m[:2])
			painting := path + paintingName + "?image_process=format,webp/quality,Q_90"
			operators[i].Skins = append(operators[i].Skins, painting)
			// 皮肤
			for c := 0; c < skinCount[name]; c++ {
				paintingName := fmt.Sprintf("立绘_%s_skin%d.png", name, c+1)
				m := utils.Md5(paintingName)
				path := "https://prts.wiki" + fmt.Sprintf("/images/%s/%s/", m[:1], m[:2])
				painting := path + paintingName + "?image_process=format,webp/quality,Q_90"
				operators[i].Skins = append(operators[i].Skins, painting)
			}
			continue
		}
		if name == "阿米娅(近卫)" {
			// 立绘
			paintingName := fmt.Sprintf("立绘_%s_2.png", name)
			m := utils.Md5(paintingName)
			path := "https://prts.wiki" + fmt.Sprintf("/images/%s/%s/", m[:1], m[:2])
			painting := path + paintingName + "?image_process=format,webp/quality,Q_90"
			operators[i].Skins = append(operators[i].Skins, painting)
			// 皮肤
			for c := 0; c < skinCount[name]; c++ {
				paintingName := fmt.Sprintf("立绘_%s_skin%d.png", name, c+1)
				m := utils.Md5(paintingName)
				path := "https://prts.wiki" + fmt.Sprintf("/images/%s/%s/", m[:1], m[:2])
				painting := path + paintingName + "?image_process=format,webp/quality,Q_90"
				operators[i].Skins = append(operators[i].Skins, painting)
			}
			continue
		}
		if operator.Rarity < 3 {
			// 立绘
			paintingName := fmt.Sprintf("立绘_%s_1.png", name)
			m := utils.Md5(paintingName)
			path := "https://prts.wiki" + fmt.Sprintf("/images/%s/%s/", m[:1], m[:2])
			painting := path + paintingName + "?image_process=format,webp/quality,Q_90"
			operators[i].Skins = append(operators[i].Skins, painting)
			// 皮肤
			for c := 0; c < skinCount[name]; c++ {
				paintingName := fmt.Sprintf("立绘_%s_skin%d.png", name, c+1)
				m := utils.Md5(paintingName)
				path := "https://prts.wiki" + fmt.Sprintf("/images/%s/%s/", m[:1], m[:2])
				painting := path + paintingName + "?image_process=format,webp/quality,Q_90"
				operators[i].Skins = append(operators[i].Skins, painting)
			}
		} else {
			// 立绘
			for e := 0; e < 2; e++ {
				paintingName := fmt.Sprintf("立绘_%s_%d.png", name, e+1)
				m := utils.Md5(paintingName)
				path := "https://prts.wiki" + fmt.Sprintf("/images/%s/%s/", m[:1], m[:2])
				painting := path + paintingName + "?image_process=format,webp/quality,Q_90"
				operators[i].Skins = append(operators[i].Skins, painting)
			}
			// 皮肤
			for c := 0; c < skinCount[name]; c++ {
				paintingName := fmt.Sprintf("立绘_%s_skin%d.png", name, c+1)
				m := utils.Md5(paintingName)
				path := "https://prts.wiki" + fmt.Sprintf("/images/%s/%s/", m[:1], m[:2])
				painting := path + paintingName + "?image_process=format,webp/quality,Q_90"
				operators[i].Skins = append(operators[i].Skins, painting)
			}
		}
	}

	defer response.Body.Close()
	utils.RedisSet("data_source", json.MustMarshalString(operators), 0)
	log.Println("数据源更新完毕")
}
