package account

import (
	"arknights_bot/config"
	"arknights_bot/utils/repo"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

// UnbindHandle 解绑角色
func UnbindHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	var players []UserPlayer
	res := repo.GetPlayersByUserId(userId).Scan(&players)
	if res.RowsAffected == 0 {
		config.Arknights.SendText(chatId, "您还未绑定任何角色！")
		return nil
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
	config.Arknights.Send(sendMessage)
	return nil
}
