package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/commandoperation"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"log"
)

var hashSize = 11

// this function is a combination of getAccount and getPlayers Function
func getAccountAndPlayers(update tgbotapi.Update) (*account.UserAccount, []account.UserPlayer, error) {
	var userAccount *account.UserAccount
	var players []account.UserPlayer
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
func playerSelector(update tgbotapi.Update, userAccount account.UserAccount, players []account.UserPlayer, operation commandoperation.OperationI, nameType string) error {
	chatId := update.Message.Chat.ID
	callBackFunction := operation.GetCallBackFunctionOnMultiPlayer(update, userAccount, chatId, nameType)
	functionHash := utils.RandStringBytesMaskImprSrcUnsafe(hashSize)
	//keep trying to make sure key not duplicate
	for ; !commandoperation.AddCallback(functionHash, callBackFunction); functionHash = utils.RandStringBytesMaskImprSrcUnsafe(hashSize) {
	}

	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, player := range players {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s(%s)", player.PlayerName, player.ServerName), fmt.Sprintf("%s,%s,%s", "player", functionHash, player.Uid)),
		))
	}
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)
	sendMessage := tgbotapi.NewMessage(chatId, operation.HintWordForPlayerSelection())
	sendMessage.ReplyMarkup = inlineKeyboardMarkup
	msg, err := bot.Arknights.Send(sendMessage)
	if err != nil {
		log.Println("can not send massage ", err)
	}
	messagecleaner.AddDelQueueFuncHash(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay, functionHash)
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
	if res.RowsAffected != 0 {
		if res.RowsAffected == -1 {
			panic("SQL ERROR check your sql config")
		}
		return &userAccount, nil
	} else {
		// 未绑定账号
		sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("未查询到绑定账号，请先进行[绑定](https://t.me/%s)。", viper.GetString("bot.name")))
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			log.Println("can not send massage ", err)
		}
		messagecleaner.AddDelQueue(chatId, messageId, 5)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil, nil
	}

}

// this function use to get the players
// Param:
//   - update tg api for update
//
// Return
//   - the pointer to array of userPlayers nil if not found
//
// Note: this function will automatically ask user to bind their characters
func getPlayers(update tgbotapi.Update) ([]account.UserPlayer, error) {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	var players []account.UserPlayer
	res := utils.GetPlayersByUserId(userId).Scan(&players)
	if res.RowsAffected != 0 {
		if res.RowsAffected == -1 {
			panic("SQL ERROR check your sql config")
		}
		return players, nil
	} else {
		sendMessage := tgbotapi.NewMessage(chatId, "您还未绑定任何角色！")
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			log.Println("can not send massage ", err)
		}
		messagecleaner.AddDelQueue(chatId, messageId, 5)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil, nil
	}

}
