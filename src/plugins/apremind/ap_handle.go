package apremind

import (
	"arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils/repo"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

// ApHandle 理智提醒
func ApHandle(update tgbotapi.Update) error {
	param := update.Message.CommandArguments()
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID

	var userAccount account.UserAccount

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

	if param == "" {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "请指定理智提醒参数，使用 /help 查看使用说明。")
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}

	switch param {
	case "on":
		// 开启理智提醒
		apRemindOn(update)
	case "off":
		// 关闭理智提醒
		apRemindOff(update)
	default:
		// 尝试解析 thr 参数
		if strings.HasPrefix(param, "thr ") {
			thresholdStr := strings.TrimPrefix(param, "thr ")
			threshold, err := strconv.Atoi(thresholdStr)
			if err != nil || threshold < 1 || threshold > 100 {
				msg, err := config.Arknights.ReplyText(chatId, messageId, "理智提醒阈值请输入1-100之间的整数！")
				if err != nil {
					return err
				}
				messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
				return nil
			}
			apSetThreshold(update, threshold)
		} else {
			msg, err := config.Arknights.ReplyText(chatId, messageId, "未知的理智提醒参数，请使用 /help 查看使用说明。")
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		}
	}
	return nil
}

// 开启理智提醒
func apRemindOn(update tgbotapi.Update) {
	message := update.Message
	userId := message.From.ID
	chatId := message.Chat.ID
	messageId := message.MessageID

	var userApRemind UserApRemind
	res := repo.GetApRemindByUserId(userId).Scan(&userApRemind)
	if res.RowsAffected > 0 {
		displayThreshold := userApRemind.ApThreshold
		if displayThreshold == 0 {
			displayThreshold = defaultApThreshold
		}
		msg, err := config.Arknights.ReplyText(chatId, messageId, fmt.Sprintf("已开启理智提醒！当前阈值 %d%%。", displayThreshold))
		if err != nil {
			return
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return
	}

	id, _ := gonanoid.New(32)
	userApRemind = UserApRemind{
		Id:          id,
		UserName:    message.From.FullName(),
		UserNumber:  userId,
		ApThreshold: defaultApThreshold,
		ApNotified:  0,
	}

	config.DBEngine.Table("user_ap_remind").Create(&userApRemind)

	ScheduleNextApCheck(userId)

	msg, err := config.Arknights.ReplyText(chatId, messageId, fmt.Sprintf("理智提醒已开启！当理智恢复到 %d%% 时将发送通知。", defaultApThreshold))
	if err != nil {
		return
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
}

// 关闭理智提醒
func apRemindOff(update tgbotapi.Update) {
	message := update.Message
	userId := message.From.ID
	chatId := message.Chat.ID
	messageId := message.MessageID

	CancelApCheck(userId)
	config.DBEngine.Exec("delete from user_ap_remind where user_number = ?", userId)

	msg, err := config.Arknights.ReplyText(chatId, messageId, "理智提醒已关闭！")
	if err != nil {
		return
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
}

// 设置理智提醒阈值
func apSetThreshold(update tgbotapi.Update, threshold int) {
	message := update.Message
	userId := message.From.ID
	chatId := message.Chat.ID
	messageId := message.MessageID

	var userApRemind UserApRemind
	res := repo.GetApRemindByUserId(userId).Scan(&userApRemind)
	if res.RowsAffected == 0 {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "请先开启理智提醒！(/ap on)")
		if err != nil {
			return
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return
	}

	config.DBEngine.Exec("update user_ap_remind set ap_threshold = ?, ap_notified = 0 where user_number = ?", threshold, userId)

	// Reschedule so the new threshold is used for the next check.
	ScheduleNextApCheck(userId)

	msg, err := config.Arknights.ReplyText(chatId, messageId, fmt.Sprintf("理智提醒阈值已设置为 %d%%", threshold))
	if err != nil {
		return
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
}
