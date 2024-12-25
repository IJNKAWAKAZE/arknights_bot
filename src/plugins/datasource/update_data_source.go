package datasource

import (
	"arknights_bot/config"
	"arknights_bot/utils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
	"github.com/starudream/go-lib/core/v2/codec/json"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
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
func UpdateDataSource() {
	go UpdateDataSourceRunner()
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
		attrs := selection.Nodes[0].Attr
		operator.Name = attrs[0].Val
		operator.Profession = Profession[attrs[1].Val]
		operator.ProfessionZH = attrs[1].Val + "干员"
		operator.Rarity, _ = strconv.Atoi(attrs[2].Val)
		operator.Logo = attrs[3].Val
		operator.Race = attrs[6].Val
		operator.NameEn = attrs[7].Val
		operator.NameJa = attrs[8].Val
		operator.Code = attrs[9].Val
		operator.HP = attrs[10].Val
		operator.ATK = attrs[11].Val
		operator.DEF = attrs[12].Val
		operator.Res = attrs[13].Val
		operator.ReDeploy = attrs[14].Val
		c := strings.Split(attrs[15].Val, "→")
		operator.Cost = c[len(c)-1]
		b := strings.Split(attrs[16].Val, "→")
		operator.Block = b[len(b)-1]
		operator.Interval = attrs[17].Val
		operator.Sex = attrs[18].Val
		operator.Position = attrs[19].Val
		operator.Tags = strings.ReplaceAll(attrs[20].Val, "支援机械", "机械")
		operator.ObtainMethod = attrs[21].Val
		// 头像
		paintingName := fmt.Sprintf("头像_%s.png", operator.Name)
		m := utils.Md5(paintingName)
		path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
		operator.Avatar = path + paintingName + "?image_process=format,webp/quality,Q_90"
		// 半身像
		paintingName = fmt.Sprintf("半身像_%s_1.png", operator.Name)
		m = utils.Md5(paintingName)
		path = "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
		operator.ThumbURL = path + paintingName + "?image_process=format,webp/quality,Q_90"
		operators = append(operators, operator)
	})

	// 老干员map
	var oldOperators = make(map[string]string)
	var os []utils.Operator
	operatorsJson := utils.RedisGet("operatorList")
	json.Unmarshal([]byte(operatorsJson), &os)
	for _, o := range os {
		oldOperators[o.Name] = o.Name
	}

	skinCount := make(map[string][]string)
	response, _ = http.Get(api + "时装回廊")
	doc, _ = goquery.NewDocumentFromReader(response.Body)
	doc.Find(".skinwrapper").Each(func(i int, selection *goquery.Selection) {
		img, _ := url.QueryUnescape(selection.Find(".charimg").First().Nodes[0].FirstChild.Attr[1].Val)
		skinName := selection.Find(".charnameEn").Text()
		compileRegex := regexp.MustCompile("_(.*?)_")
		match := compileRegex.FindStringSubmatch(img)
		if len(match) > 1 {
			name := match[1]
			skinCount[name] = append(skinCount[name], skinName)
		}
	})

	var birthdayMap = make(map[string][]utils.Operator)
	for i, operator := range operators {
		name := operator.Name
		if name == "阿米娅" {
			// 立绘
			for e := 0; e < 2; e++ {
				paintingName := fmt.Sprintf("立绘_%s_%d.png", name, e+1)
				m := utils.Md5(paintingName)
				path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
				painting := path + paintingName + "?image_process=format,webp/quality,Q_90"
				var skin utils.Skin
				skin.Url = painting
				operators[i].Skins = append(operators[i].Skins, skin)
			}
			// 精1立绘
			paintingName := fmt.Sprintf("立绘_%s_1+.png", name)
			m := utils.Md5(paintingName)
			path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
			painting := path + paintingName + "?image_process=format,webp/quality,Q_90"
			var skin utils.Skin
			skin.Url = painting
			operators[i].Skins = append(operators[i].Skins, skin)
			// 皮肤
			for c, sk := range skinCount[name] {
				paintingName := fmt.Sprintf("立绘_%s_skin%d.png", name, len(skinCount[name])-c)
				m := utils.Md5(paintingName)
				path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
				painting := path + paintingName + "?image_process=format,webp/quality,Q_90"
				var skin utils.Skin
				skin.Name = sk
				skin.Url = painting
				operators[i].Skins = append(operators[i].Skins, skin)
			}
		} else if name == "阿米娅(近卫)" || name == "阿米娅(医疗)" {
			// 立绘
			paintingName := fmt.Sprintf("立绘_%s_2.png", name)
			m := utils.Md5(paintingName)
			path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
			painting := path + paintingName + "?image_process=format,webp/quality,Q_90"
			var skin utils.Skin
			skin.Url = painting
			operators[i].Skins = append(operators[i].Skins, skin)
			// 皮肤
			for c, sk := range skinCount[name] {
				paintingName := fmt.Sprintf("立绘_%s_skin%d.png", name, len(skinCount[name])-c)
				m := utils.Md5(paintingName)
				path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
				painting := path + paintingName + "?image_process=format,webp/quality,Q_90"
				var skin utils.Skin
				skin.Name = sk
				skin.Url = painting
				operators[i].Skins = append(operators[i].Skins, skin)
			}
		}
		if operator.Rarity < 3 {
			// 立绘
			paintingName := fmt.Sprintf("立绘_%s_1.png", name)
			m := utils.Md5(paintingName)
			path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
			painting := path + paintingName + "?image_process=format,webp/quality,Q_90"
			var skin utils.Skin
			skin.Url = painting
			operators[i].Skins = append(operators[i].Skins, skin)
			// 皮肤
			for c, sk := range skinCount[name] {
				paintingName := fmt.Sprintf("立绘_%s_skin%d.png", name, len(skinCount[name])-c)
				m := utils.Md5(paintingName)
				path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
				painting := path + paintingName + "?image_process=format,webp/quality,Q_90"
				var skin utils.Skin
				skin.Name = sk
				skin.Url = painting
				operators[i].Skins = append(operators[i].Skins, skin)
			}
		} else {
			// 立绘
			for e := 0; e < 2; e++ {
				paintingName := fmt.Sprintf("立绘_%s_%d.png", name, e+1)
				m := utils.Md5(paintingName)
				path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
				painting := path + paintingName + "?image_process=format,webp/quality,Q_90"
				var skin utils.Skin
				skin.Url = painting
				operators[i].Skins = append(operators[i].Skins, skin)
			}
			// 皮肤
			for c, sk := range skinCount[name] {
				paintingName := fmt.Sprintf("立绘_%s_skin%d.png", name, len(skinCount[name])-c)
				m := utils.Md5(paintingName)
				path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
				painting := path + paintingName + "?image_process=format,webp/quality,Q_90"
				var skin utils.Skin
				skin.Name = sk
				skin.Url = painting
				operators[i].Skins = append(operators[i].Skins, skin)
			}
		}
		// 新增干员
		if _, has := oldOperators[name]; !has && config.IgnoreBirthday[name] == "" {
			response, _ := http.Get(api + name)
			doc, _ := goquery.NewDocumentFromReader(response.Body)
			doc.Find(".poem").Each(func(j int, selection *goquery.Selection) {
				text := selection.Text()
				if strings.Contains(text, "【生日】") {
					t := strings.Split(text, "\n")
					birthday := t[5][strings.Index(t[5], "】")+3:]
					reg := regexp.MustCompile("[0-9]+")
					if reg.MatchString(birthday) {
						birthdayMap[birthday] = append(birthdayMap[birthday], operators[i])
					} else {
						birthdayMap["未知"] = append(birthdayMap["未知"], operators[i])
					}
					return
				}
			})
		}
	}

	for k, v := range birthdayMap {
		utils.RedisSet("birthday:"+k, json.MustMarshalString(v), 0)
	}

	defer response.Body.Close()

	utils.RedisSet("operatorList", json.MustMarshalString(operators), 0)
	MaterialInfo()
	log.Println("数据源更新完毕")
	utils.DataNeedUpdate = true
}

