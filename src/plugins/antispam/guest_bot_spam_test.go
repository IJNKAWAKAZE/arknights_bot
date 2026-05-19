package antispam

import (
	"testing"
	"time"

	bot "arknights_bot/config"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func TestCheckGuestBotSpam(t *testing.T) {
	bot.GuestBotSpamEnabled = true
	update := tgbotapi.Update{
		GuestMessage: &tgbotapi.Message{
			MessageID: 1,
			Chat: &tgbotapi.Chat{
				ID:   -1001,
				Type: "supergroup",
			},
			From: &tgbotapi.User{
				ID:        2001,
				IsBot:     true,
				FirstName: "Ad Bot",
			},
			GuestBotCallerUser: &tgbotapi.User{
				ID:        1001,
				FirstName: "Caller",
			},
		},
	}

	if !CheckGuestBotSpam(update) {
		t.Fatal("expected guest bot spam")
	}
}

func TestCheckGuestBotSpamIgnoresNormalMessage(t *testing.T) {
	bot.GuestBotSpamEnabled = true
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			MessageID: 1,
			Chat: &tgbotapi.Chat{
				ID:   -1001,
				Type: "supergroup",
			},
			From: &tgbotapi.User{
				ID:        1001,
				FirstName: "User",
			},
			Text: "hello",
		},
	}

	if CheckGuestBotSpam(update) {
		t.Fatal("normal message must not be treated as guest bot spam")
	}
}

func TestCheckGuestBotSpamDisabled(t *testing.T) {
	old := bot.GuestBotSpamEnabled
	bot.GuestBotSpamEnabled = false
	t.Cleanup(func() {
		bot.GuestBotSpamEnabled = old
	})

	update := tgbotapi.Update{GuestMessage: guestMessage(2001, 1001)}
	if CheckGuestBotSpam(update) {
		t.Fatal("disabled guest bot spam should not match")
	}
}

func TestIsGuestBotMessageShapes(t *testing.T) {
	if !isGuestBotMessage(&tgbotapi.Message{GuestBotCallerUser: &tgbotapi.User{ID: 1}}) {
		t.Fatal("caller user should mark guest bot message")
	}
	if !isGuestBotMessage(&tgbotapi.Message{GuestBotCallerChat: &tgbotapi.Chat{ID: -1}}) {
		t.Fatal("caller chat should mark guest bot message")
	}
	if !isGuestBotMessage(&tgbotapi.Message{GuestQueryID: "guest-query"}) {
		t.Fatal("guest query id should mark guest bot message")
	}
	if isGuestBotMessage(&tgbotapi.Message{}) {
		t.Fatal("plain message should not mark guest bot message")
	}
}

func TestIsTrackableMessage(t *testing.T) {
	if !isTrackableMessage(&tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: -1001, Type: "supergroup"},
		From: &tgbotapi.User{ID: 1001, FirstName: "User"},
		Text: "normal chat",
	}) {
		t.Fatal("normal group text should be trackable")
	}
	if isTrackableMessage(&tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: -1001, Type: "supergroup"},
		From: &tgbotapi.User{ID: 2001, IsBot: true, FirstName: "Bot"},
		Text: "bot chat",
	}) {
		t.Fatal("bot message should not be trackable")
	}
	if isTrackableMessage(&tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 1001, Type: "private"},
		From: &tgbotapi.User{ID: 1001, FirstName: "User"},
		Text: "private chat",
	}) {
		t.Fatal("private message should not be trackable")
	}
	if isTrackableMessage(&tgbotapi.Message{
		Chat:     &tgbotapi.Chat{ID: -1001, Type: "supergroup"},
		From:     &tgbotapi.User{ID: 1001, FirstName: "User"},
		Text:     "/guest_spam",
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 11}},
	}) {
		t.Fatal("command message should not be trackable")
	}
}

func TestEvaluateGuestMessageGuardAndCallerShapes(t *testing.T) {
	useFakeTelegram(t)

	if decision := EvaluateGuestMessage(nil); decision.DeleteMessage || decision.Reason != "" {
		t.Fatalf("nil decision = %+v, want empty", decision)
	}
	if decision := EvaluateGuestMessage(&tgbotapi.Message{From: &tgbotapi.User{ID: 1}}); decision.DeleteMessage || decision.Reason != "" {
		t.Fatalf("missing chat decision = %+v, want empty", decision)
	}
	if decision := EvaluateGuestMessage(&tgbotapi.Message{Chat: &tgbotapi.Chat{ID: -1}}); decision.DeleteMessage || decision.Reason != "" {
		t.Fatalf("missing from decision = %+v, want empty", decision)
	}
	if decision := EvaluateGuestMessage(guestQueryOnlyMessage(2001)); decision.DeleteMessage || decision.Reason != ReasonTrusted || len(decision.Logs) != 1 {
		t.Fatalf("guest query only decision = %+v, want allow log", decision)
	}
	if decision := EvaluateGuestMessage(guestChatMessage(2002, -2002)); decision.DeleteMessage || decision.Reason != ReasonTrusted || len(decision.Logs) != 1 {
		t.Fatalf("caller chat unknown decision = %+v, want allow log", decision)
	}
}

