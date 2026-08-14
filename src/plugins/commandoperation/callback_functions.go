package commandoperation

import (
	"arknights_bot/config"
	"arknights_bot/plugins/account"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"log"
)

var callBackMap = make(map[string]MultiuserCallBackFunction)
var nextStepMap = make(map[int64]NextStepOperation)

func AddNextStep(chatID int64, operation NextStepOperation, cmd string) bool {
	config.Arknights.SetWaitMessage(chatID, cmd)
	_, hasKey := nextStepMap[chatID]
	if hasKey {
		return false
	} else {
		nextStepMap[chatID] = operation
		return true
	}
}
func HaveNextStep(chatID int64) bool {
	_, hasKey := nextStepMap[chatID]
	hasMainKey := config.Arknights.HasWaitMessage(chatID)
	if hasKey && !hasMainKey {
		delete(nextStepMap, chatID)
		hasKey = false
	}
	return hasKey
}
func GetStep(chatID int64) *NextStepOperation {
	operation, hasKey := nextStepMap[chatID]
	if hasKey {
		return &operation
	} else {
		return nil
	}
}
func RemoveNextStep(chatID int64) {
	delete(nextStepMap, chatID)
}
func RemoveCallBack(callBackHash string) {
	delete(callBackMap, callBackHash)
}
func GetCallback(callBackHash string) (MultiuserCallBackFunction, bool) {
	callback, ok := callBackMap[callBackHash]
	if ok {
		RemoveCallBack(callBackHash)
	}
	return callback, ok
}
func AddCallback(callBackHash string, function MultiuserCallBackFunction) bool {
	_, hasKey := callBackMap[callBackHash]
	if hasKey {
		return false
	} else {
		callBackMap[callBackHash] = function
		return true
	}
}

type MultiuserCallBackFunction struct {
	Function func(uid string) error
	UserId   int64
}
type NextStepOperation struct {
	PlayerID      string
	Account       account.UserAccount
	Param         string
	NextOperation Operation
}

func (op NextStepOperation) Run(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	messageID := update.Message.MessageID
	var err error
	if op.NextOperation.CheckRequirementsAndPrepare(update) {

		err = op.NextOperation.Run(op.PlayerID, op.Account, chatId, update.Message)
		if err != nil {
			sent, sendErr := config.Arknights.ReplyText(chatId, messageID, "未知错误，请重试。")
			if sendErr != nil {
				log.Printf("%v can not be send error : %v", sent, sendErr)
			}
		}
		RemoveNextStep(chatId)
	} else {
		msg, isMarkDown := op.NextOperation.HintOnRequirementsFailed()
		tgMessage := tgbotapi.NewMessage(chatId, msg+" 使用 /cancel 指令取消操作")
		tgMessage.ReplyToMessageID = messageID
		if isMarkDown {
			tgMessage.ParseMode = tgbotapi.ModeMarkdownV2
		}
		sent, sendErr := config.Arknights.Send(tgMessage)
		if sendErr != nil {
			log.Printf("%v can not be send error : %v", sent, sendErr)
		}
	}
	return err
}

func NewMultiuserCallBackFunction(Function func(uid string) error, UserId int64) MultiuserCallBackFunction {
	return MultiuserCallBackFunction{
		Function: Function,
		UserId:   UserId,
	}
}
