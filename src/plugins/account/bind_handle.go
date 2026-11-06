package account

import (
	"arknights_bot/config"
	"arknights_bot/plugins/skland"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

// BindHandle 绑定角色
func BindHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	var buttons [][]tgbotapi.InlineKeyboardButton
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("国服", fmt.Sprintf("%s,%s,%s", "chooseServer", "国服", "setToken")),
		tgbotapi.NewInlineKeyboardButtonData("国际服", fmt.Sprintf("%s,%s,%s", "chooseServer", "国际服", "setToken")),
	))
	// 扫码绑定：少一步复制 token（仅国服，国际服仍走 Token）
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📷 国服扫码登录", "qrLogin"),
	))
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)
	sendMessage := tgbotapi.NewMessage(chatId, "请选择要绑定的服务器")
	sendMessage.ReplyMarkup = inlineKeyboardMarkup
	config.Arknights.Send(sendMessage)
	return nil
}

// SetToken 设置token
func SetToken(update tgbotapi.Update) error {
	message := update.Message
	chatId := message.Chat.ID
	userId := message.From.ID
	token := message.Text

	_, _ = config.Arknights.SendChatAction(chatId, "typing")

	var userToken UserToken
	if err := json.Unmarshal([]byte(token), &userToken); err == nil {
		token = userToken.Data.Content
	}
	account, err := skland.Login(token, getServerName(chatId))
	if err != nil {
		config.Arknights.SendText(chatId, "登录失败！请检查token是否正确。")
		return err
	}
	return finishBind(chatId, userId, message.From.FullName(), token, account)
}

// CancelHandle 取消操作
func CancelHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	config.Arknights.ClearWaitMessage(userId)
	config.Arknights.SendText(chatId, "已取消操作")
	return nil
}
