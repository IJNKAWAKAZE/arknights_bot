//go:build integration

package antispam

import (
	bot "arknights_bot/config"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	integrationChatID      = int64(-100100)
	integrationGuestBotID  = int64(990001)
	integrationGuestBot2ID = int64(990002)
	integrationCallerID    = int64(880001)
)

func TestGuestSpamIntegrationLowTrustSyncAndReload(t *testing.T) {
	setupGuestSpamIntegration(t)

	message := testGuestMessage(integrationGuestBotID, integrationCallerID)
	decision := EvaluateGuestMessage(message)
	ApplyGuestSpamState(decision)

	if !decision.DeleteMessage || !decision.BlacklistBot || !decision.WarnCaller || decision.MuteCaller {
		t.Fatalf("low trust decision = %+v, want delete+blacklist+warn without first mute", decision)
	}
	if !IsBlacklisted(integrationGuestBotID) {
		t.Fatal("guest bot should be blacklisted in redis immediately")
	}

	if err := SyncCacheToDB(); err != nil {
		t.Fatalf("sync cache to db: %v", err)
	}

	clearGuestSpamRedis(t)
	cached, err := bot.GoRedis.SIsMember(redisCtx, blacklistSetKey(), integrationGuestBotID).Result()
	if err != nil {
		t.Fatalf("read redis blacklist: %v", err)
	}
	if cached {
		t.Fatal("redis blacklist should be empty before reload")
	}
	if err := LoadCacheFromDB(); err != nil {
		t.Fatalf("load cache from db: %v", err)
	}
	if !IsBlacklisted(integrationGuestBotID) {
		t.Fatal("guest bot should be reloaded from db into redis")
	}

	reloaded := EvaluateGuestMessage(message)
	if !reloaded.DeleteMessage || !reloaded.BanCallerUser || reloaded.BlacklistBot {
		t.Fatalf("blacklist decision = %+v, want delete+ban caller without re-blacklist", reloaded)
	}
}

func TestGuestSpamIntegrationHandleLowTrustDeleteFailureStillBlacklists(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)
	fake.deleteErr = errTelegram()

	err := GuestBotSpamHandle(tgbotapi.Update{GuestMessage: testGuestMessage(integrationGuestBotID, integrationCallerID)})
	if err != nil {
		t.Fatalf("handle low trust: %v", err)
	}
	if len(fake.deletes) != 1 {
		t.Fatalf("delete calls = %d, want 1", len(fake.deletes))
	}
	if !IsBlacklisted(integrationGuestBotID) {
		t.Fatal("low trust guest bot should be blacklisted even if delete fails")
	}
	logs := RecentLogs(integrationChatID, 10)
	if !hasLogAction(logs, ActionDeleteFailed) || !hasLogAction(logs, ActionAutoBlacklist) || !hasLogAction(logs, ActionWarnCaller) {
		t.Fatalf("logs = %+v, want delete_failed, auto_blacklist, warn_caller", logs)
	}
}

func TestGuestSpamIntegrationTrustedCallerAllowsUnknownBot(t *testing.T) {
	setupGuestSpamIntegration(t)

	firstSeen := startOfDay(time.Now().Add(-4 * 24 * time.Hour))
	risk := MemberRisk{
		ID:                 riskID(integrationChatID, integrationCallerID),
		ChatID:             integrationChatID,
		UserID:             integrationCallerID,
		UserName:           "Trusted Caller",
		FirstSeenAt:        firstSeen,
		LastMessageAt:      time.Now(),
		RecentMessageCount: trustMessageCount,
	}
	setMemberRisk(risk, true)
	for i := int64(0); i < trustMessageCount; i++ {
		key := memberActivityDayKey(integrationChatID, integrationCallerID, dayKey(time.Now()))
		bot.GoRedis.Incr(redisCtx, key)
		bot.GoRedis.Expire(redisCtx, key, activityDayTTL)
	}

	decision := EvaluateGuestMessage(testGuestMessage(integrationGuestBot2ID, integrationCallerID))
	ApplyGuestSpamState(decision)
	if decision.DeleteMessage || decision.BlacklistBot || decision.WarnCaller || !decision.Trusted {
		t.Fatalf("trusted decision = %+v, want allow+trusted", decision)
	}
	if IsBlacklisted(integrationGuestBot2ID) {
		t.Fatal("trusted unknown guest bot should not be blacklisted")
	}
}

func TestGuestSpamIntegrationHandleTrustedCallerNoTelegramAction(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)
	seedTrustedCaller(integrationCallerID)

	err := GuestBotSpamHandle(tgbotapi.Update{GuestMessage: testGuestMessage(integrationGuestBot2ID, integrationCallerID)})
	if err != nil {
		t.Fatalf("handle trusted guest message: %v", err)
	}
	if len(fake.deletes) != 0 || len(fake.bans) != 0 || len(fake.requests) != 0 {
		t.Fatalf("trusted caller should not call telegram actions, deletes=%d bans=%d requests=%d", len(fake.deletes), len(fake.bans), len(fake.requests))
	}
	if IsBlacklisted(integrationGuestBot2ID) {
		t.Fatal("trusted caller should not blacklist unknown guest bot")
	}
}

func TestGuestSpamIntegrationWarningEscalatesToMute(t *testing.T) {
	setupGuestSpamIntegration(t)

	first := EvaluateGuestMessage(testGuestMessage(991001, integrationCallerID))
	if first.MuteCaller {
		t.Fatalf("first warning decision = %+v, should not mute", first)
	}
	second := EvaluateGuestMessage(testGuestMessage(991002, integrationCallerID))
	if !second.MuteCaller || second.MuteDuration != 24*time.Hour {
		t.Fatalf("second warning decision = %+v, want 24h mute", second)
	}
}

