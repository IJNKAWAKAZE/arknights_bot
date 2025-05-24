package utils

import (
	"testing"
	"strings"
	"arknights_bot/utils/suffixtree"
)

// 测试专用的简化Operator结构体
type testOperator struct {
	Name   string
	Pinyin [][]string
}

// 将测试数据转换为Operator格式
func toOperator(op testOperator) Operator {
	return Operator{
		Name:   op.Name,
		Pinyin: op.Pinyin,
	}
}

func TestGetOperatorByName(t *testing.T) {
	// 保存原始状态
	oldOperators := operators
	oldOperatorMap := operatorMap
	oldOperatorTree := operatorTree
	oldIsTesting := isTesting
	
	// 测试完成后恢复
	defer func() {
		operators = oldOperators
		operatorMap = oldOperatorMap
		operatorTree = oldOperatorTree
		isTesting = oldIsTesting
	}()
	
	// 设置测试模式
	isTesting = true
	
	// 初始化测试数据
	testData := []testOperator{
		{
			Name:   "陈",
			Pinyin: [][]string{{"chen", "zhen"}},
		},
		{
			Name:   "陈晖洁",
			Pinyin: [][]string{{"chen", "zhen"}, {"hui"}, {"jie"}},
		},
		{
			Name:   "阿米娅",
			Pinyin: [][]string{{"a"}, {"mi"}, {"ya"}},
		},
	}
	
	// 转换为Operator并建立索引
	operators = make([]Operator, len(testData))
	operatorMap = make(map[string]Operator)
	operatorTree = suffixtree.NewGeneralizedSuffixTree()
	
	for index, op := range testData {
		operator := toOperator(op)
		operators[index] = operator
		
		// 原始名称索引
		operatorMap[strings.ToLower(operator.Name)] = operator
		operatorTree.Put(strings.ToLower(operator.Name), index)
		
		// 拼音组合索引
		var possibleKeys []string
		for _, charPinyin := range operator.Pinyin {
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
		
		for _, key := range possibleKeys {
			lowerKey := strings.ToLower(key)
			operatorMap[lowerKey] = operator
			operatorTree.Put(lowerKey, index)
		}
	}

	// 测试用例
	testCases := []struct {
		name     string
		expected bool
	}{
		// 单字测试
		{"陈", true},      // 中文精确匹配
		{"chen", true},   // 拼音精确匹配
		{"zhen", true},   // 多音字拼音匹配
		{"chn", false},   // 错误拼音不匹配
		{"辰", true},      // 同音字匹配
		{"尘", true},      // 同音字匹配
		{"chen2", false}, // 错误拼音不匹配
		{"ch", true},     // 部分拼音匹配
		{"foo", false},   // 完全不匹配

		// 多字测试
		{"陈晖洁", true},          // 多字中文精确匹配
		{"chenhuijie", true},    // 多字拼音精确匹配
		{"zhenhuijie", true},    // 多字多音字拼音匹配
		{"辰晖洁", true},         // 同音字匹配
		{"陈晖杰", true},         // 同音字匹配
		{"尘晖洁", true},         // 同音字匹配
		{"chenhuij", true},      // 多字部分拼音匹配
		{"chenhui", true},       // 多字不完整拼音匹配
		{"陈hui洁", false},       // 混合格式不匹配
		{"陈 晖 洁", true},       // 带空格匹配
		{"chenhuijie2", false},  // 多字错误拼音不匹配
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetOperatorByName(tc.name)
			if (result.Name != "") != tc.expected {
				t.Errorf("搜索 %s 期望 %v 但得到 %v", tc.name, tc.expected, result.Name != "")
			}
		})
	}
}

func TestPinyinSearch(t *testing.T) {
	// 保存原始全局变量
	oldOperators := operators
	oldOperatorMap := operatorMap
	oldOperatorTree := operatorTree
	
	// 测试完成后恢复
	defer func() {
		operators = oldOperators
		operatorMap = oldOperatorMap
		operatorTree = oldOperatorTree
	}()

	// 初始化测试数据
	testOperator := Operator{
		Name: "陈",
		Pinyin: [][]string{{"chen", "zhen"}}, // 添加多音字变体
	}
	operators = []Operator{testOperator}
	operatorMap = map[string]Operator{
		"陈": testOperator,
		"chen": testOperator,
		"zhen": testOperator,
	}
	operatorTree = suffixtree.NewGeneralizedSuffixTree()
	operatorTree.Put("陈", 0)
	operatorTree.Put("chen", 0)
	operatorTree.Put("zhen", 0)

	// 测试拼音搜索
	testCases := []struct {
		name     string
		expected bool
	}{
		{"陈", true},
		{"chen", true},
		{"zhen", true},
		{"foo", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 直接调用搜索逻辑，不触发updateData
			var result []Operator
			set := make(map[int]bool)
			
			for _, op := range operatorTree.Search(strings.ToLower(tc.name)) {
				if !set[op] {
					set[op] = true
					result = append(result, operators[op])
				}
			}
			
			if (len(result) > 0) != tc.expected {
				t.Errorf("搜索 %s 期望 %v 但得到 %v", tc.name, tc.expected, len(result) > 0)
			}
		})
	}
}