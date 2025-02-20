package system

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	"strconv"
)

// HeadhuntHandle 寻访模拟
func HeadhuntHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	param := update.Message.CommandArguments()
	headhuntKey := fmt.Sprintf("headhuntFlag:%d", chatId)

	if param == "" {
		if utils.RedisIsExists(headhuntKey) && utils.RedisGet(headhuntKey) == "stop" {
			sendMessage := tgbotapi.NewMessage(chatId, "模拟寻访功能已关闭！")
			msg, err := bot.Arknights.Send(sendMessage)
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
			return nil
		}
	}

	if param != "" {
		if bot.Arknights.IsAdmin(chatId, userId) {
			text := ""
			if param == "start" {
				utils.RedisSet(headhuntKey, "start", 0)
				text = "模拟寻访已开启！"
			} else if param == "stop" {
				utils.RedisSet(headhuntKey, "stop", 0)
				text = "模拟寻访已关闭！"
			}
			sendMessage := tgbotapi.NewMessage(chatId, text)
			msg, err := bot.Arknights.Send(sendMessage)
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
			return nil
		}
		sendMessage := tgbotapi.NewMessage(chatId, "无使用权限！")
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	key := fmt.Sprintf("headhuntTimes:%d", userId)
	if !update.Message.Chat.IsPrivate() {
		if !utils.RedisIsExists(key) {
			utils.RedisSet(key, "1", 0)
		} else {
			times, _ := strconv.Atoi(utils.RedisGet(key))
			headhuntTimes := bot.HeadhuntTimes
			if times == headhuntTimes {
				messagecleaner.AddDelQueue(chatId, messageId, 60)
				sendMessage := tgbotapi.NewMessage(chatId, "已达到每日次数限制！")
				sendMessage.ReplyToMessageID = messageId
				msg, err := bot.Arknights.Send(sendMessage)
				if err != nil {
					return err
				}
				messagecleaner.AddDelQueue(chatId, msg.MessageID, 60)
				return nil
			}
			utils.RedisSet(key, strconv.Itoa(times+1), 0)
		}
	}

	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	bot.Arknights.Send(sendAction)
	port := viper.GetString("http.port")
	pic, err := utils.Screenshot(fmt.Sprintf("http://localhost:%s/headhunt?userId=%d", port, userId), 0, 1)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, err.Error())
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		times, _ := strconv.Atoi(utils.RedisGet(key))
		utils.RedisSet(key, strconv.Itoa(times-1), 0)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(chatId, msg.MessageID, 5)
		return nil
	}
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	sendPhoto.ReplyToMessageID = messageId
	msg, err := bot.Arknights.Send(sendPhoto)
	if err != nil {
		return err
	}
	if !update.Message.Chat.IsPrivate() {
		messagecleaner.AddDelQueue(chatId, msg.MessageID, 600)
		messagecleaner.AddDelQueue(chatId, messageId, 600)
	}
	return nil
}
func ResetHeadhuntTimes() {
	res, ctx := utils.RedisScanKeys("headhuntTimes:*")
	for res.Next(ctx) {
		utils.RedisDel(res.Val())
	}
}
