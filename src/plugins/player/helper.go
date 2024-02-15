package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

// this function is a combination of getAccount and getPlayers Function
func getAccountAndPlayers(update tgbotapi.Update) (*account.UserAccount, *[]account.UserPlayer, error) {
	var userAccount *account.UserAccount
	var players *[]account.UserPlayer
	var err error
	userAccount, err = getAccount(update)
	if err == nil && userAccount != nil {
		players, err = getPlayers(update)
	}
	return userAccount, players, err
}
func NO_REQUIREMENT(_ tgbotapi.Update) bool { return true }

// this function will guide user to select the player and SEND CALLBACK TO callback_handle.go
// Param :
//   - chatId,userId,messageID from tg api update
//   - players : all bound account
//   - operation : all operation inside callback_handle.go
//
// Return :
//   - error : error if any
func playerSelector(update tgbotapi.Update, players []account.UserPlayer, operation PlayerOperation, perRequirement func(u2 tgbotapi.Update) bool) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	if !perRequirement(update) {
		return nil
	}
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, player := range players {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s(%s)", player.PlayerName, player.ServerName), fmt.Sprintf("%s,%d,%d,%s,%d", "player", operation, userId, player.Uid, messageId)),
		))
	}
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)
	sendMessage := tgbotapi.NewMessage(chatId, operation.getHintWordForPlayerSelection())
	sendMessage.ReplyMarkup = inlineKeyboardMarkup
	msg, _ := bot.Arknights.Send(sendMessage)
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
	return nil
}

// this function use to get account
// Param:
//   - update tg api for update
//
// Return
//   - the pointer to userAccounts nil if not found
//
// Note: this function will automatically ask user to bind their account
func getAccount(update tgbotapi.Update) (*account.UserAccount, error) {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	var userAccount account.UserAccount
	res := utils.GetAccountByUserId(userId).Scan(&userAccount)
	if res.RowsAffected == 0 {
		// 未绑定账号
		sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("未查询到绑定账号，请先进行[绑定](https://t.me/%s)。", viper.GetString("bot.name")))
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		sendMessage.ReplyToMessageID = messageId
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(chatId, messageId, 5)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil, nil
	}
	return &userAccount, nil
}

// this function use to get the players
// Param:
//   - update tg api for update
//
// Return
//   - the pointer to array of userPlayers nil if not found
//
// Note: this function will automatically ask user to bind their characters
func getPlayers(update tgbotapi.Update) (*[]account.UserPlayer, error) {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	var players []account.UserPlayer
	res := utils.GetPlayersByUserId(userId).Scan(&players)
	if res.RowsAffected == 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "您还未绑定任何角色！")
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(chatId, messageId, 5)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil, nil
	}
	return &players, nil
}