func TestGuestSpamIntegrationHandleMuteSuccessLogged(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	if err := GuestBotSpamHandle(tgbotapi.Update{GuestMessage: testGuestMessage(991201, integrationCallerID)}); err != nil {
		t.Fatalf("first warning: %v", err)
	}
	if err := GuestBotSpamHandle(tgbotapi.Update{GuestMessage: testGuestMessage(991202, integrationCallerID)}); err != nil {
		t.Fatalf("second warning: %v", err)
	}
	if len(fake.requests) != 1 {
		t.Fatalf("mute request calls = %d, want 1", len(fake.requests))
	}
	if !hasLogAction(RecentLogs(integrationChatID, 20), ActionMuteCaller) {
		t.Fatalf("logs = %+v, want mute success log", RecentLogs(integrationChatID, 20))
	}
}

func TestGuestSpamIntegrationHandleMuteFailureLogged(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)
	fake.requestErr = errors.New("mute failed")

	first := GuestBotSpamHandle(tgbotapi.Update{GuestMessage: testGuestMessage(991101, integrationCallerID)})
	second := GuestBotSpamHandle(tgbotapi.Update{GuestMessage: testGuestMessage(991102, integrationCallerID)})
	if first != nil || second != nil {
		t.Fatalf("handle warning sequence errors = %v %v", first, second)
	}
	if len(fake.requests) != 1 {
		t.Fatalf("mute request calls = %d, want 1", len(fake.requests))
	}
	if !hasLogAction(RecentLogs(integrationChatID, 20), ActionMuteCaller) {
		t.Fatalf("logs = %+v, want mute failure log", RecentLogs(integrationChatID, 20))
	}
}

func TestGuestSpamIntegrationBlacklistedCallerChat(t *testing.T) {
	setupGuestSpamIntegration(t)

	guestBotID := int64(992001)
	AddBlacklist(GuestBotBlacklist{
		ID:          blacklistID(guestBotID),
		BotID:       guestBotID,
		BotName:     "Known Spam Bot",
		BotUserName: "known_spam_bot",
		Source:      "test",
	}, true)

	decision := EvaluateGuestMessage(testGuestChatMessage(guestBotID, -100200))
	if !decision.DeleteMessage || !decision.BanCallerChat || decision.BanCallerUser || decision.BlacklistBot {
		t.Fatalf("caller chat decision = %+v, want delete+ban caller chat only", decision)
	}
}

func TestGuestSpamIntegrationHandleBlacklistedCallerActionsAndFailures(t *testing.T) {
	tests := []struct {
		name      string
		message   *tgbotapi.Message
		banErr    error
		reqErr    error
		wantBan   int
		wantReq   int
		wantLog   string
		wantNoBan bool
	}{
		{
			name:    "caller user banned",
			message: testGuestMessage(992101, integrationCallerID),
			wantBan: 1,
			wantLog: ActionBanCaller,
		},
		{
			name:    "caller user ban failure logged",
			message: testGuestMessage(992102, integrationCallerID),
			banErr:  errTelegram(),
			wantBan: 1,
			wantLog: ActionBanCaller,
		},
		{
			name:    "caller chat banned",
			message: testGuestChatMessage(992103, -100201),
			wantReq: 1,
			wantLog: ActionBanCallerChat,
		},
		{
			name:    "caller chat ban failure logged",
			message: testGuestChatMessage(992104, -100202),
			reqErr:  errTelegram(),
			wantReq: 1,
			wantLog: ActionBanCallerChat,
		},
		{
			name:      "bot caller not banned",
			message:   testBotCallerMessage(992105, integrationCallerID),
			wantNoBan: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupGuestSpamIntegration(t)
			fake := useFakeTelegram(t)
			fake.banErr = tt.banErr
			fake.requestErr = tt.reqErr
			AddBlacklist(GuestBotBlacklist{BotID: tt.message.From.ID, BotName: "Known Bot"}, true)

			if err := GuestBotSpamHandle(tgbotapi.Update{GuestMessage: tt.message}); err != nil {
				t.Fatalf("handle blacklist: %v", err)
			}
			if len(fake.bans) != tt.wantBan {
				t.Fatalf("ban calls = %d, want %d", len(fake.bans), tt.wantBan)
			}
			if len(fake.requests) != tt.wantReq {
				t.Fatalf("request calls = %d, want %d", len(fake.requests), tt.wantReq)
			}
			if tt.wantNoBan && (len(fake.bans) != 0 || len(fake.requests) != 0) {
				t.Fatalf("bot caller should not be punished, bans=%d requests=%d", len(fake.bans), len(fake.requests))
			}
			if tt.wantLog != "" && !hasLogAction(RecentLogs(integrationChatID, 10), tt.wantLog) {
				t.Fatalf("logs = %+v, want action %s", RecentLogs(integrationChatID, 10), tt.wantLog)
			}
		})
	}
}

func TestGuestSpamIntegrationActivityDateSyncAndReload(t *testing.T) {
	db := setupGuestSpamIntegration(t)

	firstSeen := startOfDay(time.Now().AddDate(0, 0, -4))
	setMemberRisk(MemberRisk{
		ID:          riskID(integrationChatID, integrationCallerID),
		ChatID:      integrationChatID,
		UserID:      integrationCallerID,
		UserName:    "Active Caller",
		FirstSeenAt: firstSeen,
	}, true)

	for i := int64(0); i < trustMessageCount; i++ {
		RecordMessageActivity(integrationChatID, integrationCallerID, "Active Caller")
	}
	if err := SyncCacheToDB(); err != nil {
		t.Fatalf("sync cache to db: %v", err)
	}

	today := dayKey(time.Now())
	var activity MemberActivity
	if err := db.Where("id = ?", activityID(integrationChatID, integrationCallerID, today)).First(&activity).Error; err != nil {
		t.Fatalf("load activity from db: %v", err)
	}
	if activity.ActivityDay.IsZero() || dayKey(activity.ActivityDay) != today {
		t.Fatalf("activity day = %v, want date bucket %s", activity.ActivityDay, today)
	}

	clearGuestSpamRedis(t)
	if err := LoadCacheFromDB(); err != nil {
		t.Fatalf("load cache from db: %v", err)
	}
	trust := TrustFor(integrationChatID, integrationCallerID)
	if !trust.Trusted || trust.RecentMessageCount < trustMessageCount {
		t.Fatalf("trust after reload = %+v, want trusted with reloaded activity", trust)
	}
}

