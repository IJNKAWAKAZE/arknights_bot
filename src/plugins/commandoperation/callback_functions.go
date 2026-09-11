package commandoperation

import (
	"arknights_bot/config"
	"arknights_bot/plugins/account"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"log"
	"sync"
)

// 指令现在由 async 调度并发执行，这两个 map 会被多个 goroutine 访问，必须加锁。
var (
	mu          sync.Mutex
	callBackMap = make(map[string]MultiuserCallBackFunction)
	nextStepMap = make(map[int64]NextStepOperation)
)

func AddNextStep(chatID int64, operation NextStepOperation, cmd string) bool {
	mu.Lock()
	_, hasKey := nextStepMap[chatID]
	if !hasKey {
		nextStepMap[chatID] = operation
	}
	mu.Unlock()
	// 必须先登记等待状态再让调用方发送提示：
	// 否则用户看到提示立刻回复时，回复可能在状态写入前就被路由，导致下一步丢失。
	config.Arknights.SetWaitMessage(chatID, cmd)
	return !hasKey
}
func HaveNextStep(chatID int64) bool {
	hasMainKey := config.Arknights.HasWaitMessage(chatID)
	mu.Lock()
	defer mu.Unlock()
	_, hasKey := nextStepMap[chatID]
	if hasKey && !hasMainKey {
		delete(nextStepMap, chatID)
		hasKey = false
	}
	return hasKey
}
func GetStep(chatID int64) *NextStepOperation {
	mu.Lock()
	defer mu.Unlock()
	operation, hasKey := nextStepMap[chatID]
	if hasKey {
		return &operation
	} else {
		return nil
	}
}
func RemoveNextStep(chatID int64) {
	mu.Lock()
	defer mu.Unlock()
	delete(nextStepMap, chatID)
}
func RemoveCallBack(callBackHash string) {
	mu.Lock()
	defer mu.Unlock()
	delete(callBackMap, callBackHash)
}
func GetCallback(callBackHash string) (MultiuserCallBackFunction, bool) {
	mu.Lock()
	defer mu.Unlock()
	callback, ok := callBackMap[callBackHash]
	if ok {
		delete(callBackMap, callBackHash)
	}
	return callback, ok
}
func AddCallback(callBackHash string, function MultiuserCallBackFunction) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, hasKey := callBackMap[callBackHash]; hasKey {
		return false
	}
	callBackMap[callBackHash] = function
	return true
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
