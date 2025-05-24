package pinyin

import (
	"strings"
	"github.com/mozillazg/go-pinyin"
)

type Args struct {
	Style     int // 拼音风格
	Heteronym bool // 是否启用多音字模式
}

const (
	Tone         = pinyin.Tone         // 带声调
	Tone2        = pinyin.Tone2        // 声调在拼音之后
	Tone3        = pinyin.Tone3        // 声调在韵母之后
	Initials     = pinyin.Initials     // 仅声母
	FirstLetter  = pinyin.FirstLetter  // 首字母
	Finals       = pinyin.Finals       // 仅韵母
	FinalsTone   = pinyin.FinalsTone   // 仅韵母带声调
	FinalsTone2  = pinyin.FinalsTone2  // 仅韵母声调在拼音之后
	FinalsTone3  = pinyin.FinalsTone3  // 仅韵母声调在韵母之后
	Normal       = pinyin.Normal       // 普通风格，不带声调
)

func NewArgs() *Args {
	return &Args{
		Style:     Normal,
		Heteronym: false,
	}
}

// Pinyin 将中文转换为拼音
func Pinyin(s string, args *Args) [][]string {
	a := pinyin.NewArgs()
	a.Style = args.Style
	a.Heteronym = args.Heteronym
	return pinyin.Pinyin(s, a)
}

// Convert 将中文转换为拼音字符串
func Convert(s string, args *Args) string {
	pys := Pinyin(s, args)
	var result strings.Builder
	for _, words := range pys {
		for _, word := range words {
			result.WriteString(word)
		}
	}
	return result.String()
}

// Homophones 获取同音字
func Homophones(s string, args *Args) map[string][]string {
	result := make(map[string][]string)
	pys := Pinyin(s, args)
	
	for _, words := range pys {
		for _, word := range words {
			if _, exists := result[word]; !exists {
				result[word] = []string{}
			}
		}
	}
	
	return result
}

// Initial 获取拼音首字母
func Initial(s string) string {
	// 非中文字符直接返回
	for _, r := range s {
		if r < 0x4E00 || r > 0x9FA5 {
			return s
		}
	}
	args := NewArgs()
	args.Style = FirstLetter
	return Convert(s, args)
}

// NameVariations 获取名字中每个字的所有拼音
func NameVariations(name string, args *Args) [][]string {
	result := make([][]string, 0)
	for _, char := range name {
		charStr := string(char)
		// 非中文字符直接作为单元素返回
		if char < 0x4E00 || char > 0x9FA5 {
			result = append(result, []string{charStr})
			continue
		}
		pinyins := Pinyin(charStr, args)
		if len(pinyins) > 0 {
			// 去重
			unique := make(map[string]bool)
			var uniquePinyins []string
			for _, py := range pinyins[0] {
				if !unique[py] {
					unique[py] = true
					uniquePinyins = append(uniquePinyins, py)
				}
			}
			result = append(result, uniquePinyins)
		} else {
			result = append(result, []string{charStr})
		}
	}
	return result
}