func TestGuestSpamIntegrationTrackActivityHandle(t *testing.T) {
	setupGuestSpamIntegration(t)

	if err := TrackActivityHandle(tgbotapi.Update{Message: trackableMessage(integrationCallerID, "hello")}); err != nil {
		t.Fatalf("track activity: %v", err)
	}
	trust := TrustFor(integrationChatID, integrationCallerID)
	if trust.RecentMessageCount != 1 {
		t.Fatalf("recent count = %d, want 1", trust.RecentMessageCount)
	}
	if ActiveUserCount(integrationChatID) != 1 {
		t.Fatalf("active users = %d, want 1", ActiveUserCount(integrationChatID))
	}
}

func TestGuestSpamIntegrationLoadCacheOnlyRestoresFreshActiveUsers(t *testing.T) {
	db := setupGuestSpamIntegration(t)

	fresh := time.Now().Add(-2 * time.Minute)
	stale := time.Now().Add(-20 * time.Minute)
	if err := db.Create(&MemberRisk{
		ID:            riskID(integrationChatID, 9001),
		ChatID:        integrationChatID,
		UserID:        9001,
		UserName:      "Fresh",
		FirstSeenAt:   startOfDay(time.Now().AddDate(0, 0, -5)),
		LastMessageAt: fresh,
	}).Error; err != nil {
		t.Fatalf("create fresh risk: %v", err)
	}
	if err := db.Create(&MemberRisk{
		ID:            riskID(integrationChatID, 9002),
		ChatID:        integrationChatID,
		UserID:        9002,
		UserName:      "Stale",
		FirstSeenAt:   startOfDay(time.Now().AddDate(0, 0, -5)),
		LastMessageAt: stale,
	}).Error; err != nil {
		t.Fatalf("create stale risk: %v", err)
	}

	clearGuestSpamRedis(t)
	if err := LoadCacheFromDB(); err != nil {
		t.Fatalf("load cache from db: %v", err)
	}
	if got := ActiveUserCount(integrationChatID); got != 1 {
		t.Fatalf("active users after reload = %d, want 1", got)
	}
	if !IsActiveUser(integrationChatID, 9001) || IsActiveUser(integrationChatID, 9002) {
		t.Fatal("reload should keep fresh user only")
	}
	ttl, err := bot.GoRedis.TTL(redisCtx, activeUsersKey(integrationChatID)).Result()
	if err != nil {
		t.Fatalf("read active users ttl: %v", err)
	}
	if ttl != -1 {
		t.Fatalf("active users ttl after reload = %s, want no expiration", ttl)
	}
}

func TestGuestSpamIntegrationLoadCacheRestoresActiveUserAtCutoffBoundary(t *testing.T) {
	db := setupGuestSpamIntegration(t)

	now := time.Now()
	cutoffAt := time.Unix(int64(activeWindowCutoff(now)), 0)
	if err := db.Create(&MemberRisk{
		ID:            riskID(integrationChatID, 9003),
		ChatID:        integrationChatID,
		UserID:        9003,
		UserName:      "Boundary",
		FirstSeenAt:   startOfDay(now.AddDate(0, 0, -5)),
		LastMessageAt: cutoffAt,
	}).Error; err != nil {
		t.Fatalf("create cutoff risk: %v", err)
	}

	clearGuestSpamRedis(t)
	if err := LoadCacheFromDB(); err != nil {
		t.Fatalf("load cache from db: %v", err)
	}
	if !IsActiveUser(integrationChatID, 9003) {
		t.Fatal("member at cutoff boundary should be restored as active")
	}
	if got := ActiveUserCount(integrationChatID); got != 1 {
		t.Fatalf("active users after cutoff-boundary reload = %d, want 1", got)
	}
}

func TestGuestSpamIntegrationVotePassedBlacklistsAndClearsVote(t *testing.T) {
	setupGuestSpamIntegration(t)

	vote := SpamVote{
		ID:                "vote-integration",
		ChatID:            integrationChatID,
		ChatName:          "Guest Spam Test",
		MessageID:         401,
		GuestBotID:        993001,
		GuestBotName:      "Voted Spam Bot",
		GuestBotUserName:  "voted_spam_bot",
		ActiveUserCount:   5,
		RequiredVoteCount: 3,
		Voters:            []int64{1, 2, 3},
		CreatedAt:         time.Now(),
		ExpiresAt:         time.Now().Add(10 * time.Minute),
	}
	SaveVote(vote)
	ApplyVotePassedState(vote)

	if _, ok := GetVote(vote.ID); ok {
		t.Fatal("vote should be removed after passing")
	}
	if !IsBlacklisted(vote.GuestBotID) {
		t.Fatal("voted guest bot should be blacklisted")
	}

	if err := SyncCacheToDB(); err != nil {
		t.Fatalf("sync cache to db: %v", err)
	}
	clearGuestSpamRedis(t)
	if err := LoadCacheFromDB(); err != nil {
		t.Fatalf("load cache from db: %v", err)
	}
	if !IsBlacklisted(vote.GuestBotID) {
		t.Fatal("voted guest bot should reload as blacklisted")
	}
}

