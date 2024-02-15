package system

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/datasource"
	"arknights_bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

func UpdateHandle(update tgbotapi.Update) (bool, error) {
	owner := viper.GetInt64("bot.owner")
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID

	getChatMemberConfig := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatId,
			UserID: userId,
		},
	}

	if utils.IsAdmin(getChatMemberConfig) || owner == userId {
		sendMessage := tgbotapi.NewMessage(chatId, "开始更新数据源")
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		datasource.UpdateDataSourceRunner()
		sendMessage = tgbotapi.NewMessage(chatId, "数据源更新结束")
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return true, nil
	}

	sendMessage := tgbotapi.NewMessage(chatId, "无使用权限！")
	sendMessage.ReplyToMessageID = messageId
	bot.Arknights.Send(sendMessage)
	return true, nil
}
