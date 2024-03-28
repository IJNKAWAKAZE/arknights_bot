package account

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"arknights_bot/utils/telebot"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// BindHandle 绑定角色
func BindHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	sendMessage := tgbotapi.NewMessage(chatId, "请输入token或使用 /cancel 指令取消操作。")
	bot.Arknights.Send(sendMessage)
	sendMessage.Text = "如何获取token\n\n" +
		"1\\.前往 [森空岛](https://www.skland.com) 登录\n" +
		"2\\.打开网址复制content中的 token  [获取token](https://web-api.skland.com/account/info/hg)\n\n" +
		"手机用户且已登录森空岛直接点击此处获取token：[获取token](https://ss.xingzhige.com/skland.html)"
	sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
	bot.Arknights.Send(sendMessage)
	telebot.WaitMessage[chatId] = "setToken"
	return nil
}

// SetToken 设置token
func SetToken(update tgbotapi.Update) error {
	message := update.Message
	chatId := message.Chat.ID
	userId := message.From.ID
	token := message.Text

	sendAction := tgbotapi.NewChatAction(chatId, "typing")
	bot.Arknights.Send(sendAction)

	var userToken UserToken
	err := json.Unmarshal([]byte(token), &userToken)
	if err == nil {
		token = userToken.Data.Content
	}
	account, err := skland.Login(token)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, "登录失败！请检查token是否正确。")
		bot.Arknights.Send(sendMessage)
		return err
	}
	// 查询账户是否存在
	var userAccount UserAccount
	res := utils.GetAccountByUserIdAndSklandId(userId, account.UserId).Scan(&userAccount)
	if res.RowsAffected > 0 {
		// 更新账户信息
		userAccount.HypergryphToken = token
		userAccount.SklandToken = account.Skland.Token
		userAccount.SklandToken = account.Skland.Cred
		bot.DBEngine.Table("user_account").Save(&userAccount)
	} else {
		// 不存在 新增账户
		id, _ := gonanoid.New(32)
		userAccount = UserAccount{
			Id:              id,
			UserName:        utils.GetFullName(message.From),
			UserNumber:      userId,
			HypergryphToken: token,
			SklandToken:     account.Skland.Token,
			SklandCred:      account.Skland.Cred,
			SklandId:        account.UserId,
		}
		bot.DBEngine.Table("user_account").Create(&userAccount)
	}
	delete(telebot.WaitMessage, chatId)
	// 获取角色列表
	players, err := skland.ArknightsPlayers(account.Skland)
	if err != nil || len(players) == 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "未查询到绑定角色！")
		bot.Arknights.Send(sendMessage)
		return err
	}

	sklandIdMap[chatId] = account.UserId
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, player := range players {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s(%s)", player.NickName, player.ChannelName), fmt.Sprintf("%s,%s,%s,%s", "bind", player.Uid, player.ChannelName, player.NickName)),
		))
	}
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)
	sendMessage := tgbotapi.NewMessage(chatId, "请选择要绑定的角色")
	sendMessage.ReplyMarkup = inlineKeyboardMarkup
	bot.Arknights.Send(sendMessage)
	return nil
}

// CancelHandle 取消操作
func CancelHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	delete(telebot.WaitMessage, chatId)
	sendMessage := tgbotapi.NewMessage(chatId, "已取消操作")
	bot.Arknights.Send(sendMessage)
	return nil
}