func TestGuestSpamIntegrationGuestSpamHandleCandidatesAndStartVote(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)
	message := commandMessage("/guest_spam", 7001)

	if err := GuestSpamHandle(tgbotapi.Update{Message: message}); err != nil {
		t.Fatalf("empty candidates command: %v", err)
	}
	if len(fake.sends) != 1 || !sentMessageContains(fake.sends[0], "最近没有可判定") {
		t.Fatalf("sends = %+v, want no candidates message", fake.sends)
	}

	recent := RecentGuestMessage{
		ChatID:           integrationChatID,
		ChatName:         "Guest Spam Test",
		MessageID:        501,
		GuestBotID:       993101,
		GuestBotName:     "Candidate Bot",
		GuestBotUserName: "candidate_bot",
		SeenAt:           time.Now(),
	}
	RecordRecentGuestMessage(recent)
	fake.sends = nil
	if err := GuestSpamHandle(tgbotapi.Update{Message: commandMessage("/guest_spam", 7001)}); err != nil {
		t.Fatalf("candidate list command: %v", err)
	}
	if len(fake.sends) != 1 || !sentMessageContains(fake.sends[0], "请选择") {
		t.Fatalf("sends = %+v, want candidate list", fake.sends)
	}

	for _, userID := range []int64{1, 2, 3} {
		writeActiveUsers(integrationChatID, time.Now(), userID)
	}
	fake.sends = nil
	if err := GuestSpamHandle(tgbotapi.Update{Message: commandMessage("/guest_spam @candidate_bot", 7001)}); err != nil {
		t.Fatalf("start vote command: %v", err)
	}
	if len(fake.sends) != 1 || !sentMessageContains(fake.sends[0], "是否将 guest bot") {
		t.Fatalf("sends = %+v, want vote message", fake.sends)
	}
	if len(voteIDs(t)) != 1 {
		t.Fatalf("vote dirty ids = %v, want one vote", voteIDs(t))
	}
}

func TestGuestSpamIntegrationGuestSpamHandleSendFailures(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)
	fake.sendErr = errTelegram()

	err := GuestSpamHandle(tgbotapi.Update{Message: commandMessage("/guest_spam", 7001)})
	if err == nil {
		t.Fatal("empty candidates send failure should return error")
	}
	if len(fake.queuedDeletes) != 1 {
		t.Fatalf("queued deletes = %d, want command cleanup queued", len(fake.queuedDeletes))
	}

	RecordRecentGuestMessage(RecentGuestMessage{
		ChatID:           integrationChatID,
		ChatName:         "Guest Spam Test",
		MessageID:        505,
		GuestBotID:       993105,
		GuestBotName:     "Send Fail Bot",
		GuestBotUserName: "send_fail_bot",
		SeenAt:           time.Now(),
	})
	writeActiveUsers(integrationChatID, time.Now(), 1, 2, 3)
	err = GuestSpamHandle(tgbotapi.Update{Message: commandMessage("/guest_spam 505", 7001)})
	if err == nil {
		t.Fatal("start vote send failure should return error")
	}
	if len(voteIDs(t)) != 0 {
		t.Fatalf("vote dirty ids = %v, want no vote saved on send failure", voteIDs(t))
	}
	if hasLogAction(RecentLogs(integrationChatID, 10), ActionVoteStarted) {
		t.Fatalf("logs = %+v, want no vote started log on send failure", RecentLogs(integrationChatID, 10))
	}
}

func TestGuestSpamIntegrationGuestSpamHandleInsufficientActiveUsers(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)
	RecordRecentGuestMessage(RecentGuestMessage{
		ChatID:           integrationChatID,
		ChatName:         "Guest Spam Test",
		MessageID:        502,
		GuestBotID:       993102,
		GuestBotName:     "Quiet Bot",
		GuestBotUserName: "quiet_bot",
		SeenAt:           time.Now(),
	})

	if err := GuestSpamHandle(tgbotapi.Update{Message: commandMessage("/guest_spam 502", 7001)}); err != nil {
		t.Fatalf("insufficient active command: %v", err)
	}
	if len(fake.sends) != 1 || !sentMessageContains(fake.sends[0], "活跃人数少于 3") {
		t.Fatalf("sends = %+v, want invalid vote message", fake.sends)
	}
	if !hasLogAction(RecentLogs(integrationChatID, 10), ActionVoteInvalid) {
		t.Fatalf("logs = %+v, want vote invalid", RecentLogs(integrationChatID, 10))
	}
	if len(voteIDs(t)) != 0 {
		t.Fatalf("vote dirty ids = %v, want none", voteIDs(t))
	}
}

func TestGuestSpamIntegrationSelectRecentGuestCallbackPaths(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)
	callback := selectCallback("select-expired", 503, 7001)

	if err := SelectRecentGuestCallback(tgbotapi.Update{CallbackQuery: callback}); err != nil {
		t.Fatalf("expired select callback: %v", err)
	}
	if len(fake.callbacks) != 1 || !fake.callbacks[0].showAlert || fake.callbacks[0].text != "候选消息已过期" {
		t.Fatalf("callbacks = %+v, want expired alert", fake.callbacks)
	}

	RecordRecentGuestMessage(RecentGuestMessage{
		ChatID:           integrationChatID,
		ChatName:         "Guest Spam Test",
		MessageID:        504,
		GuestBotID:       993104,
		GuestBotName:     "Selected Bot",
		GuestBotUserName: "selected_bot",
		SeenAt:           time.Now(),
	})
	fake.callbacks = nil
	fake.callbackDelete = nil
	fake.sends = nil
	if err := SelectRecentGuestCallback(tgbotapi.Update{CallbackQuery: selectCallback("select-low-active", 504, 7001)}); err != nil {
		t.Fatalf("low active select callback: %v", err)
	}
	if len(fake.callbacks) != 0 || len(fake.callbackDelete) != 0 {
		t.Fatalf("low active should not answer started/delete callback, callbacks=%+v deletes=%+v", fake.callbacks, fake.callbackDelete)
	}
	if len(fake.sends) != 1 || !sentMessageContains(fake.sends[0], "活跃人数少于 3") {
		t.Fatalf("low active sends=%+v, want invalid vote message", fake.sends)
	}

	for _, userID := range []int64{1, 2, 3} {
		writeActiveUsers(integrationChatID, time.Now(), userID)
	}
	fake.callbacks = nil
	fake.callbackDelete = nil
	if err := SelectRecentGuestCallback(tgbotapi.Update{CallbackQuery: selectCallback("select-ok", 504, 7001)}); err != nil {
		t.Fatalf("valid select callback: %v", err)
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "已发起投票" || len(fake.callbackDelete) != 1 {
		t.Fatalf("callbacks=%+v deletes=%+v, want start answer and callback delete", fake.callbacks, fake.callbackDelete)
	}
}

