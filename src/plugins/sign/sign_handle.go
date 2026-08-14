package sign

import (
	"arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils/repo"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
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

	res := repo.GetAccountByUserId(userId).Scan(&userAccount)
	if res.RowsAffected == 0 {
		// 未绑定账号
		msg, err := config.Arknights.SendMarkdownV2(chatId, fmt.Sprintf("未查询到绑定账号，请先进行[绑定](https://t.me/%s)。", viper.GetString("bot.name")), messageId)
		messagecleaner.AddDelQueue(chatId, messageId, 5)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}

	// 获取绑定角色
	res = repo.GetPlayersByUserId(userId).Scan(&players)
	if res.RowsAffected == 0 {
		msg, err := config.Arknights.SendText(chatId, "您还未绑定任何角色！")
		messagecleaner.AddDelQueue(chatId, messageId, 5)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}

	if param != "" {
		switch param {
		case "auto":
			// 开启自动签到
			autoSign(update)
		case "stop":
			// 关闭自动签到
			stopSign(update)
		case "notify_all":
			// 设置通知模式为全部通知
			setNotifyMode(update, 0)
		case "notify_fail":
			// 设置通知模式为仅失败通知
			setNotifyMode(update, 1)
		case "notify_success":
			// 设置通知模式为仅成功通知
			setNotifyMode(update, 2)
		default:
			sendMessage := tgbotapi.NewMessage(chatId, "未知的签到指令参数，请使用 /help 查看使用说明。")
			sendMessage.ReplyToMessageID = messageId
			msg, err := config.Arknights.Send(sendMessage)
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
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
		msg, err := config.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
	} else {
		// 绑定单个角色执行签到
		repo.GetAccountByUid(userId, players[0].Uid).Scan(&userAccount)
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

	_, _ = config.Arknights.SendChatAction(chatId, "typing")

	award, hasSigned, err := skland.SignGamePlayer(player.Uid, skAccount, account.ServerName)
	if err != nil {
		log.Println(playerName, err)
		msg, err := config.Arknights.SendText(chatId, fmt.Sprintf("角色 %s 签到失败！\n失败原因:%s", playerName, err.Error()))
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}
	// 今日已完成签到
	if hasSigned {
		config.Arknights.SendText(chatId, fmt.Sprintf("角色 %s 今天已经签到过了", playerName))
		return nil
	}
	// 签到成功
	config.Arknights.SendText(chatId, fmt.Sprintf("角色 %s 签到成功!\n今日奖励：%s", playerName, award))
	return nil
}

// 开启自动签到
func autoSign(update tgbotapi.Update) {
	message := update.Message
	userId := message.From.ID
	chatId := message.Chat.ID
	messageId := message.MessageID
	var userSign UserSign
	res := repo.GetAutoSignByUserId(userId).Scan(&userSign)
	if res.RowsAffected > 0 {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "已开启自动签到！")
		if err != nil {
			return
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return
	}
	id, _ := gonanoid.New(32)
	userSign = UserSign{
		Id:         id,
		UserName:   message.From.FullName(),
		UserNumber: userId,
		NotifyMode: 0,
	}

	config.DBEngine.Table("user_sign").Create(&userSign)

	msg, err := config.Arknights.ReplyText(chatId, messageId, "开启自动签到成功！")
	if err != nil {
		return
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
}

// 关闭自动签到
func stopSign(update tgbotapi.Update) {
	message := update.Message
	userId := message.From.ID
	chatId := message.Chat.ID
	messageId := message.MessageID

	config.DBEngine.Exec("delete from user_sign where user_number = ?", userId)

	msg, err := config.Arknights.ReplyText(chatId, messageId, "已关闭自动签到！")
	if err != nil {
		return
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
}

// 设置签到通知模式
func setNotifyMode(update tgbotapi.Update, mode int) {
	message := update.Message
	userId := message.From.ID
	chatId := message.Chat.ID
	messageId := message.MessageID

	var userSign UserSign
	res := repo.GetAutoSignByUserId(userId).Scan(&userSign)
	if res.RowsAffected == 0 {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "请先开启自动签到！(/sign auto)")
		if err != nil {
			return
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return
	}

	config.DBEngine.Exec("update user_sign set notify_mode = ? where user_number = ?", mode, userId)

	var modeText string
	switch mode {
	case 0:
		modeText = "全部通知"
	case 1:
		modeText = "仅失败时通知"
	case 2:
		modeText = "仅成功时通知"
	}

	msg, err := config.Arknights.ReplyText(chatId, messageId, fmt.Sprintf("签到通知模式已设置为：%s", modeText))
	if err != nil {
		return
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
}
