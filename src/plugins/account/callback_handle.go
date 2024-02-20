package account

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"strings"
)

// ChoosePlayer 选择绑定角色
func ChoosePlayer(callBack tgbotapi.Update) (bool, error) {
	callbackQuery := callBack.CallbackQuery
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 4 {
		return true, nil
	}

	userId := callbackQuery.From.ID
	chatId := callbackQuery.Message.Chat.ID

	uid := d[1]
	serverName := d[2]
	playerName := d[3]

	var userAccount UserAccount
	var userPlayer UserPlayer
	utils.GetAccountByUserId(userId).Scan(&userAccount)
	res := utils.GetPlayerByUserId(userId, uid).Scan(&userPlayer)
	if res.RowsAffected == 0 {
		id, _ := gonanoid.New(32)
		userPlayer = UserPlayer{
			Id:         id,
			AccountId:  userAccount.Id,
			UserName:   userAccount.UserName,
			UserNumber: userAccount.UserNumber,
			Uid:        uid,
			ServerName: serverName,
			PlayerName: playerName,
		}
		bot.DBEngine.Table("user_player").Create(&userPlayer)
	} else {
		userPlayer.PlayerName = playerName
		bot.DBEngine.Table("user_player").Save(&userPlayer)
		sendMessage := tgbotapi.NewMessage(chatId, "此角色已绑定，更新角色信息。")
		bot.Arknights.Send(sendMessage)
		return true, nil
	}
	sendMessage := tgbotapi.NewMessage(chatId, "角色绑定成功！")
	bot.Arknights.Send(sendMessage)
	return true, nil
}

// UnbindPlayer 解绑角色
func UnbindPlayer(callBack tgbotapi.Update) (bool, error) {
	callbackQuery := callBack.CallbackQuery
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 2 {
		return true, nil
	}

	userId := callbackQuery.From.ID
	chatId := callbackQuery.Message.Chat.ID
	messageId := callbackQuery.Message.MessageID

	uid := d[1]
	bot.DBEngine.Exec("delete from user_player where user_number = ? and uid = ?", userId, uid)
	sendMessage := tgbotapi.NewMessage(chatId, "角色解绑成功！")
	bot.Arknights.Send(sendMessage)
	delMsg := tgbotapi.NewDeleteMessage(chatId, messageId)
	bot.Arknights.Send(delMsg)
	return true, nil
}

// ChooseBTokenPlayer 选择设置BToken角色
func ChooseBTokenPlayer(callBack tgbotapi.Update) (bool, error) {
	callbackQuery := callBack.CallbackQuery
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 2 {
		return true, nil
	}

	userId := callbackQuery.From.ID
	chatId := callbackQuery.Message.Chat.ID
	messageId := callbackQuery.Message.MessageID
	uid := d[1]

	delMsg := tgbotapi.NewDeleteMessage(chatId, messageId)
	bot.Arknights.Send(delMsg)
	WaitBToken(chatId, userId, uid)
	return true, nil
}
