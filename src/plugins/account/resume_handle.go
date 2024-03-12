package account

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	"arknights_bot/utils/telebot"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"unicode/utf8"
)

var setResume = make(map[int64]string)

func ResumeHandle(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	var players []UserPlayer
	res := utils.GetPlayersByUserId(userId).Scan(&players)
	if res.RowsAffected == 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "您还未绑定任何角色！")
		bot.Arknights.Send(sendMessage)
		return true, nil
	}
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, player := range players {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s(%s)", player.PlayerName, player.ServerName), fmt.Sprintf("%s,%s", "resume", player.Uid)),
		))
	}
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)
	sendMessage := tgbotapi.NewMessage(chatId, "请选择要设置的角色")
	sendMessage.ReplyMarkup = inlineKeyboardMarkup
	bot.Arknights.Send(sendMessage)
	return true, nil
}

func WaitResume(chatId, userId int64, uid string) {
	setResume[userId] = uid
	sendMessage := tgbotapi.NewMessage(chatId, "请输入签名(最多30字)，输入null设置签名为空 /cancel 指令取消操作。")
	bot.Arknights.Send(sendMessage)
	telebot.WaitMessage[chatId] = "resume"
}

// Resume 设置签名
func Resume(update tgbotapi.Update) (bool, error) {
	message := update.Message
	chatId := message.Chat.ID
	userId := message.From.ID
	resume := message.Text

	if utf8.RuneCountInString(resume) > 30 {
		sendMessage := tgbotapi.NewMessage(chatId, "超出签名长度限制")
		bot.Arknights.Send(sendMessage)
		return true, nil
	}
	if resume == "null" {
		resume = ""
	}
	bot.DBEngine.Exec("update user_player set resume = ? where user_number = ? and uid = ?", resume, userId, setResume[userId])
	sendMessage := tgbotapi.NewMessage(chatId, "角色签名设置成功！")
	bot.Arknights.Send(sendMessage)
	delete(telebot.WaitMessage, chatId)
	delete(setResume, userId)
	return true, nil
}
