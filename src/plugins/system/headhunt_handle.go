package system

import (
	"arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils/cache"
	"arknights_bot/utils/media"
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
		if cache.RedisIsExists(headhuntKey) && cache.RedisGet(headhuntKey) == "off" {
			sendMessage := tgbotapi.NewMessage(chatId, "模拟寻访功能已关闭！")
			msg, err := config.Arknights.Send(sendMessage)
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
			return nil
		}
	}

	if param != "" {
		if config.Arknights.IsAdmin(chatId, userId) {
			text := ""
			if param == "on" {
				cache.RedisSet(headhuntKey, "on", 0)
				text = "模拟寻访已开启！"
			} else if param == "off" {
				cache.RedisSet(headhuntKey, "off", 0)
				text = "模拟寻访已关闭！"
			}
			sendMessage := tgbotapi.NewMessage(chatId, text)
			msg, err := config.Arknights.Send(sendMessage)
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
			return nil
		}
		sendMessage := tgbotapi.NewMessage(chatId, "无使用权限！")
		msg, err := config.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}

	key := fmt.Sprintf("headhuntTimes:%d", userId)
	if !update.Message.Chat.IsPrivate() {
		if !cache.RedisIsExists(key) {
			cache.RedisSet(key, "1", 0)
		} else {
			times, _ := strconv.Atoi(cache.RedisGet(key))
			headhuntTimes := config.HeadhuntTimes
			if times == headhuntTimes {
				messagecleaner.AddDelQueue(chatId, messageId, 60)
				sendMessage := tgbotapi.NewMessage(chatId, "已达到每日次数限制！")
				sendMessage.ReplyToMessageID = messageId
				msg, err := config.Arknights.Send(sendMessage)
				if err != nil {
					return err
				}
				messagecleaner.AddDelQueue(chatId, msg.MessageID, 60)
				return nil
			}
			cache.RedisSet(key, strconv.Itoa(times+1), 0)
		}
	}

	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	config.Arknights.Send(sendAction)
	port := viper.GetString("http.port")
	pic, err := media.Screenshot(fmt.Sprintf("http://localhost:%s/headhunt?userId=%d", port, userId), 0, 1)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, err.Error())
		sendMessage.ReplyToMessageID = messageId
		msg, err := config.Arknights.Send(sendMessage)
		times, _ := strconv.Atoi(cache.RedisGet(key))
		cache.RedisSet(key, strconv.Itoa(times-1), 0)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(chatId, msg.MessageID, 5)
		return nil
	}
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	sendPhoto.ReplyToMessageID = messageId
	msg, err := config.Arknights.Send(sendPhoto)
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
	res, ctx := cache.RedisScanKeys("headhuntTimes:*")
	for res.Next(ctx) {
		cache.RedisDel(res.Val())
	}
}
