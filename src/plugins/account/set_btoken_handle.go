package account

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"arknights_bot/utils/telebot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SetBTokenHandle 设置btoken

func SetBTokenHandle(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID

	var userAccount UserAccount

	res := utils.GetAccountByUserId(userId).Scan(&userAccount)
	if res.RowsAffected == 0 {
		// 未绑定账号
		sendMessage := tgbotapi.NewMessage(chatId, "未查询到绑定账号，请先进行绑定。")
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		bot.Arknights.Send(sendMessage)
		return true, nil
	}
	sendMessage := tgbotapi.NewMessage(chatId, "请输入btoken或使用 /cancel 指令取消操作。")
	bot.Arknights.Send(sendMessage)
	sendMessage.Text = "如何获取token\n\n" +
		"1\\.前往 [官网](https://ak.hypergryph.com/user/bilibili/login) 登录\n" +
		"2\\.打开网址复制content中的 token  [获取token](https://web-api.hypergryph.com/account/info/ak-b)"
	sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
	bot.Arknights.Send(sendMessage)
	telebot.WaitMessage[chatId] = "bToken"
	return true, nil
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
	// 查查询账户信息
	var userAccount UserAccount
	res := utils.GetAccountByUserId(userId).Scan(&userAccount)
	if res.RowsAffected > 0 {
		// 更新账户信息
		userAccount.BToken = token
		bot.DBEngine.Table("user_account").Save(&userAccount)
		sendMessage := tgbotapi.NewMessage(chatId, "设置BToken成功！")
		bot.Arknights.Send(sendMessage)
	}
	delete(telebot.WaitMessage, chatId)
	return true, nil
}
