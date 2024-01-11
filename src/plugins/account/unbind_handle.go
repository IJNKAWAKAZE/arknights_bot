package account

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// UnbindHandle 解绑角色
func UnbindHandle(update tgbotapi.Update) (bool, error) {
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
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s(%s)", player.PlayerName, player.ServerName), fmt.Sprintf("%s,%s", "unbind", player.Uid)),
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
