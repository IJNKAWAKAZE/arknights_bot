package account

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"arknights_bot/utils/telebot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SetTokenHandle 重设token

func SetTokenHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID

	var userAccount UserAccount

	res := utils.GetAccountByUserId(userId).Scan(&userAccount)
	if res.RowsAffected == 0 {
		// 未绑定账号
		sendMessage := tgbotapi.NewMessage(chatId, "未查询到绑定账号，请先进行绑定。")
		bot.Arknights.Send(sendMessage)
		return nil
	}
	sendMessage := tgbotapi.NewMessage(chatId, "请输入新token或使用 /cancel 指令取消操作。")
	bot.Arknights.Send(sendMessage)
	sendMessage.Text = "如何获取token\n\n" +
		"1\\.前往 [森空岛](https://www.skland.com) 登录\n" +
		"2\\.打开网址复制content中的 token  [获取token](https://web-api.skland.com/account/info/hg)\n\n" +
		"手机用户且已登录森空岛直接点击此处获取token：[获取token](https://ss.xingzhige.com/skland.html)"
	sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
	bot.Arknights.Send(sendMessage)
	telebot.WaitMessage[chatId] = "resetToken"
	return nil
}

// ResetToken 重设token
func ResetToken(update tgbotapi.Update) error {
	message := update.Message
	chatId := message.Chat.ID
	userId := message.From.ID
	token := message.Text

	sendAction := tgbotapi.NewChatAction(chatId, "typing")
	bot.Arknights.Send(sendAction)

	account, err := skland.Login(token)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, "登录失败！请检查token是否正确。")
		bot.Arknights.Send(sendMessage)
		return err
	}
	// 查查询账户信息
	var userAccount UserAccount
	res := utils.GetAccountByUserId(userId).Scan(&userAccount)
	if res.RowsAffected > 0 {
		// 更新账户信息
		userAccount.HypergryphToken = token
		userAccount.SklandToken = account.Skland.Token
		userAccount.SklandToken = account.Skland.Cred
		bot.DBEngine.Table("user_account").Save(&userAccount)
		sendMessage := tgbotapi.NewMessage(chatId, "重设token成功！")
		bot.Arknights.Send(sendMessage)
	}
	delete(telebot.WaitMessage, chatId)
	return nil
}
