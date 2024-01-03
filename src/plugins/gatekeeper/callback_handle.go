package gatekeeper

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

func CallBackData(callBack tgbotapi.Update) (bool, error) {
	callbackQuery := callBack.CallbackQuery
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 4 {
		return true, nil
	}

	userId, _ := strconv.ParseInt(d[1], 10, 64)
	chatId := callbackQuery.Message.Chat.ID
	name := utils.GetFullName(callbackQuery.From)

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
			pass(chatId, userId, callbackQuery, fmt.Sprintf("管理员通过了<a href=\"tg://user?id=%d\">%s</a>的验证", userId, d[3]))
		}

		if d[2] == "BAN" {
			chatMember := tgbotapi.ChatMemberConfig{ChatID: chatId, UserID: userId}
			ban(chatId, userId, callbackQuery, chatMember, fmt.Sprintf("管理员封禁了<a href=\"tg://user?id=%d\">%s</a>", userId, d[3]))
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
	go utils.DelayDelMsg(msg.Chat.ID, msg.MessageID, viper.GetDuration("bot.msg_del_delay"))
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
	go utils.DelayDelMsg(msg.Chat.ID, msg.MessageID, viper.GetDuration("bot.msg_del_delay"))
}

func ChoosePlayer(callBack tgbotapi.Update) (bool, error) {
	callbackQuery := callBack.CallbackQuery
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 4 {
		return true, nil
	}

	userId := callbackQuery.From.ID
	chatId := callbackQuery.Message.Chat.ID
	messageId := callbackQuery.Message.MessageID

	uid := d[1]
	serverName := d[2]
	playerName := d[3]

	var userAccount account.UserAccount
	var userPlayer account.UserPlayer
	utils.GetAccountByUserId(userId).Scan(&userAccount)
	res := bot.DBEngine.Raw("select * from user_player where user_number = ? and uid = ?", userId, uid).Scan(&userPlayer)
	if res.RowsAffected == 0 {
		id, _ := gonanoid.New(32)
		userPlayer = account.UserPlayer{
			Id:         id,
			AccountId:  userAccount.Id,
			UserName:   userAccount.UserName,
			UserNumber: userAccount.UserNumber,
			Uid:        uid,
			ServerName: serverName,
			PlayerName: playerName,
		}
		bot.DBEngine.Table("user_player").Save(&userPlayer)
	}
	sendMessage := tgbotapi.NewMessage(chatId, "角色绑定成功！")
	bot.Arknights.Send(sendMessage)
	delMsg := tgbotapi.NewDeleteMessage(chatId, messageId)
	bot.Arknights.Send(delMsg)
	return true, nil
}

func UnbindPlayer(callBack tgbotapi.Update) (bool, error) {
	callbackQuery := callBack.CallbackQuery
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 2 {
		return true, nil
	}

	userId := callbackQuery.From.ID
	chatId := callbackQuery.Message.Chat.ID
	messageId := callbackQuery.Message.MessageID

	uid := d[1]
	bot.DBEngine.Exec("delete from user_player where user_number = ? and uid = ?", userId, uid)
	sendMessage := tgbotapi.NewMessage(chatId, "角色解绑成功！")
	bot.Arknights.Send(sendMessage)
	delMsg := tgbotapi.NewDeleteMessage(chatId, messageId)
	bot.Arknights.Send(delMsg)
	return true, nil
}