func TestGuestSpamIntegrationSpamVoteCallbackPaths(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	if err := SpamVoteCallback(tgbotapi.Update{}); err != nil {
		t.Fatalf("nil callback: %v", err)
	}
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: &tgbotapi.CallbackQuery{Data: "bad"}}); err != nil {
		t.Fatalf("bad callback data: %v", err)
	}

	SaveVote(SpamVote{ID: "unknown-action", ExpiresAt: time.Now().Add(time.Hour), RequiredVoteCount: 1})
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("unknown-cb", "guestspam_vote,wat,unknown-action", 7001)}); err != nil {
		t.Fatalf("unknown action callback: %v", err)
	}
	vote, ok := GetVote("unknown-action")
	if !ok || len(vote.Voters) != 0 {
		t.Fatalf("unknown action vote = %+v ok=%v, want unchanged", vote, ok)
	}

	SaveVote(SpamVote{ID: "cancel-vote", ChatID: integrationChatID, StarterUserID: 7001, ExpiresAt: time.Now().Add(time.Hour)})
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("cancel-cb", "guestspam_vote,cancel,cancel-vote", 7001)}); err != nil {
		t.Fatalf("cancel callback: %v", err)
	}
	if _, ok := GetVote("cancel-vote"); ok {
		t.Fatal("cancel should delete vote")
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "已取消" || len(fake.callbackDelete) != 1 {
		t.Fatalf("cancel callbacks=%+v deletes=%+v", fake.callbacks, fake.callbackDelete)
	}

	SaveVote(SpamVote{ID: "cancel-denied", ChatID: integrationChatID, StarterUserID: 7001, ExpiresAt: time.Now().Add(time.Hour)})
	fake.callbacks = nil
	fake.callbackDelete = nil
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("cancel-denied-cb", "guestspam_vote,cancel,cancel-denied", 7002)}); err != nil {
		t.Fatalf("cancel denied callback: %v", err)
	}
	if _, ok := GetVote("cancel-denied"); !ok {
		t.Fatal("unauthorized cancel should keep vote")
	}
	if len(fake.callbacks) != 1 || !fake.callbacks[0].showAlert || fake.callbacks[0].text != "只有发起者或管理员可以取消投票" || len(fake.callbackDelete) != 0 {
		t.Fatalf("cancel denied callbacks=%+v deletes=%+v", fake.callbacks, fake.callbackDelete)
	}

	fake.admins[7002] = true
	fake.callbacks = nil
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("cancel-admin-cb", "guestspam_vote,cancel,cancel-denied", 7002)}); err != nil {
		t.Fatalf("admin cancel callback: %v", err)
	}
	if _, ok := GetVote("cancel-denied"); ok {
		t.Fatal("admin cancel should delete vote")
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "已取消" || len(fake.callbackDelete) != 1 {
		t.Fatalf("admin cancel callbacks=%+v deletes=%+v", fake.callbacks, fake.callbackDelete)
	}
	delete(fake.admins, 7002)

	fake.callbacks = nil
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("expired-cb", "guestspam_vote,vote,missing-vote", 7001)}); err != nil {
		t.Fatalf("expired callback: %v", err)
	}
	if len(fake.callbacks) != 1 || !fake.callbacks[0].showAlert || fake.callbacks[0].text != "投票已过期" {
		t.Fatalf("expired callbacks=%+v", fake.callbacks)
	}

	SaveVote(SpamVote{ID: "bot-vote", ExpiresAt: time.Now().Add(time.Hour), RequiredVoteCount: 2})
	fake.callbacks = nil
	botVoteCallback := voteCallback("bot-cb", "guestspam_vote,vote,bot-vote", 0)
	botVoteCallback.From.IsBot = true
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: botVoteCallback}); err != nil {
		t.Fatalf("bot vote callback: %v", err)
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "机器人不能参与投票" {
		t.Fatalf("bot callbacks=%+v", fake.callbacks)
	}

	SaveVote(SpamVote{ID: "dup-vote", ExpiresAt: time.Now().Add(time.Hour), RequiredVoteCount: 2, Voters: []int64{7001}})
	fake.callbacks = nil
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("dup-cb", "guestspam_vote,vote,dup-vote", 7001)}); err != nil {
		t.Fatalf("duplicate vote callback: %v", err)
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "你已经投过票了" {
		t.Fatalf("duplicate callbacks=%+v", fake.callbacks)
	}

	writeActiveUsers(integrationChatID, time.Now(), 7001)
	SaveVote(SpamVote{ID: "partial-vote", ChatID: integrationChatID, ExpiresAt: time.Now().Add(time.Hour), RequiredVoteCount: 2})
	fake.callbacks = nil
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("partial-cb", "guestspam_vote,vote,partial-vote", 7001)}); err != nil {
		t.Fatalf("partial vote callback: %v", err)
	}
	vote, ok = GetVote("partial-vote")
	if !ok || len(vote.Voters) != 1 || vote.Voters[0] != 7001 || vote.VoteScore != activeVoteWeight {
		t.Fatalf("partial vote = %+v ok=%v, want one voter", vote, ok)
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "已投票 1/2（活跃成员 +1）" {
		t.Fatalf("partial callbacks=%+v", fake.callbacks)
	}

	writeActiveUsers(integrationChatID, time.Now(), 7002)
	SaveVote(SpamVote{
		ID:                "pass-vote",
		ChatID:            integrationChatID,
		ChatName:          "Guest Spam Test",
		MessageID:         601,
		GuestBotID:        993201,
		GuestBotName:      "Pass Bot",
		GuestBotUserName:  "pass_bot",
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: 2,
		Voters:            []int64{7001},
	})
	fake.callbacks = nil
	fake.callbackDelete = nil
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("pass-cb", "guestspam_vote,vote,pass-vote", 7002)}); err != nil {
		t.Fatalf("pass vote callback: %v", err)
	}
	if _, ok := GetVote("pass-vote"); ok {
		t.Fatal("passing vote should delete vote")
	}
	if !IsBlacklisted(993201) {
		t.Fatal("passing vote should blacklist guest bot")
	}
	if len(fake.deletes) != 1 || len(fake.callbacks) != 1 || fake.callbacks[0].text != "投票通过，已拉黑并删除消息" || len(fake.callbackDelete) != 1 {
		t.Fatalf("pass deletes=%+v callbacks=%+v callbackDeletes=%+v", fake.deletes, fake.callbacks, fake.callbackDelete)
	}
}

