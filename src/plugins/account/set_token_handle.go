package account

import (
	"arknights_bot/config"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils/repo"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

// SetTokenHandle 重设token
func SetTokenHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID

	var userAccount UserAccount

	res := repo.GetAccountByUserId(userId).Scan(&userAccount)
	if res.RowsAffected == 0 {
		// 未绑定账号
		config.Arknights.SendText(chatId, "未查询到绑定账号，请先进行绑定。")
		return nil
	}
	var buttons [][]tgbotapi.InlineKeyboardButton
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("国服", fmt.Sprintf("%s,%s,%s", "chooseServer", "国服", "resetToken")),
		tgbotapi.NewInlineKeyboardButtonData("国际服", fmt.Sprintf("%s,%s, %s", "chooseServer", "国际服", "resetToken")),
	))
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)
	sendMessage := tgbotapi.NewMessage(chatId, "请选择要绑定的服务器")
	sendMessage.ReplyMarkup = inlineKeyboardMarkup
	config.Arknights.Send(sendMessage)
	return nil
}

// ResetToken 重设token
func ResetToken(update tgbotapi.Update) error {
	message := update.Message
	chatId := message.Chat.ID
	userId := message.From.ID
	token := message.Text

	_, _ = config.Arknights.SendChatAction(chatId, "typing")

	var userToken UserToken
	err := json.Unmarshal([]byte(token), &userToken)
	if err == nil {
		token = userToken.Data.Content
	}
	account, err := skland.Login(token, serverNameMap[chatId])
	if err != nil {
		config.Arknights.SendText(chatId, "登录失败！请检查token是否正确。")
		return err
	}
	// 查查询账户信息
	var userAccount UserAccount
	res := repo.GetAccountByUserIdAndSklandId(userId, account.UserId).Scan(&userAccount)
	if res.RowsAffected > 0 {
		// 更新账户信息
		userAccount.HypergryphToken = token
		userAccount.SklandToken = account.Skland.Token
		userAccount.SklandCred = account.Skland.Cred
		config.DBEngine.Table("user_account").Save(&userAccount)
		config.Arknights.SendText(chatId, "重设token成功！")
	}
	config.Arknights.ClearWaitMessage(userId)
	return nil
}
