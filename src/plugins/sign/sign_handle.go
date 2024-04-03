package sign

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/spf13/viper"
	"log"
)

// SignHandle 森空岛签到
func SignHandle(update tgbotapi.Update) error {
	param := update.Message.CommandArguments()
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID

	var userAccount account.UserAccount
	var players []account.UserPlayer

	res := utils.GetAccountByUserId(userId).Scan(&userAccount)
	if res.RowsAffected == 0 {
		// 未绑定账号
		sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("未查询到绑定账号，请先进行[绑定](https://t.me/%s)。", viper.GetString("bot.name")))
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		sendMessage.ReplyToMessageID = messageId
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(chatId, messageId, 5)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	// 获取绑定角色
	res = utils.GetPlayersByUserId(userId).Scan(&players)
	if res.RowsAffected == 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "您还未绑定任何角色！")
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(chatId, messageId, 5)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	if param != "" {
		if param == "auto" {
			// 开启自动签到
			autoSign(update)
		} else if param == "stop" {
			// 关闭自动签到
			stopSign(update)
		}
		return nil
	}

	if res.RowsAffected > 1 {
		// 绑定多个角色进行选择
		var buttons [][]tgbotapi.InlineKeyboardButton
		for _, player := range players {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s(%s)", player.PlayerName, player.ServerName), fmt.Sprintf("%s,%d,%s", "sign", userId, player.Uid)),
			))
		}
		inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
			buttons...,
		)
		sendMessage := tgbotapi.NewMessage(chatId, "请选择要签到的角色")
		sendMessage.ReplyMarkup = inlineKeyboardMarkup
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
	} else {
		// 绑定单个角色执行签到
		utils.GetAccountByUid(userId, players[0].Uid).Scan(&userAccount)
		return Sign(players[0], userAccount, chatId)
	}
	return nil
}

func Sign(player account.UserPlayer, account account.UserAccount, chatId int64) error {
	var skAccount skland.Account
	playerName := player.PlayerName
	skAccount.Hypergryph.Token = account.HypergryphToken
	skAccount.Skland.Token = account.SklandToken
	skAccount.Skland.Cred = account.SklandCred

	sendAction := tgbotapi.NewChatAction(chatId, "typing")
	bot.Arknights.Send(sendAction)

	award, hasSigned, err := skland.SignGamePlayer(player.Uid, skAccount)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("角色 %s 签到失败！\nmsg:%s", playerName, err.Error()))
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		log.Println(playerName, err)
		return err
	}
	// 今日已完成签到
	if hasSigned {
		sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("角色 %s 今天已经签到过了", playerName))
		bot.Arknights.Send(sendMessage)
		return nil
	}
	// 签到成功
	sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("角色 %s 签到成功!\n今日奖励：%s", playerName, award))
	bot.Arknights.Send(sendMessage)
	return nil
}

// 开启自动签到
func autoSign(update tgbotapi.Update) {
	message := update.Message
	userId := message.From.ID
	chatId := message.Chat.ID
	messageId := message.MessageID
	var userSign UserSign
	res := utils.GetAutoSignByUserId(userId).Scan(&userSign)
	if res.RowsAffected > 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "已开启自动签到！")
		sendMessage.ReplyToMessageID = messageId
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return
	}
	id, _ := gonanoid.New(32)
	userSign = UserSign{
		Id:         id,
		UserName:   utils.GetFullName(message.From),
		UserNumber: userId,
	}

	bot.DBEngine.Table("user_sign").Create(&userSign)

	sendMessage := tgbotapi.NewMessage(chatId, "开启自动签到成功！")
	sendMessage.ReplyToMessageID = messageId
	msg, _ := bot.Arknights.Send(sendMessage)
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
}

// 关闭自动签到
func stopSign(update tgbotapi.Update) {
	message := update.Message
	userId := message.From.ID
	chatId := message.Chat.ID
	messageId := message.MessageID

	bot.DBEngine.Exec("delete from user_sign where user_number = ?", userId)

	sendMessage := tgbotapi.NewMessage(chatId, "已关闭自动签到！")
	sendMessage.ReplyToMessageID = messageId
	msg, _ := bot.Arknights.Send(sendMessage)
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
}