func MaterialInfo() {
	var itemMap = make(map[string]string)
	var materialMap = make(map[string][]utils.Material)
	res, err := http.Get(viper.GetString("api.item"))
	if err != nil {
		log.Println(err)
		return
	}
	read, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	j := gjson.ParseBytes(read)
	for _, item := range j.Get("data").Array() {
		itemMap[item.Get("itemId").String()] = item.Get("itemName").String()
	}

	res, err = http.Get(viper.GetString("api.stage_result"))
	if err != nil {
		log.Println(err)
		return
	}
	read, err = io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	j = gjson.ParseBytes(read)
	for _, d := range j.Get("data.recommendedStageList").Array() {
		for _, item := range d.Get("stageResultList").Array() {
			var material utils.Material
			material.ZoneName = item.Get("zoneName").String()
			material.Code = item.Get("stageCode").String()
			material.Name = item.Get("itemName").String()
			// 图标
			paintingName := fmt.Sprintf("道具_带框_%s.png", material.Name)
			m := utils.Md5(paintingName)
			path := "https://media.prts.wiki/thumb" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
			pic := path + paintingName + "/75px-" + paintingName
			material.Icon = pic
			material.ApExpect = fmt.Sprintf("%.1f", item.Get("apExpect").Float())
			material.KnockRating = fmt.Sprintf("%.1f%%", item.Get("knockRating").Float()*100)
			material.SecondaryItem = itemMap[item.Get("secondaryItemId").String()]
			if material.SecondaryItem != "" {
				paintingName := fmt.Sprintf("道具_带框_%s.png", material.SecondaryItem)
				m := utils.Md5(paintingName)
				path := "https://media.prts.wiki/thumb" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
				pic := path + paintingName + "/75px-" + paintingName
				material.SecondaryItemIcon = pic
			}
			material.StageEfficiency = fmt.Sprintf("%.1f%%", item.Get("stageEfficiency").Float()*100)
			materialMap[material.Name] = append(materialMap[material.Name], material)
		}
	}
	utils.RedisSet("materialMap", json.MustMarshalString(materialMap), 0)
}
