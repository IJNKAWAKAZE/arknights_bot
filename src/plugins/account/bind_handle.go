package account

import (
	"arknights_bot/config"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils/repo"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// BindHandle 绑定角色
func BindHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	var buttons [][]tgbotapi.InlineKeyboardButton
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("国服", fmt.Sprintf("%s,%s,%s", "chooseServer", "国服", "setToken")),
		tgbotapi.NewInlineKeyboardButtonData("国际服", fmt.Sprintf("%s,%s,%s", "chooseServer", "国际服", "setToken")),
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

	sendAction := tgbotapi.NewChatAction(chatId, "typing")
	config.Arknights.Send(sendAction)

	var userToken UserToken
	err := json.Unmarshal([]byte(token), &userToken)
	if err == nil {
		token = userToken.Data.Content
	}
	account, err := skland.Login(token, serverNameMap[chatId])
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, "登录失败！请检查token是否正确。")
		config.Arknights.Send(sendMessage)
		return err
	}
	// 查询账户是否存在
	var userAccount UserAccount
	res := repo.GetAccountByUserIdAndSklandId(userId, account.UserId).Scan(&userAccount)
	if res.RowsAffected > 0 {
		// 更新账户信息
		userAccount.HypergryphToken = token
		userAccount.SklandToken = account.Skland.Token
		userAccount.SklandCred = account.Skland.Cred
		config.DBEngine.Table("user_account").Save(&userAccount)
	} else {
		// 不存在 新增账户
		id, _ := gonanoid.New(32)
		userAccount = UserAccount{
			Id:              id,
			UserName:        message.From.FullName(),
			UserNumber:      userId,
			HypergryphToken: token,
			SklandToken:     account.Skland.Token,
			SklandCred:      account.Skland.Cred,
			SklandId:        account.UserId,
			ServerName:      serverNameMap[chatId],
		}
		config.DBEngine.Table("user_account").Create(&userAccount)
	}
	delete(tgbotapi.WaitMessage, chatId)
	// 获取角色列表
	players, err := skland.ArknightsPlayers(account.Skland, userAccount.ServerName)
	if err != nil || len(players) == 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "未查询到绑定角色！")
		config.Arknights.Send(sendMessage)
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
	config.Arknights.Send(sendMessage)
	return nil
}

// CancelHandle 取消操作
func CancelHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	delete(tgbotapi.WaitMessage, chatId)
	sendMessage := tgbotapi.NewMessage(chatId, "已取消操作")
	config.Arknights.Send(sendMessage)
	return nil
}
