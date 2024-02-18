package system

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"strconv"
)

// HeadhuntHandle 寻访模拟
func HeadhuntHandle(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	param := update.Message.CommandArguments()
	headhuntKey := fmt.Sprintf("headhuntFlag:%d", chatId)
	messagecleaner.AddDelQueue(chatId, messageId, 60)

	if param == "" {
		if utils.RedisIsExists(headhuntKey) && utils.RedisGet(headhuntKey) == "stop" {
			sendMessage := tgbotapi.NewMessage(chatId, "模拟寻访功能已关闭！")
			msg, _ := bot.Arknights.Send(sendMessage)
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
			return true, nil
		}
	}

	if param != "" {
		if utils.IsAdmin(chatId, userId) {
			text := ""
			if param == "start" {
				utils.RedisSet(headhuntKey, "start", 0)
				text = "模拟寻访已开启！"
			} else if param == "stop" {
				utils.RedisSet(headhuntKey, "stop", 0)
				text = "模拟寻访已关闭！"
			}
			sendMessage := tgbotapi.NewMessage(chatId, text)
			msg, _ := bot.Arknights.Send(sendMessage)
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
			return true, nil
		}
		sendMessage := tgbotapi.NewMessage(chatId, "无使用权限！")
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return true, nil
	}

	if !update.Message.Chat.IsPrivate() {
		key := fmt.Sprintf("headhuntTimes:%d", userId)
		if !utils.RedisIsExists(key) {
			utils.RedisSet(key, "1", 0)
		} else {
			times, _ := strconv.Atoi(utils.RedisGet(key))
			headhuntTimes := bot.HeadhuntTimes
			if times == headhuntTimes {
				sendMessage := tgbotapi.NewMessage(chatId, "已达到每日次数限制！")
				sendMessage.ReplyToMessageID = messageId
				msg, _ := bot.Arknights.Send(sendMessage)
				messagecleaner.AddDelQueue(chatId, msg.MessageID, 60)
				return true, nil
			}
			utils.RedisSet(key, strconv.Itoa(times+1), 0)
		}
	}

	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	bot.Arknights.Send(sendAction)
	port := viper.GetString("http.port")
	pic := utils.Screenshot(fmt.Sprintf("http://localhost:%s/headhunt?userId=%d", port, userId), 0)
	if pic == nil {
		sendMessage := tgbotapi.NewMessage(chatId, "失败啦，啊哈哈哈！")
		sendMessage.ReplyToMessageID = messageId
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(chatId, msg.MessageID, 5)
		return true, nil
	}
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	sendPhoto.ReplyToMessageID = messageId
	msg, _ := bot.Arknights.Send(sendPhoto)
	//messagecleaner.AddDelQueue(chatId, msg.MessageID, 60)
	return true, nil
}
func ResetHeadhuntTimes() func() {
	resetHeadhuntTimes := func() {
		res, ctx := utils.RedisScanKeys("headhuntTimes:*")
		for res.Next(ctx) {
			utils.RedisDel(res.Val())
		}
	}
	return resetHeadhuntTimes
}
