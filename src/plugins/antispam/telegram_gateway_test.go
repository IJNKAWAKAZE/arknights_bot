package antispam

import (
	"errors"
	"testing"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

type fakeTelegramGateway struct {
	sends          []tgbotapi.Chattable
	requests       []tgbotapi.Chattable
	deletes        []deleteCall
	callbacks      []callbackCall
	callbackDelete []deleteCall
	bans           []banCall
	unbans         []banCall
	restricts      []restrictCall
	calls          []string
	queuedDeletes  []queueDeleteCall
	admins         map[int64]bool
	beforeSend     func()
	sendErr        error
	requestErr     error
	deleteErr      error
	banErr         error
	unbanErr       error
	restrictErr    error
	nextMessageID  int
}

type deleteCall struct {
	chatID    int64
	messageID int
}

type callbackCall struct {
	id        string
	showAlert bool
	text      string
}

type banCall struct {
	chatID int64
	userID int64
}

type restrictCall struct {
	chatID      int64
	userID      int64
	permissions string
}

type queueDeleteCall struct {
	chatID    int64
	messageID int
	delay     float64
}

func (fake *fakeTelegramGateway) Send(config tgbotapi.Chattable) (tgbotapi.Message, error) {
	if fake.beforeSend != nil {
		fake.beforeSend()
	}
	fake.sends = append(fake.sends, config)
	if fake.sendErr != nil {
		return tgbotapi.Message{}, fake.sendErr
	}
	fake.nextMessageID++
	return tgbotapi.Message{
		MessageID: fake.nextMessageID,
		Chat:      &tgbotapi.Chat{ID: -100100, Type: "supergroup"},
	}, nil
}

func (fake *fakeTelegramGateway) Request(config tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	fake.requests = append(fake.requests, config)
	if fake.requestErr != nil {
		return nil, fake.requestErr
	}
	return &tgbotapi.APIResponse{Ok: true}, nil
}

func (fake *fakeTelegramGateway) DeleteMessage(chatID int64, messageID int) (*tgbotapi.APIResponse, error) {
	fake.deletes = append(fake.deletes, deleteCall{chatID: chatID, messageID: messageID})
	if fake.deleteErr != nil {
		return nil, fake.deleteErr
	}
	return &tgbotapi.APIResponse{Ok: true}, nil
}

func (fake *fakeTelegramGateway) AnswerCallback(callbackID string, showAlert bool, text string) (*tgbotapi.APIResponse, error) {
	fake.callbacks = append(fake.callbacks, callbackCall{id: callbackID, showAlert: showAlert, text: text})
	return &tgbotapi.APIResponse{Ok: true}, nil
}

func (fake *fakeTelegramGateway) DeleteCallbackMessage(callback *tgbotapi.CallbackQuery) (*tgbotapi.APIResponse, error) {
	if callback != nil && callback.Message != nil && callback.Message.Chat != nil {
		fake.callbackDelete = append(fake.callbackDelete, deleteCall{chatID: callback.Message.Chat.ID, messageID: callback.Message.MessageID})
	}
	return &tgbotapi.APIResponse{Ok: true}, nil
}

func (fake *fakeTelegramGateway) BanChatMember(chatID, userID int64) (*tgbotapi.APIResponse, error) {
	fake.bans = append(fake.bans, banCall{chatID: chatID, userID: userID})
	fake.calls = append(fake.calls, "ban")
	if fake.banErr != nil {
		return nil, fake.banErr
	}
	return &tgbotapi.APIResponse{Ok: true}, nil
}

func (fake *fakeTelegramGateway) UnbanChatMember(chatID, userID int64) (*tgbotapi.APIResponse, error) {
	fake.unbans = append(fake.unbans, banCall{chatID: chatID, userID: userID})
	fake.calls = append(fake.calls, "unban")
	if fake.unbanErr != nil {
		return nil, fake.unbanErr
	}
	return &tgbotapi.APIResponse{Ok: true}, nil
}

func (fake *fakeTelegramGateway) RestrictChatMember(chatID, userID int64, permissions string) (*tgbotapi.APIResponse, error) {
	fake.restricts = append(fake.restricts, restrictCall{chatID: chatID, userID: userID, permissions: permissions})
	fake.calls = append(fake.calls, "restrict")
	if fake.restrictErr != nil {
		return nil, fake.restrictErr
	}
	return &tgbotapi.APIResponse{Ok: true}, nil
}

func (fake *fakeTelegramGateway) IsAdmin(chatID, userID int64) bool {
	return fake.admins[userID]
}

func (fake *fakeTelegramGateway) QueueDelete(chatID int64, messageID int, delay float64) {
	fake.queuedDeletes = append(fake.queuedDeletes, queueDeleteCall{chatID: chatID, messageID: messageID, delay: delay})
}

func useFakeTelegram(t *testing.T) *fakeTelegramGateway {
	t.Helper()
	fake := &fakeTelegramGateway{
		admins:        make(map[int64]bool),
		nextMessageID: 10000,
	}
	old := guestSpamTelegram
	guestSpamTelegram = fake
	t.Cleanup(func() {
		guestSpamTelegram = old
	})
	return fake
}

func errTelegram() error {
	return errors.New("telegram failed")
}
