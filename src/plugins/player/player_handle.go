package player

import (
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/commandoperation"
	"arknights_bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var inited = false

// PlayerHandle 角色信息查询
func PlayerHandle(update tgbotapi.Update) (bool, error) {
	if !inited {
		initFactory()
		inited = true
	}
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	var userAccount account.UserAccount
	var players []account.UserPlayer
	var operationP *commandoperation.OperationI
	userAccountP, playersP, err := getAccountAndPlayers(update)
	if err != nil || userAccountP == nil || playersP == nil {
		return true, err
	}
	command := update.Message.Command()
	if commandoperation.HaveNextStep(chatId) {
		return true, commandoperation.GetStep(chatId).Run(update)
	}
	if len(command) != 0 {
		operationP = playerOperationFactory(command)
	}
	if operationP == nil {
		log.Printf("Unmatched Handle %s", update.Message.Command())
		return true, nil
	}
	operation := *operationP
	players = playersP
	userAccount = *userAccountP
	if players == nil || len(players) == 0 {
		log.Printf("Code reach impossible point players = %v after getPlayer warp", players)
	}
	if !operation.CheckRequirementsAndPrepare(update) {
		msg, isMarkDown := operation.HintOnRequirementsFailed()
		utils.SendMessage(chatId, msg, isMarkDown, &messageId)
		return true, nil
	}
	if len(players) > 1 {
		return true, playerSelector(update, userAccount, players, operation, command)
	}
	return operation.Run(players[0].Uid, userAccount, chatId, update.Message)
}