func TestGuestSpamIntegrationSpamVoteWeightTiers(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	SaveVote(SpamVote{
		ID:                "inactive-vote",
		ChatID:            integrationChatID,
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: 1,
	})
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("inactive-cb", "guestspam_vote,vote,inactive-vote", 7101)}); err != nil {
		t.Fatalf("inactive vote callback: %v", err)
	}
	vote, ok := GetVote("inactive-vote")
	if !ok || len(vote.Voters) != 0 || vote.VoteScore != 0 {
		t.Fatalf("inactive vote = %+v ok=%v, want no voter and zero score", vote, ok)
	}
	if len(fake.callbacks) != 1 || !fake.callbacks[0].showAlert || fake.callbacks[0].text != "最近不够活跃，不能参与投票" {
		t.Fatalf("inactive callbacks=%+v", fake.callbacks)
	}

	writeActiveUsers(integrationChatID, time.Now(), 7102)
	seedTrustedCaller(7103)
	fake.admins[7104] = true
	fake.callbacks = nil

	SaveVote(SpamVote{
		ID:                "weighted-vote",
		ChatID:            integrationChatID,
		ChatName:          "Guest Spam Test",
		MessageID:         701,
		GuestBotID:        993301,
		GuestBotName:      "Weighted Bot",
		GuestBotUserName:  "weighted_bot",
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: activeVoteWeight + trustedVoteWeight + adminVoteWeight,
	})

	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("active-cb", "guestspam_vote,vote,weighted-vote", 7102)}); err != nil {
		t.Fatalf("active vote callback: %v", err)
	}
	vote, ok = GetVote("weighted-vote")
	if !ok || vote.VoteScore != activeVoteWeight {
		t.Fatalf("after active vote = %+v ok=%v, want score %d", vote, ok, activeVoteWeight)
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "已投票 1/6（活跃成员 +1）" {
		t.Fatalf("active callbacks=%+v", fake.callbacks)
	}

	fake.callbacks = nil
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("trusted-cb", "guestspam_vote,vote,weighted-vote", 7103)}); err != nil {
		t.Fatalf("trusted vote callback: %v", err)
	}
	vote, ok = GetVote("weighted-vote")
	if !ok || vote.VoteScore != activeVoteWeight+trustedVoteWeight {
		t.Fatalf("after trusted vote = %+v ok=%v, want score %d", vote, ok, activeVoteWeight+trustedVoteWeight)
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "已投票 3/6（可信成员 +2）" {
		t.Fatalf("trusted callbacks=%+v", fake.callbacks)
	}

	fake.callbacks = nil
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("admin-cb", "guestspam_vote,vote,weighted-vote", 7104)}); err != nil {
		t.Fatalf("admin vote callback: %v", err)
	}
	if _, ok := GetVote("weighted-vote"); ok {
		t.Fatal("admin weighted vote should pass and delete vote")
	}
	if !IsBlacklisted(993301) {
		t.Fatal("weighted vote should blacklist guest bot")
	}
	if len(fake.deletes) != 1 || len(fake.callbacks) != 1 || fake.callbacks[0].text != "投票通过，已拉黑并删除消息" {
		t.Fatalf("admin pass deletes=%+v callbacks=%+v", fake.deletes, fake.callbacks)
	}
}

