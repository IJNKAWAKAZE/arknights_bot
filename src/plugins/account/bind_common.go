package account

import (
	"arknights_bot/config"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils/repo"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// finishBind 登录成功后的收尾：落库 + 抛出可选角色
func finishBind(chatId, userId int64, userName, token string, account skland.Account) error {
	var userAccount UserAccount
	res := repo.GetAccountByUserIdAndSklandId(userId, account.UserId).Scan(&userAccount)
	if err := repo.CheckDB(res, chatId); err != nil {
		return err
	}
	if res.RowsAffected > 0 {
		// 更新账户信息
		userAccount.HypergryphToken = token
		userAccount.SklandToken = account.Skland.Token
		userAccount.SklandCred = account.Skland.Cred
		if err := repo.CheckDB(config.DBEngine.Table("user_account").Save(&userAccount), chatId); err != nil {
			return err
		}
	} else {
		// 不存在 新增账户
		id, _ := gonanoid.New(32)
		userAccount = UserAccount{
			Id:              id,
			UserName:        userName,
			UserNumber:      userId,
			HypergryphToken: token,
			SklandToken:     account.Skland.Token,
			SklandCred:      account.Skland.Cred,
			SklandId:        account.UserId,
			ServerName:      getServerName(chatId),
		}
		if err := repo.CheckDB(config.DBEngine.Table("user_account").Create(&userAccount), chatId); err != nil {
			return err
		}
	}
	config.Arknights.ClearWaitMessage(userId)

	// 获取角色列表
	players, err := skland.ArknightsPlayers(account.Skland, userAccount.ServerName)
	if err != nil || len(players) == 0 {
		config.Arknights.SendText(chatId, "未查询到绑定角色！")
		return err
	}

	setSklandId(chatId, account.UserId)
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, player := range players {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s(%s)", player.NickName, player.ChannelName), fmt.Sprintf("%s,%s,%s,%s", "bind", player.Uid, player.ChannelName, player.NickName)),
		))
	}
	sendMessage := tgbotapi.NewMessage(chatId, "请选择要绑定的角色")
	sendMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)
	config.Arknights.Send(sendMessage)
	return nil
}

// finishReset 重设 token 后的收尾：只刷新凭据，不重选角色
func finishReset(chatId, userId int64, token string, account skland.Account) error {
	var userAccount UserAccount
	res := repo.GetAccountByUserIdAndSklandId(userId, account.UserId).Scan(&userAccount)
	if err := repo.CheckDB(res, chatId); err != nil {
		return err
	}
	config.Arknights.ClearWaitMessage(userId)
	if res.RowsAffected == 0 {
		// 换了账号扫码/换 token 时会对不上，提示走绑定
		config.Arknights.SendText(chatId, "未查询到绑定账号，请先进行绑定。")
		return nil
	}

	userAccount.HypergryphToken = token
	userAccount.SklandToken = account.Skland.Token
	userAccount.SklandCred = account.Skland.Cred
	if err := repo.CheckDB(config.DBEngine.Table("user_account").Save(&userAccount), chatId); err != nil {
		return err
	}
	config.Arknights.SendText(chatId, "重设token成功！")
	return nil
}
