package gatekeeper

import (
	"arknights_bot/config"
	"strings"
	"unicode"
)

// normalizeText 去掉空白并转为小写，用于宽松匹配广告词
func normalizeText(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return unicode.ToLower(r)
	}, s)
}

// hasAdWord 判断 bio 是否包含广告关键词
func hasAdWord(bio string) bool {
	config.DataMu.RLock()
	defer config.DataMu.RUnlock()
	b := normalizeText(bio)
	for _, word := range config.ADWords {
		if word != "" && strings.Contains(b, normalizeText(word)) {
			return true
		}
	}
	return false
}