func TestGuestSpamIntegrationGuestSpamLogHandlePaths(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)
	nonAdmin := commandMessage("/guest_spam_log", 8001)
	if err := GuestSpamLogHandle(tgbotapi.Update{Message: nonAdmin}); err != nil {
		t.Fatalf("non-admin log command: %v", err)
	}
	if len(fake.sends) != 1 || !sentMessageContains(fake.sends[0], "无使用权限") {
		t.Fatalf("non-admin sends=%+v", fake.sends)
	}

	fake.admins[8001] = true
	fake.sends = nil
	if err := GuestSpamLogHandle(tgbotapi.Update{Message: commandMessage("/guest_spam_log", 8001)}); err != nil {
		t.Fatalf("empty log command: %v", err)
	}
	if len(fake.sends) != 1 || !sentMessageContains(fake.sends[0], "暂无 guest spam 日志") {
		t.Fatalf("empty log sends=%+v", fake.sends)
	}

	fake.sends = nil
	if err := GuestSpamLogHandle(tgbotapi.Update{Message: commandMessage("/guest_spam_log restore bad", 8001)}); err != nil {
		t.Fatalf("bad restore command: %v", err)
	}
	if len(fake.sends) != 1 || !sentMessageContains(fake.sends[0], "用户 ID 格式错误") {
		t.Fatalf("bad restore sends=%+v", fake.sends)
	}

	AddWarning(integrationChatID, integrationCallerID, "Caller")
	fake.sends = nil
	if err := GuestSpamLogHandle(tgbotapi.Update{Message: commandMessage("/guest_spam_log restore 880001", 8001)}); err != nil {
		t.Fatalf("restore command: %v", err)
	}
	risk, ok := getMemberRisk(integrationChatID, integrationCallerID)
	if !ok || risk.WarningCount != 0 || risk.MuteLevel != 0 {
		t.Fatalf("risk after restore = %+v ok=%v, want cleared", risk, ok)
	}
	if len(fake.unbans) != 1 || len(fake.restricts) != 1 || len(fake.sends) != 1 || !sentMessageContains(fake.sends[0], "已恢复") {
		t.Fatalf("restore unbans=%+v restricts=%+v sends=%+v", fake.unbans, fake.restricts, fake.sends)
	}
	if fake.restricts[0].permissions != tgbotapi.AllPermissions {
		t.Fatalf("restore permissions = %q, want %q", fake.restricts[0].permissions, tgbotapi.AllPermissions)
	}

	AddWarning(integrationChatID, integrationCallerID, "Caller")
	risk, _ = getMemberRisk(integrationChatID, integrationCallerID)
	risk.MuteLevel = 2
	setMemberRisk(risk, true)
	fake.unbans = nil
	fake.restricts = nil
	fake.restrictErr = errTelegram()
	fake.sends = nil
	err := GuestSpamLogHandle(tgbotapi.Update{Message: commandMessage("/guest_spam_log restore 880001", 8001)})
	if err == nil {
		t.Fatal("restore should return telegram error when unrestrict fails")
	}
	risk, ok = getMemberRisk(integrationChatID, integrationCallerID)
	if !ok || risk.WarningCount == 0 || risk.MuteLevel == 0 {
		t.Fatalf("failed restore risk = %+v ok=%v, want warnings kept", risk, ok)
	}
	if len(fake.sends) != 0 {
		t.Fatalf("failed restore should not send success message, sends=%+v", fake.sends)
	}
	if len(fake.unbans) != 1 || len(fake.restricts) != 1 {
		t.Fatalf("failed unrestrict should still attempt unban+restrict once, unbans=%+v restricts=%+v", fake.unbans, fake.restricts)
	}
	logs := RecentLogs(integrationChatID, 10)
	if len(logs) == 0 || !strings.Contains(logs[0].Detail, "unrestrict failed:") {
		t.Fatalf("logs = %+v, want latest detail to contain unrestrict failed:", logs)
	}

	AddWarning(integrationChatID, integrationCallerID, "Caller")
	risk, _ = getMemberRisk(integrationChatID, integrationCallerID)
	risk.MuteLevel = 2
	setMemberRisk(risk, true)
	fake.unbanErr = errTelegram()
	fake.restrictErr = nil
	fake.unbans = nil
	fake.restricts = nil
	fake.sends = nil
	err = GuestSpamLogHandle(tgbotapi.Update{Message: commandMessage("/guest_spam_log restore 880001", 8001)})
	if err == nil {
		t.Fatal("restore should return telegram error when unban fails")
	}
	risk, ok = getMemberRisk(integrationChatID, integrationCallerID)
	if !ok || risk.WarningCount == 0 || risk.MuteLevel == 0 {
		t.Fatalf("failed unban risk = %+v ok=%v, want warnings kept", risk, ok)
	}
	if len(fake.unbans) != 1 || len(fake.restricts) != 0 {
		t.Fatalf("failed unban should stop before restrict, unbans=%+v restricts=%+v", fake.unbans, fake.restricts)
	}
	if len(fake.sends) != 0 {
		t.Fatalf("failed unban should not send success message, sends=%+v", fake.sends)
	}
	logs = RecentLogs(integrationChatID, 10)
	if len(logs) == 0 || !strings.Contains(logs[0].Detail, "unban failed:") {
		t.Fatalf("logs = %+v, want latest detail to contain unban failed:", logs)
	}

	fake.sends = nil
	fake.unbanErr = nil
	fake.restrictErr = nil
	if err := GuestSpamLogHandle(tgbotapi.Update{Message: commandMessage("/guest_spam_log", 8001)}); err != nil {
		t.Fatalf("log command: %v", err)
	}
	if len(fake.sends) != 1 || !sentMessageContains(fake.sends[0], "最近 guest spam 日志") {
		t.Fatalf("log sends=%+v", fake.sends)
	}
}

func TestGuestSpamIntegrationRestoreClearsWarnings(t *testing.T) {
	setupGuestSpamIntegration(t)

	risk, _ := AddWarning(integrationChatID, integrationCallerID, "Caller")
	risk.WarningCount = 1
	risk.MuteLevel = 2
	setMemberRisk(risk, true)

	RestoreCallerState(integrationChatID, integrationCallerID, &tgbotapi.Message{
		Chat: &tgbotapi.Chat{
			ID:    integrationChatID,
			Title: "Guest Spam Test",
		},
		From: &tgbotapi.User{
			ID:        42,
			FirstName: "Admin",
		},
	})

	risk, ok := getMemberRisk(integrationChatID, integrationCallerID)
	if !ok {
		t.Fatal("member risk should still exist")
	}
	if risk.WarningCount != 0 || risk.MuteLevel != 0 {
		t.Fatalf("restored risk = %+v, want cleared warnings and mute level", risk)
	}
}

func setupGuestSpamIntegration(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("GUEST_SPAM_TEST_MYSQL_DSN")
	redisAddr := os.Getenv("GUEST_SPAM_TEST_REDIS_ADDR")
	if dsn == "" || redisAddr == "" {
		t.Skip("set GUEST_SPAM_TEST_MYSQL_DSN and GUEST_SPAM_TEST_REDIS_ADDR to run integration tests")
		return nil
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open mysql: %v", err)
	}
	bot.DBEngine = db

	bot.GoRedis = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: os.Getenv("GUEST_SPAM_TEST_REDIS_PASSWORD"),
		DB:       0,
	})
	if err := bot.GoRedis.Ping(context.Background()).Err(); err != nil {
		t.Fatalf("ping redis: %v", err)
	}
	t.Cleanup(func() {
		clearGuestSpamRedis(t)
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
		_ = bot.GoRedis.Close()
		bot.DBEngine = nil
		bot.GoRedis = nil
	})

	migrateGuestSpamTestSchema(t, db)
	clearGuestSpamTables(t, db)
	clearGuestSpamRedis(t)
	return db
}

