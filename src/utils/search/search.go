package search

import (
	"arknights_bot/utils/cache"
	"arknights_bot/utils/hashutil"
	"arknights_bot/utils/model"
	"arknights_bot/utils/pinyin"
	"arknights_bot/utils/suffixtree"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var operatorMap = make(map[string]model.Operator)
var recruitOperatorList []model.Operator
var DataNeedUpdate = true
var operatorTree suffixtree.GST
var enemyTree suffixtree.GST
var enemyArray []pair
var operators []model.Operator

// dataMu 保护干员/敌人/物品数据的并发读写，防止重建期间读取到混合状态
var dataMu sync.RWMutex

// SetDataNeedUpdate 标记数据需要重新加载（数据源更新完成后调用）
func SetDataNeedUpdate() {
	dataMu.Lock()
	DataNeedUpdate = true
	dataMu.Unlock()
}

// ensureData 数据就绪检查，仅在需要重建时获取写锁（double-check）
func ensureData() {
	dataMu.RLock()
	need := DataNeedUpdate
	dataMu.RUnlock()
	if !need {
		return
	}
	dataMu.Lock()
	defer dataMu.Unlock()
	if !DataNeedUpdate {
		return
	}
	updateData()
}

func GetOperators() []model.Operator {
	ensureData()
	dataMu.RLock()
	defer dataMu.RUnlock()
	return operators
}
func updateData() {
	if !DataNeedUpdate {
		return
	}
	// 重置
	recruitOperatorList = nil

	//operators
	operatorsJson := cache.RedisGet("operatorList")
	json.Unmarshal([]byte(operatorsJson), &operators)
	operatorMap = make(map[string]model.Operator)
	operatorTree = suffixtree.NewGeneralizedSuffixTree()
	for index, operator := range operators {
		// 生成所有拼音变体并建立索引
		pinyinArgs := pinyin.NewArgs()
		pinyinArgs.Style = pinyin.Normal
		pinyinArgs.Heteronym = true // 启用多音字模式

		// 原始名称索引
		operatorMap[strings.ToLower(operator.Name)] = operator
		operatorTree.Put(strings.ToLower(operator.Name), index)

		// 生成所有拼音变体
		variations := pinyin.NameVariations(operator.Name, pinyinArgs)
		operator.Pinyin = variations // 存储拼音变体供后续使用

		// 生成所有可能的拼音组合并建立索引
		var possibleKeys []string
		for _, charPinyin := range variations {
			if len(possibleKeys) == 0 {
				possibleKeys = charPinyin
				continue
			}
			var newKeys []string
			for _, key := range possibleKeys {
				for _, py := range charPinyin {
					newKeys = append(newKeys, key+py)
				}
			}
			possibleKeys = newKeys
		}

		// 索引所有拼音组合
		for _, key := range possibleKeys {
			lowerKey := strings.ToLower(key)
			if _, exists := operatorMap[lowerKey]; !exists {
				operatorMap[lowerKey] = operator
				//operatorTree.Put(lowerKey, index)
			}
		}
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

	//set flag
	DataNeedUpdate = false

}

type pair struct {
	a, b interface{}
}

func fetchEnemiesData() ([]pair, suffixtree.GST) {
	makeurl := func(n string) string {
		paintingName := fmt.Sprintf("头像_敌人_%s.png", n)
		m := hashutil.Md5(paintingName)
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

var isTesting = false

func GetOperatorByName(name string) model.Operator {
	if !isTesting {
		ensureData()
	}
	dataMu.RLock()
	defer dataMu.RUnlock()

	// 先尝试精确匹配
	lowerName := strings.ToLower(name)
	if op, ok := operatorMap[lowerName]; ok {
		return op
	}

	// 生成输入名的所有拼音组合
	pinyinArgs := pinyin.NewArgs()
	pinyinArgs.Style = pinyin.Normal
	pinyinArgs.Heteronym = true
	inputPinyin := pinyin.Pinyin(name, pinyinArgs)

	// 生成所有可能的拼音组合
	var possibleKeys []string
	for _, charPinyin := range inputPinyin {
		if len(possibleKeys) == 0 {
			possibleKeys = charPinyin
			continue
		}
		var newKeys []string
		for _, key := range possibleKeys {
			for _, py := range charPinyin {
				newKeys = append(newKeys, key+py)
			}
		}
		possibleKeys = newKeys
	}

	// 检查每个拼音组合是否匹配
	for _, key := range possibleKeys {
		if op, ok := operatorMap[strings.ToLower(key)]; ok {
			return op
		}
	}

	// 最后尝试后缀树搜索
	if indices := operatorTree.Search(lowerName); len(indices) > 0 {
		return operators[indices[0]]
	}

	return model.Operator{} // 未找到返回空结构体
}

func GetOperatorsByName(name string) []model.Operator {
	ensureData()
	dataMu.RLock()
	defer dataMu.RUnlock()
	var operatorList []model.Operator
	var set = make(map[int]bool)

	// 原始名称搜索
	for _, op := range operatorTree.Search(strings.ToLower(name)) {
		_, contain := set[op]
		if !contain {
			set[op] = true
			operatorList = append(operatorList, operators[op])
		}
	}

	// 拼音搜索 - 使用预先生成的拼音索引
	pinyinArgs := pinyin.NewArgs()
	pinyinArgs.Style = pinyin.Normal
	pinyinArgs.Heteronym = true
	inputPinyin := pinyin.Pinyin(name, pinyinArgs)

	// 生成所有可能的拼音组合
	var possibleKeys []string
	for _, charPinyin := range inputPinyin {
		if len(possibleKeys) == 0 {
			possibleKeys = charPinyin
			continue
		}
		var newKeys []string
		for _, key := range possibleKeys {
			for _, py := range charPinyin {
				newKeys = append(newKeys, key+py)
			}
		}
		possibleKeys = newKeys
	}

	// 搜索所有拼音组合
	for _, key := range possibleKeys {
		for _, op := range operatorTree.Search(strings.ToLower(key)) {
			_, contain := set[op]
			if !contain {
				set[op] = true
				operatorList = append(operatorList, operators[op])
			}
		}
	}

	return operatorList
}

func GetRecruitOperatorList() []model.Operator {
	ensureData()
	dataMu.RLock()
	defer dataMu.RUnlock()
	return recruitOperatorList
}

func GetEnemiesByName(name string) map[string]string {
	ensureData()
	dataMu.RLock()
	defer dataMu.RUnlock()
	var enemyMap = make(map[string]string)
	for _, index := range enemyTree.Search(strings.ToLower(name)) {
		a := enemyArray[index]
		enemyMap[a.a.(string)] = a.b.(string)
	}
	return enemyMap
}
