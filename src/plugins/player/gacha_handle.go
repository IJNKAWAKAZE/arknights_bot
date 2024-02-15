package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

// Gacha 抽卡记录

func Gacha(uid string, account account.UserAccount, chatId int64, messageId int) (bool, error) {
	var userGacha []UserGacha
	res := utils.GetUserGacha(account.UserNumber, uid).Scan(&userGacha)
	if res.RowsAffected == 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "不存在抽卡记录，请先同步！")
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return true, nil
	}

	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	bot.Arknights.Send(sendAction)

	port := viper.GetString("http.port")
	pic := utils.Screenshot(fmt.Sprintf("http://localhost:%s/gacha?userId=%d&uid=%s", port, account.UserNumber, uid), 3000)
	if pic == nil {
		sendMessage := tgbotapi.NewMessage(chatId, "生成图片失败，token可能已失效请重设token。")
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return true, nil
	}
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	sendPhoto.ReplyToMessageID = messageId
	bot.Arknights.Send(sendPhoto)
	return true, nil
}
