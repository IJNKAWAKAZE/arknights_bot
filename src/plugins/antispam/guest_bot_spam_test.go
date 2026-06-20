package antispam

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	bot "arknights_bot/config"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
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

func TestActiveUsersExpireIndividually(t *testing.T) {
	setupGuestSpamRedisOnlyForUnitTest(t)

	now := time.Now()
	writeActiveUserScore(-100100, 1001, now.Add(-2*time.Minute))
	writeActiveUserScore(-100100, 1002, now.Add(-12*time.Minute))

	if got := ActiveUserCount(-100100); got != 1 {
		t.Fatalf("active users = %d, want 1 after pruning expired members", got)
	}
	if !IsActiveUser(-100100, 1001) {
		t.Fatal("fresh user should stay active")
	}
	if IsActiveUser(-100100, 1002) {
		t.Fatal("expired user should be inactive")
	}
}

func TestActiveUserAtExpiryCutoffRemainsActive(t *testing.T) {
	setupGuestSpamRedisOnlyForUnitTest(t)

	now := time.Now().Truncate(time.Second)
	writeActiveUserScore(-100100, 1001, now.Add(-activeWindowTTL))
	pruneActiveUsers(-100100, now)

	if got := ActiveUserCount(-100100); got != 1 {
		t.Fatalf("active users at cutoff = %d, want 1", got)
	}
	if !IsActiveUser(-100100, 1001) {
		t.Fatal("user at exact cutoff should stay active")
	}
}

func TestTrackActivityHandleRefreshesSortedSetScore(t *testing.T) {
	setupGuestSpamRedisOnlyForUnitTest(t)

	staleAt := time.Now().Add(-activeWindowTTL - 2*time.Minute).Truncate(time.Second)
	writeActiveUserScore(testIntegrationChatID, 880001, staleAt)
	if IsActiveUser(testIntegrationChatID, 880001) {
		t.Fatal("stale active user should be inactive before refresh")
	}

	if err := TrackActivityHandle(tgbotapi.Update{Message: trackableGuestSpamMessage(880001, "hello")}); err != nil {
		t.Fatalf("track activity: %v", err)
	}
	if got := ActiveUserCount(testIntegrationChatID); got != 1 {
		t.Fatalf("active users = %d, want 1", got)
	}
	if !IsActiveUser(testIntegrationChatID, 880001) {
		t.Fatal("tracked user should be active")
	}
	score, err := bot.GoRedis.ZScore(redisCtx, activeUsersKey(testIntegrationChatID), "880001").Result()
	if err != nil {
		t.Fatalf("read refreshed score: %v", err)
	}
	if score <= float64(staleAt.Unix()) {
		t.Fatalf("refreshed score = %f, want > %d", score, staleAt.Unix())
	}
}

func TestRecordMessageActivityDoesNotSetActiveUsersTTL(t *testing.T) {
	setupGuestSpamRedisOnlyForUnitTest(t)

	RecordMessageActivity(testIntegrationChatID, 880001, "User")

	ttl, err := bot.GoRedis.TTL(redisCtx, activeUsersKey(testIntegrationChatID)).Result()
	if err != nil {
		t.Fatalf("read active users ttl: %v", err)
	}
	if ttl != -1 {
		t.Fatalf("active users ttl = %s, want no expiration", ttl)
	}
}

func writeActiveUserScore(chatID int64, userID int64, when time.Time) {
	if bot.GoRedis == nil {
		return
	}
	if err := bot.GoRedis.ZAdd(redisCtx, activeUsersKey(chatID), &redis.Z{
		Score:  float64(when.Unix()),
		Member: strconv.FormatInt(userID, 10),
	}).Err(); err != nil {
		panic(err)
	}
}

func writeActiveUsers(chatID int64, when time.Time, userIDs ...int64) {
	for _, userID := range userIDs {
		writeActiveUserScore(chatID, userID, when)
	}
}

const testIntegrationChatID = int64(-100100)

func setupGuestSpamRedisOnlyForUnitTest(t *testing.T) {
	t.Helper()
	mini := miniredis.RunT(t)

	bot.GoRedis = redis.NewClient(&redis.Options{
		Addr: mini.Addr(),
		DB:   0,
	})
	if err := bot.GoRedis.Ping(redisCtx).Err(); err != nil {
		t.Fatalf("ping redis: %v", err)
	}
	clearActiveUserTestRedis(t)
	t.Cleanup(func() {
		clearActiveUserTestRedis(t)
		_ = bot.GoRedis.Close()
		bot.GoRedis = nil
		mini.Close()
	})
}

func clearActiveUserTestRedis(t *testing.T) {
	t.Helper()
	if bot.GoRedis == nil {
		return
	}
	iter := bot.GoRedis.Scan(redisCtx, 0, redisPrefix+":*", 0).Iterator()
	for iter.Next(redisCtx) {
		if err := bot.GoRedis.Del(redisCtx, iter.Val()).Err(); err != nil {
			t.Fatalf("delete redis key %s: %v", iter.Val(), err)
		}
	}
	if err := iter.Err(); err != nil {
		t.Fatalf("scan redis: %v", err)
	}
}

func trackableGuestSpamMessage(userID int64, text string) *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: int(userID % 100000),
		Chat: &tgbotapi.Chat{
			ID:    testIntegrationChatID,
			Type:  "supergroup",
			Title: "Guest Spam Test",
		},
		From: &tgbotapi.User{
			ID:        userID,
			FirstName: "User",
			UserName:  fmt.Sprintf("user_%d", userID),
		},
		Text: text,
	}
}
