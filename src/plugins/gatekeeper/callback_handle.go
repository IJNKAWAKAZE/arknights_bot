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

	if len(d) < 5 {
		return true, nil
	}

	userId, _ := strconv.ParseInt(d[1], 10, 64)
	chatId := callbackQuery.Message.Chat.ID
	joinMessageId, _ := strconv.Atoi(d[4])

	if d[2] == "PASS" || d[2] == "BAN" {
		getChatMemberConfig := tgbotapi.GetChatMemberConfig{
			ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
				ChatID: chatId,
				UserID: callbackQuery.From.ID,
			},
		}

		if !utils.IsAdmin(getChatMemberConfig) {
			answer := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "无使用权限！")
			bot.Arknights.Send(answer)
			return true, nil
		}

		if d[2] == "PASS" {
			pass(chatId, userId, callbackQuery)
		}

		if d[2] == "BAN" {
			chatMember := tgbotapi.ChatMemberConfig{ChatID: chatId, UserID: userId}
			ban(chatId, userId, callbackQuery, chatMember, joinMessageId)
		}

		return true, nil
	}

	if userId != callbackQuery.From.ID {
		answer := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "这不是你的验证！")
		bot.Arknights.Send(answer)
		return true, nil
	}

	if d[2] != d[3] {
		answer := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "验证未通过，请一分钟后再试！")
		bot.Arknights.Send(answer)
		chatMember := tgbotapi.ChatMemberConfig{ChatID: chatId, UserID: userId}
		ban(chatId, userId, callbackQuery, chatMember, joinMessageId)
		go unban(chatMember)
		return true, nil
	}

	answer := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "验证通过！")
	bot.Arknights.Send(answer)
	pass(chatId, userId, callbackQuery)

	return true, nil
}

func pass(chatId int64, userId int64, callbackQuery *tgbotapi.CallbackQuery) string {
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
	return val
}

func ban(chatId int64, userId int64, callbackQuery *tgbotapi.CallbackQuery, chatMember tgbotapi.ChatMemberConfig, joinMessageId int) {
	kickChatMemberConfig := tgbotapi.KickChatMemberConfig{
		ChatMemberConfig: chatMember,
	}
	bot.Arknights.Send(kickChatMemberConfig)
	delMsg := tgbotapi.NewDeleteMessage(chatId, callbackQuery.Message.MessageID)
	bot.Arknights.Send(delMsg)
	val := fmt.Sprintf("verify%d%d", chatId, userId)
	utils.RedisDelSetItem("verify", val)
	delJoinMessage := tgbotapi.NewDeleteMessage(chatId, joinMessageId)
	bot.Arknights.Send(delJoinMessage)
}
