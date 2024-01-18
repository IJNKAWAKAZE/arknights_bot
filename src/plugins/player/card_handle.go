package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

// CardHandle 我的名片
func CardHandle(players []account.UserPlayer, userAccount account.UserAccount, chatId int64, userId int64, messageId int) (bool, error) {
	if len(players) > 1 {
		// 绑定多个角色进行选择
		var buttons [][]tgbotapi.InlineKeyboardButton
		for _, player := range players {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s(%s)", player.PlayerName, player.ServerName), fmt.Sprintf("%s,%s,%d,%s,%d", "player", OP_CARD, userId, player.Uid, messageId)),
			))
		}
		inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
			buttons...,
		)
		sendMessage := tgbotapi.NewMessage(chatId, "请选择要查询的角色")
		sendMessage.ReplyMarkup = inlineKeyboardMarkup
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
	} else {
		// 绑定单个角色
		return Card(players[0].Uid, userAccount, chatId, messageId)
	}

	return true, nil
}

func Card(uid string, account account.UserAccount, chatId int64, messageId int) (bool, error) {

	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	bot.Arknights.Send(sendAction)

	port := viper.GetString("http.port")
	pic := utils.Screenshot(fmt.Sprintf("http://localhost:%s/card?userId=%d&uid=%s", port, account.UserNumber, uid), 0)
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
