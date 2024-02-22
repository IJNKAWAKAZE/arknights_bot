package account

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"arknights_bot/utils/telebot"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bToken = make(map[int64]string)

// SetBTokenHandle 设置btoken
func SetBTokenHandle(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID

	var players []UserPlayer
	res := utils.GetBPlayersByUserId(userId).Scan(&players)
	if res.RowsAffected == 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "您还未绑定B服角色！")
		bot.Arknights.Send(sendMessage)
		return true, nil
	}
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, player := range players {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s(%s)", player.PlayerName, player.ServerName), fmt.Sprintf("%s,%s", "setbtoken", player.Uid)),
		))
	}
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)
	sendMessage := tgbotapi.NewMessage(chatId, "请选择要解绑的角色")
	sendMessage.ReplyMarkup = inlineKeyboardMarkup
	bot.Arknights.Send(sendMessage)
	return true, nil
}

func WaitBToken(chatId, userId int64, uid string) {
	bToken[userId] = uid
	sendMessage := tgbotapi.NewMessage(chatId, "请输入btoken或使用 /cancel 指令取消操作。")
	bot.Arknights.Send(sendMessage)
	sendMessage.Text = "如何获取token\n\n" +
		"1\\.前往 [官网](https://ak.hypergryph.com/user/bilibili/login) 登录\n" +
		"2\\.打开网址复制content中的 token  [获取token](https://web-api.hypergryph.com/account/info/ak-b)"
	sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
	bot.Arknights.Send(sendMessage)
	telebot.WaitMessage[chatId] = "bToken"
}

// SetBToken 设置btoken
func SetBToken(update tgbotapi.Update) (bool, error) {
	message := update.Message
	chatId := message.Chat.ID
	userId := message.From.ID
	token := message.Text

	sendAction := tgbotapi.NewChatAction(chatId, "typing")
	bot.Arknights.Send(sendAction)

	err := skland.CheckBToken(token)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, "请检查token是否正确。")
		bot.Arknights.Send(sendMessage)
		return true, err
	}

	bot.DBEngine.Exec("update user_player set b_token = ? where user_number = ? and uid = ?", token, userId, bToken[userId])
	sendMessage := tgbotapi.NewMessage(chatId, "角色BToken设置成功！")
	bot.Arknights.Send(sendMessage)
	delete(telebot.WaitMessage, chatId)
	delete(bToken, userId)
	return true, nil
}
