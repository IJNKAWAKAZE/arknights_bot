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
