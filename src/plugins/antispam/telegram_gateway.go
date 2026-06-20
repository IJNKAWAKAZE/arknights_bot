package antispam

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

type telegramGateway interface {
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
	Request(tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	DeleteMessage(chatID int64, messageID int) (*tgbotapi.APIResponse, error)
	AnswerCallback(callbackID string, showAlert bool, text string) (*tgbotapi.APIResponse, error)
	DeleteCallbackMessage(callback *tgbotapi.CallbackQuery) (*tgbotapi.APIResponse, error)
	BanChatMember(chatID, userID int64) (*tgbotapi.APIResponse, error)
	UnbanChatMember(chatID, userID int64) (*tgbotapi.APIResponse, error)
	RestrictChatMember(chatID, userID int64, permissions string) (*tgbotapi.APIResponse, error)
	IsAdmin(chatID, userID int64) bool
	QueueDelete(chatID int64, messageID int, delay float64)
}

var guestSpamTelegram telegramGateway = liveTelegramGateway{}

type liveTelegramGateway struct{}

func (liveTelegramGateway) Send(config tgbotapi.Chattable) (tgbotapi.Message, error) {
	return bot.Arknights.Send(config)
}

func (liveTelegramGateway) Request(config tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	return bot.Arknights.Request(config)
}

func (liveTelegramGateway) DeleteMessage(chatID int64, messageID int) (*tgbotapi.APIResponse, error) {
	return bot.Arknights.Request(tgbotapi.NewDeleteMessage(chatID, messageID))
}

func (liveTelegramGateway) AnswerCallback(callbackID string, showAlert bool, text string) (*tgbotapi.APIResponse, error) {
	config := tgbotapi.NewCallback(callbackID, text)
	if showAlert {
		config = tgbotapi.NewCallbackWithAlert(callbackID, text)
	}
	return bot.Arknights.Request(config)
}

func (liveTelegramGateway) DeleteCallbackMessage(callback *tgbotapi.CallbackQuery) (*tgbotapi.APIResponse, error) {
	if callback == nil || callback.Message == nil || callback.Message.Chat == nil {
		return nil, nil
	}
	return bot.Arknights.Request(tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID))
}

func (liveTelegramGateway) BanChatMember(chatID, userID int64) (*tgbotapi.APIResponse, error) {
	return bot.Arknights.BanChatMember(chatID, userID)
}

func (liveTelegramGateway) UnbanChatMember(chatID, userID int64) (*tgbotapi.APIResponse, error) {
	return bot.Arknights.UnbanChatMember(chatID, userID)
}

func (liveTelegramGateway) RestrictChatMember(chatID, userID int64, permissions string) (*tgbotapi.APIResponse, error) {
	return bot.Arknights.RestrictChatMember(chatID, userID, permissions)
}

func (liveTelegramGateway) IsAdmin(chatID, userID int64) bool {
	return bot.Arknights.IsAdmin(chatID, userID)
}

func (liveTelegramGateway) QueueDelete(chatID int64, messageID int, delay float64) {
	messagecleaner.AddDelQueue(chatID, messageID, delay)
}
