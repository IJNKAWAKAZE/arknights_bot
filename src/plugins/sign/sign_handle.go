package sign

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"strings"
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
		msg, err := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(chatId, messageId, 5)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	// 获取绑定角色
	res = utils.GetPlayersByUserId(userId).Scan(&players)
	if res.RowsAffected == 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "您还未绑定任何角色！")
		msg, err := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(chatId, messageId, 5)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
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
		case "ap_on":
			// 开启理智提醒
			setApRemind(update, 1)
		case "ap_off":
			// 关闭理智提醒
			setApRemind(update, 0)
		default:
			// 尝试解析 ap_threshold 参数
			if strings.HasPrefix(param, "ap_threshold ") {
				thresholdStr := strings.TrimPrefix(param, "ap_threshold ")
				threshold, err := strconv.Atoi(thresholdStr)
				if err != nil || threshold < 1 || threshold > 100 {
					sendMessage := tgbotapi.NewMessage(chatId, "理智提醒阈值请输入1-100之间的整数！")
					sendMessage.ReplyToMessageID = messageId
					msg, err := bot.Arknights.Send(sendMessage)
					if err != nil {
						return err
					}
					messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
					return nil
				}
				setApThreshold(update, threshold)
			} else {
				sendMessage := tgbotapi.NewMessage(chatId, "未知的签到指令参数，请使用 /help 查看使用说明。")
				sendMessage.ReplyToMessageID = messageId
				msg, err := bot.Arknights.Send(sendMessage)
				if err != nil {
					return err
				}
				messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
			}
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
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
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
		log.Println(playerName, err)
		sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("角色 %s 签到失败！\n失败原因:%s", playerName, err.Error()))
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
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
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return
	}
	id, _ := gonanoid.New(32)
	userSign = UserSign{
		Id:          id,
		UserName:    message.From.FullName(),
		UserNumber:  userId,
		NotifyMode:  0,
		ApRemind:    0,
		ApThreshold: defaultApThreshold,
		ApNotified:  0,
	}

	bot.DBEngine.Table("user_sign").Create(&userSign)

	sendMessage := tgbotapi.NewMessage(chatId, "开启自动签到成功！")
	sendMessage.ReplyToMessageID = messageId
	msg, err := bot.Arknights.Send(sendMessage)
	if err != nil {
		return
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
}

// 关闭自动签到
func stopSign(update tgbotapi.Update) {
	message := update.Message
	userId := message.From.ID
	chatId := message.Chat.ID
	messageId := message.MessageID

	CancelApCheck(userId)
	bot.DBEngine.Exec("delete from user_sign where user_number = ?", userId)

	sendMessage := tgbotapi.NewMessage(chatId, "已关闭自动签到！")
	sendMessage.ReplyToMessageID = messageId
	msg, err := bot.Arknights.Send(sendMessage)
	if err != nil {
		return
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
}

// 设置签到通知模式
func setNotifyMode(update tgbotapi.Update, mode int) {
	message := update.Message
	userId := message.From.ID
	chatId := message.Chat.ID
	messageId := message.MessageID

	var userSign UserSign
	res := utils.GetAutoSignByUserId(userId).Scan(&userSign)
	if res.RowsAffected == 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "请先开启自动签到！(/sign auto)")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return
	}

	bot.DBEngine.Exec("update user_sign set notify_mode = ? where user_number = ?", mode, userId)

	var modeText string
	switch mode {
	case 0:
		modeText = "全部通知"
	case 1:
		modeText = "仅失败时通知"
	case 2:
		modeText = "仅成功时通知"
	}

	sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("签到通知模式已设置为：%s", modeText))
	sendMessage.ReplyToMessageID = messageId
	msg, err := bot.Arknights.Send(sendMessage)
	if err != nil {
		return
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
}

// 设置理智提醒开关
func setApRemind(update tgbotapi.Update, status int) {
	message := update.Message
	userId := message.From.ID
	chatId := message.Chat.ID
	messageId := message.MessageID

	var userSign UserSign
	res := utils.GetAutoSignByUserId(userId).Scan(&userSign)
	if res.RowsAffected == 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "请先开启自动签到！(/sign auto)")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return
	}

	bot.DBEngine.Exec("update user_sign set ap_remind = ?, ap_notified = 0 where user_number = ?", status, userId)

	var text string
	if status == 1 {
		displayThreshold := userSign.ApThreshold
		if displayThreshold == 0 {
			displayThreshold = defaultApThreshold
		}
		text = fmt.Sprintf("理智提醒已开启！当理智恢复到 %d%% 时将发送通知。", displayThreshold)
		ScheduleNextApCheck(userId)
	} else {
		text = "理智提醒已关闭！"
		CancelApCheck(userId)
	}

	sendMessage := tgbotapi.NewMessage(chatId, text)
	sendMessage.ReplyToMessageID = messageId
	msg, err := bot.Arknights.Send(sendMessage)
	if err != nil {
		return
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
}

// 设置理智提醒阈值
func setApThreshold(update tgbotapi.Update, threshold int) {
	message := update.Message
	userId := message.From.ID
	chatId := message.Chat.ID
	messageId := message.MessageID

	var userSign UserSign
	res := utils.GetAutoSignByUserId(userId).Scan(&userSign)
	if res.RowsAffected == 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "请先开启自动签到！(/sign auto)")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return
	}

	bot.DBEngine.Exec("update user_sign set ap_threshold = ?, ap_notified = 0 where user_number = ?", threshold, userId)

	// Reschedule so the new threshold is used for the next check.
	if userSign.ApRemind == 1 {
		ScheduleNextApCheck(userId)
	}

	sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("理智提醒阈值已设置为 %d%%", threshold))
	sendMessage.ReplyToMessageID = messageId
	msg, err := bot.Arknights.Send(sendMessage)
	if err != nil {
		return
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
}
