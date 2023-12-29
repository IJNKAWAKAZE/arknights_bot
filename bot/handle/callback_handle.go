package handle

import (
	"arknights_bot/bot/utils"
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
		memberInfo, _ := utils.GetChatMemberInfo(getChatMemberConfig)

		if memberInfo.Status != "creator" && memberInfo.Status != "administrator" {
			answer := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "无使用权限！")
			utils.SendCallbackAnswer(answer)
			return true, nil
		}

		if d[1] == "PASS" {
			pass(chatId, userId, callbackQuery, "管理员通过了<a href=\"tg://user?id="+strconv.FormatInt(userId, 10)+"\">"+d[2]+"</a>的验证")
		}

		if d[1] == "BAN" {
			chatMember := tgbotapi.ChatMemberConfig{ChatID: chatId, UserID: userId}
			ban(chatId, userId, callbackQuery, chatMember, "管理员封禁了<a href=\"tg://user?id="+strconv.FormatInt(userId, 10)+"\">"+d[2]+"</a>")
		}

		return true, nil
	}

	if userId != callbackQuery.From.ID {
		answer := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "这不是你的验证！")
		utils.SendCallbackAnswer(answer)
		return true, nil
	}

	if d[1] != d[2] {
		answer := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "验证未通过！")
		utils.SendCallbackAnswer(answer)
		chatMember := tgbotapi.ChatMemberConfig{ChatID: chatId, UserID: userId}
		ban(chatId, userId, callbackQuery, chatMember, "<a href=\"tg://user?id="+strconv.FormatInt(userId, 10)+"\">"+name+"</a>未通过验证，已被踢出。")
		go unban(chatMember)
		return true, nil
	}

	answer := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "验证通过！")
	utils.SendCallbackAnswer(answer)
	pass(chatId, userId, callbackQuery, "<a href=\"tg://user?id="+strconv.FormatInt(userId, 10)+"\">"+name+"</a>已完成验证")

	return true, nil
}

func pass(chatId int64, userId int64, callbackQuery *tgbotapi.CallbackQuery, text string) string {
	utils.SetMemberPermissions(tgbotapi.RestrictChatMemberConfig{
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
	utils.DeleteMessage(delMsg)
	sendMessage := tgbotapi.NewMessage(chatId, text)
	sendMessage.ParseMode = tgbotapi.ModeHTML
	msg, _ := utils.SendMessage(sendMessage)
	utils.AddDelQueue(msg.Chat.ID, msg.MessageID, 1)
	return val
}

func ban(chatId int64, userId int64, callbackQuery *tgbotapi.CallbackQuery, chatMember tgbotapi.ChatMemberConfig, text string) {
	kickChatMemberConfig := tgbotapi.KickChatMemberConfig{
		ChatMemberConfig: chatMember,
	}
	utils.KickChatMember(kickChatMemberConfig)
	delMsg := tgbotapi.NewDeleteMessage(chatId, callbackQuery.Message.MessageID)
	utils.DeleteMessage(delMsg)
	cid := strconv.FormatInt(chatId, 10)
	uid := strconv.FormatInt(userId, 10)
	val := "verify" + cid + uid
	utils.RedisDelSetItem("verify", val)
	sendMessage := tgbotapi.NewMessage(chatId, text)
	sendMessage.ParseMode = tgbotapi.ModeHTML
	msg, _ := utils.SendMessage(sendMessage)
	utils.AddDelQueue(msg.Chat.ID, msg.MessageID, 1)
}
