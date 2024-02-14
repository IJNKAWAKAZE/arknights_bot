package system

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
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
	chatType := update.Message.Chat.Type
	if chatType != "private" {
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
		return true, nil
	}
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	sendPhoto.ReplyToMessageID = messageId
	msg, _ := bot.Arknights.Send(sendPhoto)
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
