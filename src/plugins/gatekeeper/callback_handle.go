package gatekeeper

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

func CallBackData(callBack tgbotapi.Update) (bool, error) {
	callbackQuery := callBack.CallbackQuery
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 3 {
		return true, nil
	}

	userId, _ := strconv.ParseInt(d[0], 10, 64)
	chatId := callbackQuery.Message.Chat.ID
	name := utils.GetFullName(callbackQuery.From)

	if d[1] == "PASS" || d[1] == "BAN" {
		getChatMemberConfig := tgbotapi.GetChatMemberConfig{
			ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
				ChatID: chatId,
				UserID: callbackQuery.From.ID,
			},
		}
		memberInfo, _ := bot.Arknights.GetChatMember(getChatMemberConfig)

		if memberInfo.Status != "creator" && memberInfo.Status != "administrator" {
			answer := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "无使用权限！")
			bot.Arknights.Send(answer)
			return true, nil
		}

		if d[1] == "PASS" {
			pass(chatId, userId, callbackQuery, fmt.Sprintf("管理员通过了<a href=\"tg://user?id=%d\">%s</a>的验证", userId, d[2]))
		}

		if d[1] == "BAN" {
			chatMember := tgbotapi.ChatMemberConfig{ChatID: chatId, UserID: userId}
			ban(chatId, userId, callbackQuery, chatMember, fmt.Sprintf("管理员封禁了<a href=\"tg://user?id=%d\">%s</a>", userId, d[2]))
		}

		return true, nil
	}

	if userId != callbackQuery.From.ID {
		answer := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "这不是你的验证！")
		bot.Arknights.Send(answer)
		return true, nil
	}

	if d[1] != d[2] {
		answer := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "验证未通过，请一分钟后再试！")
		bot.Arknights.Send(answer)
		chatMember := tgbotapi.ChatMemberConfig{ChatID: chatId, UserID: userId}
		ban(chatId, userId, callbackQuery, chatMember, fmt.Sprintf("<a href=\"tg://user?id=%d\">%s</a>未通过验证，已被踢出。", userId, name))
		go unban(chatMember)
		return true, nil
	}

	answer := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "验证通过！")
	bot.Arknights.Send(answer)
	pass(chatId, userId, callbackQuery, fmt.Sprintf("<a href=\"tg://user?id=%d\">%s</a>已完成验证。", userId, name))

	return true, nil
}

func pass(chatId int64, userId int64, callbackQuery *tgbotapi.CallbackQuery, text string) string {
	bot.Arknights.Send(tgbotapi.RestrictChatMemberConfig{
		Permissions: &tgbotapi.ChatPermissions{
			CanSendMessages:       true,
			CanSendMediaMessages:  true,
			CanSendPolls:          true,
			CanSendOtherMessages:  true,
			CanAddWebPagePreviews: true,
			CanInviteUsers:        true,
			CanChangeInfo:         true,
			CanPinMessages:        true,
		},
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatId,
			UserID: userId,
		},
	})
	val := fmt.Sprintf("verify%d%d", chatId, userId)
	utils.RedisDelSetItem("verify", val)
	delMsg := tgbotapi.NewDeleteMessage(chatId, callbackQuery.Message.MessageID)
	bot.Arknights.Send(delMsg)
	sendMessage := tgbotapi.NewMessage(chatId, text)
	sendMessage.ParseMode = tgbotapi.ModeHTML
	msg, _ := bot.Arknights.Send(sendMessage)
	utils.AddDelQueue(msg.Chat.ID, msg.MessageID, 1)
	return val
}

func ban(chatId int64, userId int64, callbackQuery *tgbotapi.CallbackQuery, chatMember tgbotapi.ChatMemberConfig, text string) {
	kickChatMemberConfig := tgbotapi.KickChatMemberConfig{
		ChatMemberConfig: chatMember,
	}
	bot.Arknights.Send(kickChatMemberConfig)
	delMsg := tgbotapi.NewDeleteMessage(chatId, callbackQuery.Message.MessageID)
	bot.Arknights.Send(delMsg)
	val := fmt.Sprintf("verify%d%d", chatId, userId)
	utils.RedisDelSetItem("verify", val)
	sendMessage := tgbotapi.NewMessage(chatId, text)
	sendMessage.ParseMode = tgbotapi.ModeHTML
	msg, _ := bot.Arknights.Send(sendMessage)
	utils.AddDelQueue(msg.Chat.ID, msg.MessageID, 1)
}