func TestCommandAndCallbackGuards(t *testing.T) {
	useFakeTelegram(t)

	if err := GuestSpamHandle(tgbotapi.Update{}); err != nil {
		t.Fatalf("nil guest spam command: %v", err)
	}
	if err := GuestSpamHandle(tgbotapi.Update{Message: &tgbotapi.Message{From: &tgbotapi.User{ID: 1}}}); err != nil {
		t.Fatalf("missing chat guest spam command: %v", err)
	}
	if err := GuestSpamHandle(tgbotapi.Update{Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: -1}}}); err != nil {
		t.Fatalf("missing from guest spam command: %v", err)
	}
	if err := GuestSpamLogHandle(tgbotapi.Update{}); err != nil {
		t.Fatalf("nil log command: %v", err)
	}
	if err := SelectRecentGuestCallback(tgbotapi.Update{}); err != nil {
		t.Fatalf("nil select callback: %v", err)
	}
	if err := SelectRecentGuestCallback(tgbotapi.Update{CallbackQuery: &tgbotapi.CallbackQuery{Data: "guestspam_select,1"}}); err != nil {
		t.Fatalf("select callback without message: %v", err)
	}
	if err := SelectRecentGuestCallback(tgbotapi.Update{CallbackQuery: &tgbotapi.CallbackQuery{Data: "bad", Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: -1}}}}); err != nil {
		t.Fatalf("bad select callback data: %v", err)
	}
	if err := SpamVoteCallback(tgbotapi.Update{}); err != nil {
		t.Fatalf("nil vote callback: %v", err)
	}
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: &tgbotapi.CallbackQuery{Data: "bad"}}); err != nil {
		t.Fatalf("bad vote callback data: %v", err)
	}
	if started, err := startSpamVote(nil, RecentGuestMessage{}); err != nil || started {
		t.Fatalf("nil start vote = (%v,%v), want (false,nil)", started, err)
	}
}

func TestTrustedRisk(t *testing.T) {
	now := time.Now()
	risk := MemberRisk{
		FirstSeenAt:        now.Add(-4 * 24 * time.Hour),
		RecentMessageCount: trustMessageCount,
	}
	if !isTrustedRiskAt(risk, now) {
		t.Fatal("old enough user with enough normal messages should be trusted")
	}
	risk.FirstSeenAt = now.Add(-2 * 24 * time.Hour)
	if isTrustedRiskAt(risk, now) {
		t.Fatal("fresh user should stay low trust")
	}
	risk.FirstSeenAt = now.Add(-4 * 24 * time.Hour)
	risk.RecentMessageCount = trustMessageCount - 1
	if isTrustedRiskAt(risk, now) {
		t.Fatal("user with insufficient recent messages should stay low trust")
	}
}

func guestMessage(guestBotID, callerID int64) *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: int(guestBotID % 100000),
		Chat: &tgbotapi.Chat{
			ID:    -100100,
			Type:  "supergroup",
			Title: "Guest Spam Test",
		},
		From: &tgbotapi.User{
			ID:        guestBotID,
			IsBot:     true,
			FirstName: "Guest",
			UserName:  "guest_bot",
		},
		GuestBotCallerUser: &tgbotapi.User{
			ID:        callerID,
			FirstName: "Caller",
			UserName:  "caller",
		},
	}
}

func botCallerMessage(guestBotID, callerID int64) *tgbotapi.Message {
	message := guestMessage(guestBotID, callerID)
	message.GuestBotCallerUser.IsBot = true
	return message
}

func guestChatMessage(guestBotID, callerChatID int64) *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: int(guestBotID % 100000),
		Chat: &tgbotapi.Chat{
			ID:    -100100,
			Type:  "supergroup",
			Title: "Guest Spam Test",
		},
		From: &tgbotapi.User{
			ID:        guestBotID,
			IsBot:     true,
			FirstName: "Guest",
			UserName:  "guest_bot",
		},
		GuestBotCallerChat: &tgbotapi.Chat{
			ID:    callerChatID,
			Type:  "channel",
			Title: "Caller Channel",
		},
	}
}

func guestQueryOnlyMessage(guestBotID int64) *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: int(guestBotID % 100000),
		Chat: &tgbotapi.Chat{
			ID:    -100100,
			Type:  "supergroup",
			Title: "Guest Spam Test",
		},
		From: &tgbotapi.User{
			ID:        guestBotID,
			IsBot:     true,
			FirstName: "Guest",
			UserName:  "guest_bot",
		},
		GuestQueryID: "guest-query",
	}
}

func hasLogAction(logs []SpamLog, action string) bool {
	for _, item := range logs {
		if item.Action == action {
			return true
		}
	}
	return false
}

func TestWarningMuteDuration(t *testing.T) {
	if got := MuteDuration(1); got != 24*time.Hour {
		t.Fatalf("level 1 mute = %s, want 24h", got)
	}
	if got := MuteDuration(2); got != 7*24*time.Hour {
		t.Fatalf("level 2 mute = %s, want 7d", got)
	}
	if got := MuteDuration(9); got != 30*24*time.Hour {
		t.Fatalf("high level mute = %s, want 30d", got)
	}
}

func TestRequiredVoteCount(t *testing.T) {
	if _, ok := requiredVoteCount(2); ok {
		t.Fatal("active user count below 3 should invalidate vote")
	}
	if got, ok := requiredVoteCount(3); !ok || got != 2 {
		t.Fatalf("active=3 required=(%d,%v), want (2,true)", got, ok)
	}
	if got, ok := requiredVoteCount(4); !ok || got != 2 {
		t.Fatalf("active=4 required=(%d,%v), want (2,true)", got, ok)
	}
	if got, ok := requiredVoteCount(5); !ok || got != 3 {
		t.Fatalf("active=5 required=(%d,%v), want (3,true)", got, ok)
	}
}
