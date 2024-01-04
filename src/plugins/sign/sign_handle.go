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
	"strings"
)

// SignHandle 森空岛签到
func SignHandle(update tgbotapi.Update) (bool, error) {
	cmd := strings.Split(update.Message.Text, " ")
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID

	messagecleaner.AddDelQueue(chatId, messageId, 5)

	var userAccount account.UserAccount
	var players []account.UserPlayer
	var skPlayer skland.Player
	var skAccount skland.Account

	res := utils.GetAccountByUserId(userId).Scan(&userAccount)
	if res.RowsAffected == 0 {
		// 未绑定账号
		sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("未查询到绑定账号，请先进行[绑定](https://t.me/%s)。", viper.GetString("bot.name")))
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return true, nil
	}

	if len(cmd) > 1 {
		param := cmd[1]
		if param == "auto" {
			// 开启自动签到
			autoSign(update)
		} else if param == "stop" {
			// 关闭自动签到
			stopSign(update)
		}
		return true, nil
	}

	// 获取绑定角色
	res = utils.GetPlayersByUserId(userId).Scan(&players)
	if res.RowsAffected == 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "您还未绑定任何角色！")
		bot.Arknights.Send(sendMessage)
		return true, nil
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
		skPlayer.NickName = players[0].PlayerName
		skPlayer.ChannelName = players[0].ServerName
		skPlayer.Uid = players[0].Uid
		skAccount.Hypergryph.Token = userAccount.HypergryphToken
		skAccount.Skland.Token = userAccount.SklandToken
		skAccount.Skland.Cred = userAccount.SklandCred
		record, err := skland.SignGamePlayer(&skPlayer, skAccount)
		if err != nil {
			return true, err
		}
		// 今日已完成签到
		if record.HasSigned {
			sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("角色 %s 今天已经签到过了", players[0].PlayerName))
			msg, _ := bot.Arknights.Send(sendMessage)
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
			return true, nil
		}
		// 签到成功
		sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("角色 %s 签到成功!\n今日奖励：%s", players[0].PlayerName, record.Award))
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
	}
	return true, nil
}

// 开启自动签到
func autoSign(update tgbotapi.Update) {
	message := update.Message
	userId := message.From.ID
	chatId := message.Chat.ID
	var userSign UserSign
	res := utils.GetAutoSignByUserId(userId).Scan(&userSign)
	if res.RowsAffected > 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "已开启自动签到！")
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
	msg, _ := bot.Arknights.Send(sendMessage)
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
}

// 关闭自动签到
func stopSign(update tgbotapi.Update) {
	message := update.Message
	userId := message.From.ID
	chatId := message.Chat.ID

	bot.DBEngine.Exec("delete from user_sign where user_number = ?", userId)

	sendMessage := tgbotapi.NewMessage(chatId, "已关闭自动签到！")
	msg, _ := bot.Arknights.Send(sendMessage)
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
}
