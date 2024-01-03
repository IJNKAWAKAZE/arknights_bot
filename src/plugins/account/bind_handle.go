package account

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"arknights_bot/utils/telebot"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func BindHandle(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	// TODO 发送获取token教程
	telebot.WaitMessage[chatId] = "setToken"
	return true, nil
}

func SetToken(update tgbotapi.Update) (bool, error) {
	message := update.Message
	chatId := message.Chat.ID
	userId := message.From.ID
	token := message.Text
	account, err := skland.Login(token)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, "登录失败！请检查token是否正确。")
		bot.Arknights.Send(sendMessage)
		return true, err
	}
	var userAccount UserAccount
	// 查查询账户是否存在
	res := bot.DBEngine.Raw("select * from user_account where user_number = ?", userId).Scan(&userAccount)
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
		}
		bot.DBEngine.Table("user_account").Create(&userAccount)
	}
	delete(telebot.WaitMessage, chatId)
	// 获取角色列表
	players, err := skland.ArknihghtsPlayers(account.Skland)
	if err != nil || len(players) == 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "未查询到绑定角色！")
		bot.Arknights.Send(sendMessage)
		return true, err
	}

	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, player := range players {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s(%s)", player.NickName, player.ChannelName), fmt.Sprintf("%s,%s,%s", "bind", player.Uid, player.ChannelName)),
		))
	}
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)
	sendMessage := tgbotapi.NewMessage(chatId, "请选择要绑定的账号")
	sendMessage.ReplyMarkup = inlineKeyboardMarkup
	bot.Arknights.Send(sendMessage)
	return true, nil
}
