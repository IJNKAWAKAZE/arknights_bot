package gatekeeper

import (
	bot "arknights_bot/config"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"strings"
)

// isAdUser 检查昵称/用户名/简介是否命中广告词（大小写不敏感）
func isAdUser(member tgbotapi.User, bio string) bool {
	fields := []string{member.FirstName, member.LastName, member.UserName, bio}
	for _, word := range bot.ADWords {
		lowerWord := strings.ToLower(strings.TrimSpace(word))
		if lowerWord == "" {
			continue
		}
		for _, field := range fields {
			if strings.Contains(strings.ToLower(field), lowerWord) {
				return true
			}
		}
	}
	return false
}