func migrateGuestSpamTestSchema(t *testing.T, db *gorm.DB) {
	t.Helper()
	// The integration suite owns a disposable schema; production migrations stay manual.
	if err := db.AutoMigrate(&MemberRisk{}, &MemberActivity{}, &GuestBotBlacklist{}, &SpamLog{}); err != nil {
		t.Fatalf("auto migrate guest spam schema: %v", err)
	}
}

func clearGuestSpamTables(t *testing.T, db *gorm.DB) {
	t.Helper()
	tables := []string{
		"guest_spam_log",
		"guest_spam_bot_blacklist",
		"guest_spam_member_activity",
		"guest_spam_member_risk",
	}
	for _, table := range tables {
		if err := db.Exec("DELETE FROM " + table).Error; err != nil {
			t.Fatalf("clear %s: %v", table, err)
		}
	}
}

func clearGuestSpamRedis(t *testing.T) {
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

func testGuestMessage(guestBotID, callerID int64) *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: int(guestBotID % 100000),
		Chat: &tgbotapi.Chat{
			ID:    integrationChatID,
			Type:  "supergroup",
			Title: "Guest Spam Test",
		},
		From: &tgbotapi.User{
			ID:        guestBotID,
			IsBot:     true,
			FirstName: "Guest",
			UserName:  fmt.Sprintf("guest_%d_bot", guestBotID),
		},
		GuestBotCallerUser: &tgbotapi.User{
			ID:        callerID,
			FirstName: "Caller",
			UserName:  fmt.Sprintf("caller_%d", callerID),
		},
	}
}

func testBotCallerMessage(guestBotID, callerID int64) *tgbotapi.Message {
	message := testGuestMessage(guestBotID, callerID)
	message.GuestBotCallerUser.IsBot = true
	return message
}

func testGuestChatMessage(guestBotID, callerChatID int64) *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: int(guestBotID % 100000),
		Chat: &tgbotapi.Chat{
			ID:    integrationChatID,
			Type:  "supergroup",
			Title: "Guest Spam Test",
		},
		From: &tgbotapi.User{
			ID:        guestBotID,
			IsBot:     true,
			FirstName: "Guest",
			UserName:  fmt.Sprintf("guest_%d_bot", guestBotID),
		},
		GuestBotCallerChat: &tgbotapi.Chat{
			ID:    callerChatID,
			Type:  "channel",
			Title: "Caller Channel",
		},
	}
}

func seedTrustedCaller(userID int64) {
	firstSeen := startOfDay(time.Now().Add(-4 * 24 * time.Hour))
	risk := MemberRisk{
		ID:                 riskID(integrationChatID, userID),
		ChatID:             integrationChatID,
		UserID:             userID,
		UserName:           "Trusted Caller",
		FirstSeenAt:        firstSeen,
		LastMessageAt:      time.Now(),
		RecentMessageCount: trustMessageCount,
	}
	setMemberRisk(risk, true)
	for i := int64(0); i < trustMessageCount; i++ {
		key := memberActivityDayKey(integrationChatID, userID, dayKey(time.Now()))
		bot.GoRedis.Incr(redisCtx, key)
		bot.GoRedis.Expire(redisCtx, key, activityDayTTL)
	}
}

func trackableMessage(userID int64, text string) *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: int(userID % 100000),
		Chat: &tgbotapi.Chat{
			ID:    integrationChatID,
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

func commandMessage(text string, userID int64) *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: int(userID % 100000),
		Chat: &tgbotapi.Chat{
			ID:    integrationChatID,
			Type:  "supergroup",
			Title: "Guest Spam Test",
		},
		From: &tgbotapi.User{
			ID:        userID,
			FirstName: "Admin",
			UserName:  fmt.Sprintf("admin_%d", userID),
		},
		Text:     text,
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: len(strings.Fields(text)[0])}},
	}
}

func selectCallback(id string, messageID int, userID int64) *tgbotapi.CallbackQuery {
	return &tgbotapi.CallbackQuery{
		ID:   id,
		Data: fmt.Sprintf("guestspam_select,%d", messageID),
		From: &tgbotapi.User{
			ID:        userID,
			FirstName: "Voter",
		},
		Message: &tgbotapi.Message{
			MessageID: 9001,
			Chat: &tgbotapi.Chat{
				ID:    integrationChatID,
				Type:  "supergroup",
				Title: "Guest Spam Test",
			},
		},
	}
}

func voteCallback(id string, data string, userID int64) *tgbotapi.CallbackQuery {
	return &tgbotapi.CallbackQuery{
		ID:   id,
		Data: data,
		From: &tgbotapi.User{
			ID:        userID,
			FirstName: "Voter",
		},
		Message: &tgbotapi.Message{
			MessageID: 9002,
			Chat: &tgbotapi.Chat{
				ID:    integrationChatID,
				Type:  "supergroup",
				Title: "Guest Spam Test",
			},
		},
	}
}

func sentMessageContains(config tgbotapi.Chattable, text string) bool {
	message, ok := config.(tgbotapi.MessageConfig)
	return ok && strings.Contains(message.Text, text)
}

func voteIDs(t *testing.T) []string {
	t.Helper()
	ids, err := bot.GoRedis.SMembers(redisCtx, voteDirtySetKey()).Result()
	if err != nil {
		t.Fatalf("read vote dirty set: %v", err)
	}
	return ids
}
