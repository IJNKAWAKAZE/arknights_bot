package commandoperation

import (
	"arknights_bot/plugins/account"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var OperationTypeMaps = make(map[string]OperationI)

// OperationI don't embed this class at most of the time embed OperationAbstract is more useful
type OperationI interface {

	//this requirement of the string it can change base on the role of the user
	CheckRequirementsAndPrepare(update tgbotapi.Update) bool
	// run this operation (This function must be implemented by subClass)
	Run(uid string, userAccount account.UserAccount, chatId int64, message *tgbotapi.Message) error
	// when the argument did not match the requirement ask why
	HintOnRequirementsFailed() (string, bool)
	HintWordForPlayerSelection() string
	GetCallBackFunctionOnMultiPlayer(update tgbotapi.Update, account account.UserAccount, chatId int64, getTypeName string) MultiuserCallBackFunction
	NextStepOperation(playerUID string, userAccount account.UserAccount, param string) *NextStepOperation
}
type OperationAbstract struct {
	OperationI
}

func mapToStatic() {

}

func (_ OperationAbstract) NextStepOperation(playerUID string, userAccount account.UserAccount, param string) *NextStepOperation {
	return nil
}
func (operation OperationAbstract) CheckRequirementsAndPrepare(_ tgbotapi.Update) bool {
	return true
}
func (operation OperationAbstract) HintOnRequirementsFailed() (string, bool) {
	return "", false
}
func (operation OperationAbstract) GetCallBackFunctionOnMultiPlayer(update tgbotapi.Update, account account.UserAccount, chatId int64, getTypeName string) MultiuserCallBackFunction {
	newOP := OperationTypeMaps[getTypeName]
	return NewMultiuserCallBackFunction(
		func(uid string) error {
			return newOP.Run(uid, account, chatId, update.Message)
		},
		update.Message.From.ID)
}
func (_ OperationAbstract) HintWordForPlayerSelection() string {
	return "请选择要查询的角色"
}

type MultiStepOperation struct {
	OperationAbstract
}

func (_ MultiStepOperation) NextStepOperation(playerUID string, userAccount account.UserAccount, param string) *NextStepOperation {
	log.Println("Abstract method used ")
	return nil
}

// since we use chat id to flage next step so it need to be private method
func (_ MultiStepOperation) CheckRequirementsAndPrepare(update tgbotapi.Update) bool {
	return update.Message.Chat.IsPrivate()
}
