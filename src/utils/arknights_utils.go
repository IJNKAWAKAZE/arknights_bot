package utils

import (
	"arknights_bot/utils/suffixtree"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Operator struct {
	Name         string `json:"name"`         // 名字
	NameEn       string `json:"nameEn"`       // 英文名
	NameJa       string `json:"nameJp"`       // 日文名
	Code         string `json:"code"`         // 编号
	Race         string `json:"race"`         // 种族
	Profession   string `json:"profession"`   // 职业
	ProfessionZH string `json:"professionZH"` // 职业
	Rarity       int    `json:"rarity"`       // 稀有度
	Avatar       string `json:"avatar"`       // 头像
	ThumbURL     string `json:"thumbURL"`     // 半身像
	Skins        []Skin `json:"skins"`        // 皮肤
	HP           string `json:"hp"`           // 生命值
	ATK          string `json:"atk"`          // 攻击
	DEF          string `json:"def"`          // 防御
	Res          string `json:"res"`          // 法抗
	ReDeploy     string `json:"reDeploy"`     // 再部署时间
	Cost         string `json:"cost"`         // 费用
	Block        string `json:"block"`        // 阻挡数
	Interval     string `json:"interval"`     // 攻击间隔
	Sex          string `json:"sex"`          // 性别
	Position     string `json:"position"`     // 部署位
	Logo         string `json:"logo"`         // 所属
	ObtainMethod string `json:"obtainMethod"` // 获取方式
	Tags         string `json:"tags"`         // 标签
}

type Skin struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Material struct {
	ZoneName          string `json:"zoneName"`          // 关卡
	Code              string `json:"code"`              // 编码
	Name              string `json:"name"`              // 主材料名称
	Icon              string `json:"icon"`              // 主材料图标
	KnockRating       string `json:"knockRating"`       // 主产物掉率
	ApExpect          string `json:"apExpect"`          // 期望理智
	SecondaryItem     string `json:"SecondaryItem"`     // 副产物名称
	SecondaryItemIcon string `json:"SecondaryItemIcon"` // 副产物图标
	StageEfficiency   string `json:"stageEfficiency"`   // 关卡效率
}

var operatorMap = make(map[string]Operator)
var recruitOperatorList []Operator
var DataNeedUpdate = true
var operatorTree suffixtree.GST
var itemTree suffixtree.GST
var itemArray []string
var enemyTree suffixtree.GST
var enemyArray []pair
var operators []Operator

func GetOperators() []Operator {
	updateData()
	return operators
}
func updateData() {
	if !DataNeedUpdate {
		return
	}
	//operators
	operatorsJson := RedisGet("operatorList")
	json.Unmarshal([]byte(operatorsJson), &operators)
	operatorMap = make(map[string]Operator)
	operatorTree = suffixtree.NewGeneralizedSuffixTree()
	for index, operator := range operators {
		operatorMap[strings.ToLower(operator.Name)] = operator
		operatorTree.Put(strings.ToLower(operator.Name), index)
		if strings.Contains(operator.ObtainMethod, "公开招募") {
			recruitOperatorList = append(recruitOperatorList, operator)
		}
	}
	//enemy
	func() {
		resultArray, resultTree := fetchEnemiesData()
		enemyArray = resultArray
		enemyTree = resultTree
		defer func() {
			if err := recover(); err != nil {
				log.Fatal("Can not update enemy")
			}
		}()
	}()

	//items
	//gjson.Parse(RedisGet("materialMap"))
	//TODO: check what is the materialMap
	//set flag
	DataNeedUpdate = false

}

type pair struct {
	a, b interface{}
}

func fetchEnemiesData() ([]pair, suffixtree.GST) {
	makeurl := func(n string) string {
		paintingName := fmt.Sprintf("头像_敌人_%s.png", n)
		m := Md5(paintingName)
		path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
		return path + url.PathEscape(paintingName)
	}
	emeryTree := suffixtree.NewGeneralizedSuffixTree()
	var newEnemyArray []pair
	api := viper.GetString("api.enemy")
	response, _ := http.Get(api)
	e, _ := io.ReadAll(response.Body)
	defer response.Body.Close()
	enemyJson := gjson.ParseBytes(e)
	for index, en := range enemyJson.Array() {
		n := en.Get("name").String()
		newEnemyArray = append(newEnemyArray, pair{n, makeurl(n)})
		emeryTree.Put(strings.ToLower(n), index)
	}
	return newEnemyArray, emeryTree
}
func GetOperatorByName(name string) Operator {
	updateData()
	return operatorMap[name]
}

func GetOperatorsByName(name string) []Operator {
	updateData()
	var operatorList []Operator
	for _, op := range operatorTree.Search(strings.ToLower(name)) {
		print(operators[op].Name)
		operatorList = append(operatorList, operators[op])
	}
	println()
	return operatorList
}

func GetRecruitOperatorList() []Operator {
	updateData()
	return recruitOperatorList
}

func GetEnemiesByName(name string) map[string]string {
	updateData()
	var enemyMap = make(map[string]string)
	for _, index := range enemyTree.Search(strings.ToLower(name)) {
		a := enemyArray[index]
		enemyMap[a.a.(string)] = a.b.(string)
	}
	return enemyMap
}

func GetItemsByName(name string) map[string]string {
	var materialMap = make(map[string]string)
	materialJson := RedisGet("materialMap")
	gjson.Parse(materialJson).ForEach(func(key, value gjson.Result) bool {
		if strings.Contains(strings.ToLower(key.String()), strings.ToLower(name)) {
			materialMap[key.String()] = value.Get("0.icon").String()
		}
		return true
	})
	return materialMap
}

func GetItemByName(name string) []Material {
	var materials []Material
	materialJson := RedisGet("materialMap")
	gjson.Parse(materialJson).ForEach(func(key, value gjson.Result) bool {
		if strings.ToLower(key.String()) == strings.ToLower(name) {
			json.Unmarshal([]byte(value.String()), &materials)
			return false
		}
		return true
	})
	return materials
}
