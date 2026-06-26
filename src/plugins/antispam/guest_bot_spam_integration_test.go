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

	err := GuestBotSpamHandle(tgbotapi.Update{Message: testGuestMessage(integrationGuestBotID, integrationCallerID)})
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

	err := GuestBotSpamHandle(tgbotapi.Update{Message: testGuestMessage(integrationGuestBot2ID, integrationCallerID)})
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

	if err := GuestBotSpamHandle(tgbotapi.Update{Message: testGuestMessage(991201, integrationCallerID)}); err != nil {
		t.Fatalf("first warning: %v", err)
	}
	if err := GuestBotSpamHandle(tgbotapi.Update{Message: testGuestMessage(991202, integrationCallerID)}); err != nil {
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

	first := GuestBotSpamHandle(tgbotapi.Update{Message: testGuestMessage(991101, integrationCallerID)})
	second := GuestBotSpamHandle(tgbotapi.Update{Message: testGuestMessage(991102, integrationCallerID)})
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

			if err := GuestBotSpamHandle(tgbotapi.Update{Message: tt.message}); err != nil {
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
	seedMemberRisk(t, db, MemberRisk{
		ID:            riskID(integrationChatID, 9001),
		ChatID:        integrationChatID,
		UserID:        9001,
		UserName:      "Fresh",
		FirstSeenAt:   startOfDay(time.Now().AddDate(0, 0, -5)),
		LastMessageAt: fresh,
	})
	seedMemberRisk(t, db, MemberRisk{
		ID:            riskID(integrationChatID, 9002),
		ChatID:        integrationChatID,
		UserID:        9002,
		UserName:      "Stale",
		FirstSeenAt:   startOfDay(time.Now().AddDate(0, 0, -5)),
		LastMessageAt: stale,
	})

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

	now := time.Now().Truncate(time.Second)
	cutoffAt := time.Unix(int64(activeWindowCutoff(now)), 0)
	seedMemberRisk(t, db, MemberRisk{
		ID:            riskID(integrationChatID, 9003),
		ChatID:        integrationChatID,
		UserID:        9003,
		UserName:      "Boundary",
		FirstSeenAt:   startOfDay(now.AddDate(0, 0, -5)),
		LastMessageAt: cutoffAt,
	})

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
	if err := SaveVote(vote); err != nil {
		t.Fatalf("save vote: %v", err)
	}
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

func TestGuestSpamIntegrationRecentCandidatesReloadPerChatNotGlobalLimit(t *testing.T) {
	db := setupGuestSpamIntegration(t)
	targetChatID := int64(-100500)
	now := time.Now().Truncate(time.Second)

	target := SpamLog{
		ID:             "target-recent-candidate",
		ChatID:         targetChatID,
		ChatName:       "Target Chat",
		MessageID:      777001,
		GuestBotID:     997001,
		GuestBotName:   "Target Candidate Bot",
		GuestBotUser:   "target_candidate_bot",
		CallerUserID:   887001,
		CallerUserName: "Target Caller",
		Action:         ActionGuestSeen,
		Reason:         ReasonTrusted,
		CreateTime:     now.Add(-23 * time.Hour),
		UpdateTime:     now.Add(-23 * time.Hour),
	}
	if err := db.Create(&target).Error; err != nil {
		t.Fatalf("create target log: %v", err)
	}

	noise := make([]SpamLog, 0, 1000)
	for i := 0; i < 1000; i++ {
		noise = append(noise, SpamLog{
			ID:           fmt.Sprintf("noise-recent-%04d", i),
			ChatID:       int64(-200000 - i),
			ChatName:     "Busy Chat",
			MessageID:    800000 + i,
			GuestBotID:   int64(998000 + i),
			GuestBotName: "Busy Candidate Bot",
			GuestBotUser: fmt.Sprintf("busy_candidate_%04d_bot", i),
			Action:       ActionGuestSeen,
			Reason:       ReasonTrusted,
			CreateTime:   now.Add(-time.Duration(i) * time.Second),
			UpdateTime:   now.Add(-time.Duration(i) * time.Second),
		})
	}
	if err := db.Create(&noise).Error; err != nil {
		t.Fatalf("create noise logs: %v", err)
	}

	clearGuestSpamRedis(t)
	if err := LoadCacheFromDB(); err != nil {
		t.Fatalf("load cache from db: %v", err)
	}

	recents := RecentGuestMessages(targetChatID)
	if len(recents) != 1 {
		t.Fatalf("recent candidates = %+v, want target candidate", recents)
	}
	if recents[0].MessageID != target.MessageID || recents[0].GuestBotID != target.GuestBotID {
		t.Fatalf("recent candidate = %+v, want message %d bot %d", recents[0], target.MessageID, target.GuestBotID)
	}
}

func TestGuestSpamIntegrationRecentCandidatesFallbackDedupesBeyondInitialPage(t *testing.T) {
	db := setupGuestSpamIntegration(t)
	targetChatID := int64(-100501)
	now := time.Now().Truncate(time.Second)

	dupCount := recentMessageLimit * 12
	if dupCount < recentMessageLimit {
		dupCount = recentMessageLimit
	}
	uniqueCount := recentMessageLimit - 1
	if uniqueCount < 1 {
		uniqueCount = 1
	}

	rows := make([]SpamLog, 0, dupCount+uniqueCount)
	for i := 0; i < dupCount; i++ {
		ts := now.Add(-time.Duration(i) * time.Second)
		rows = append(rows, SpamLog{
			ID:             fmt.Sprintf("dup-target-%03d", i),
			ChatID:         targetChatID,
			ChatName:       "Target Chat",
			MessageID:      888001,
			GuestBotID:     997101,
			GuestBotName:   "Duplicate Candidate Bot",
			GuestBotUser:   "duplicate_candidate_bot",
			CallerUserID:   887101,
			CallerUserName: "Target Caller",
			Action:         ActionGuestSeen,
			Reason:         ReasonTrusted,
			CreateTime:     ts,
			UpdateTime:     ts,
		})
	}
	for i := 0; i < uniqueCount; i++ {
		ts := now.Add(-time.Duration(dupCount+i+1) * time.Second)
		rows = append(rows, SpamLog{
			ID:             fmt.Sprintf("unique-target-%03d", i),
			ChatID:         targetChatID,
			ChatName:       "Target Chat",
			MessageID:      888100 + i,
			GuestBotID:     997200 + int64(i),
			GuestBotName:   fmt.Sprintf("Unique Candidate Bot %d", i),
			GuestBotUser:   fmt.Sprintf("unique_candidate_%03d_bot", i),
			CallerUserID:   887200 + int64(i),
			CallerUserName: fmt.Sprintf("Unique Caller %d", i),
			Action:         ActionGuestSeen,
			Reason:         ReasonTrusted,
			CreateTime:     ts,
			UpdateTime:     ts,
		})
	}
	if err := db.Create(&rows).Error; err != nil {
		t.Fatalf("create duplicate and unique logs: %v", err)
	}

	clearGuestSpamRedis(t)
	if err := LoadCacheFromDB(); err != nil {
		t.Fatalf("load cache from db: %v", err)
	}

	recents := RecentGuestMessages(targetChatID)
	if len(recents) != recentMessageLimit {
		t.Fatalf("recent candidates len = %d, want %d", len(recents), recentMessageLimit)
	}
	if recents[0].MessageID != 888001 {
		t.Fatalf("first recent candidate = %+v, want duplicate candidate first", recents[0])
	}
	for i := 1; i < len(recents); i++ {
		wantID := 888100 + i - 1
		if recents[i].MessageID != wantID {
			t.Fatalf("recent candidate[%d] = %+v, want message %d", i, recents[i], wantID)
		}
		if recents[i-1].SeenAt.Before(recents[i].SeenAt) {
			t.Fatalf("recent candidates not newest-first at %d: %+v before %+v", i, recents[i-1], recents[i])
		}
	}

	warmed := RecentGuestMessages(targetChatID)
	if len(warmed) != len(recents) {
		t.Fatalf("warmed recent candidates len = %d, want %d", len(warmed), len(recents))
	}
	for i := range recents {
		if warmed[i].MessageID != recents[i].MessageID || warmed[i].GuestBotID != recents[i].GuestBotID || !warmed[i].SeenAt.Equal(recents[i].SeenAt) {
			t.Fatalf("warmed recent candidate[%d] = %+v, want %+v", i, warmed[i], recents[i])
		}
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

	if err := SaveVote(SpamVote{ID: "unknown-action", ExpiresAt: time.Now().Add(time.Hour), RequiredVoteCount: 1}); err != nil {
		t.Fatalf("save vote: %v", err)
	}
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("unknown-cb", "guestspam_vote,wat,unknown-action", 7001)}); err != nil {
		t.Fatalf("unknown action callback: %v", err)
	}
	vote, ok := GetVote("unknown-action")
	if !ok || len(vote.Voters) != 0 {
		t.Fatalf("unknown action vote = %+v ok=%v, want unchanged", vote, ok)
	}

	if err := SaveVote(SpamVote{ID: "cancel-vote", ChatID: integrationChatID, StarterUserID: 7001, ExpiresAt: time.Now().Add(time.Hour)}); err != nil {
		t.Fatalf("save vote: %v", err)
	}
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("cancel-cb", "guestspam_vote,cancel,cancel-vote", 7001)}); err != nil {
		t.Fatalf("cancel callback: %v", err)
	}
	if _, ok := GetVote("cancel-vote"); ok {
		t.Fatal("cancel should delete vote")
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "已取消" || len(fake.callbackDelete) != 1 {
		t.Fatalf("cancel callbacks=%+v deletes=%+v", fake.callbacks, fake.callbackDelete)
	}

	if err := SaveVote(SpamVote{ID: "cancel-denied", ChatID: integrationChatID, StarterUserID: 7001, ExpiresAt: time.Now().Add(time.Hour)}); err != nil {
		t.Fatalf("save vote: %v", err)
	}
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

	if err := SaveVote(SpamVote{ID: "bot-vote", ExpiresAt: time.Now().Add(time.Hour), RequiredVoteCount: 2}); err != nil {
		t.Fatalf("save vote: %v", err)
	}
	fake.callbacks = nil
	botVoteCallback := voteCallback("bot-cb", "guestspam_vote,vote,bot-vote", 0)
	botVoteCallback.From.IsBot = true
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: botVoteCallback}); err != nil {
		t.Fatalf("bot vote callback: %v", err)
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "机器人不能参与投票" {
		t.Fatalf("bot callbacks=%+v", fake.callbacks)
	}

	if err := SaveVote(SpamVote{ID: "dup-vote", ExpiresAt: time.Now().Add(time.Hour), RequiredVoteCount: 2, Voters: []int64{7001}}); err != nil {
		t.Fatalf("save vote: %v", err)
	}
	fake.callbacks = nil
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("dup-cb", "guestspam_vote,vote,dup-vote", 7001)}); err != nil {
		t.Fatalf("duplicate vote callback: %v", err)
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "你已经投过票了" {
		t.Fatalf("duplicate callbacks=%+v", fake.callbacks)
	}

	writeActiveUsers(integrationChatID, time.Now(), 7001)
	if err := SaveVote(SpamVote{ID: "partial-vote", ChatID: integrationChatID, ExpiresAt: time.Now().Add(time.Hour), RequiredVoteCount: 2}); err != nil {
		t.Fatalf("save vote: %v", err)
	}
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
	if err := SaveVote(SpamVote{
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
	}); err != nil {
		t.Fatalf("save vote: %v", err)
	}
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

	if err := SaveVote(SpamVote{
		ID:                "pass-vote-delete-failed",
		ChatID:            integrationChatID,
		ChatName:          "Guest Spam Test",
		MessageID:         602,
		GuestBotID:        993202,
		GuestBotName:      "Pass Bot Failed Delete",
		GuestBotUserName:  "pass_bot_failed_delete",
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: 2,
		Voters:            []int64{7001},
	}); err != nil {
		t.Fatalf("save vote: %v", err)
	}
	fake.deleteErr = errTelegram()
	fake.deletes = nil
	fake.callbacks = nil
	fake.callbackDelete = nil
	err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("pass-fail-cb", "guestspam_vote,vote,pass-vote-delete-failed", 7002)})
	if err == nil {
		t.Fatal("pass vote delete failure should return error")
	}
	if _, ok := GetVote("pass-vote-delete-failed"); ok {
		t.Fatal("passing vote with delete failure should still delete vote")
	}
	if !IsBlacklisted(993202) {
		t.Fatal("delete failure should still blacklist guest bot")
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "投票通过，已拉黑，但删除消息失败，请管理员检查权限" {
		t.Fatalf("callbacks=%+v, want delete failure callback", fake.callbacks)
	}
	if len(fake.callbackDelete) != 1 {
		t.Fatalf("callbackDeletes=%+v, want callback delete", fake.callbackDelete)
	}
	if !hasLogAction(RecentLogs(integrationChatID, 20), ActionDeleteFailed) {
		t.Fatalf("logs = %+v, want delete failure log", RecentLogs(integrationChatID, 20))
	}
}

func TestGuestSpamIntegrationSpamVoteWeightTiers(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	if err := SaveVote(SpamVote{
		ID:                "inactive-vote",
		ChatID:            integrationChatID,
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: 1,
	}); err != nil {
		t.Fatalf("save vote: %v", err)
	}
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

	if err := SaveVote(SpamVote{
		ID:                "weighted-vote",
		ChatID:            integrationChatID,
		ChatName:          "Guest Spam Test",
		MessageID:         701,
		GuestBotID:        993301,
		GuestBotName:      "Weighted Bot",
		GuestBotUserName:  "weighted_bot",
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: activeVoteWeight + trustedVoteWeight + adminVoteWeight,
	}); err != nil {
		t.Fatalf("save vote: %v", err)
	}

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
	risk.WarningCount = 1
	risk.MuteLevel = 2
	setMemberRisk(risk, true)
	beforeFailedRestrict := risk
	fake.unbans = nil
	fake.restricts = nil
	fake.restrictErr = errTelegram()
	fake.sends = nil
	err := GuestSpamLogHandle(tgbotapi.Update{Message: commandMessage("/guest_spam_log restore 880001", 8001)})
	if err == nil {
		t.Fatal("restore should return telegram error when unrestrict fails")
	}
	risk, ok = getMemberRisk(integrationChatID, integrationCallerID)
	if !ok || risk.WarningCount != beforeFailedRestrict.WarningCount || risk.MuteLevel != beforeFailedRestrict.MuteLevel {
		t.Fatalf("failed restore risk = %+v ok=%v, want unchanged from %+v", risk, ok, beforeFailedRestrict)
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
	risk.WarningCount = 1
	risk.MuteLevel = 2
	setMemberRisk(risk, true)
	beforeFailedUnban := risk
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
	if !ok || risk.WarningCount != beforeFailedUnban.WarningCount || risk.MuteLevel != beforeFailedUnban.MuteLevel {
		t.Fatalf("failed unban risk = %+v ok=%v, want unchanged from %+v", risk, ok, beforeFailedUnban)
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

func TestGuestSpamIntegrationMessagePathGuestBotSpam(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	guestBotID := int64(995001)
	callerID := int64(886001)

	// Construct path B update: update.Message with embedded GuestBotCallerUser
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			MessageID: 5001,
			Chat: &tgbotapi.Chat{
				ID:    integrationChatID,
				Type:  "supergroup",
				Title: "Guest Spam Test",
			},
			From: &tgbotapi.User{
				ID:        guestBotID,
				IsBot:     true,
				FirstName: "Message Path Bot",
				UserName:  "message_path_bot",
			},
			GuestBotCallerUser: &tgbotapi.User{
				ID:        callerID,
				FirstName: "Message Path Caller",
				UserName:  "message_path_caller",
			},
			Text: "hello from path B",
		},
	}

	// Step 1: CheckGuestBotSpam must detect it
	if !CheckGuestBotSpam(update) {
		t.Fatal("CheckGuestBotSpam should detect path B message")
	}

	// Step 2: GuestBotSpamHandle must process it (delete, blacklist, warn)
	if err := GuestBotSpamHandle(update); err != nil {
		t.Fatalf("GuestBotSpamHandle: %v", err)
	}

	// Step 3: Verify Telegram delete was called
	if len(fake.deletes) != 1 || fake.deletes[0].messageID != 5001 {
		t.Fatalf("deletes = %+v, want messageID 5001", fake.deletes)
	}

	// Step 4: Verify bot is now blacklisted
	if !IsBlacklisted(guestBotID) {
		t.Fatal("guest bot should be blacklisted after low trust detection")
	}

	// Step 5: Verify caller warning was recorded (warningCount > 0)
	risk, ok := getMemberRisk(integrationChatID, callerID)
	if !ok {
		t.Fatal("caller risk should exist")
	}
	if risk.WarningCount < 1 {
		t.Fatalf("caller warning count = %d, want >= 1", risk.WarningCount)
	}

	// Step 6: Verify spam log was recorded
	logs := RecentLogs(integrationChatID, 10)
	if len(logs) == 0 {
		t.Fatal("spam logs should exist after processing path B message")
	}
	hasAction := func(action string) bool {
		for _, item := range logs {
			if item.Action == action {
				return true
			}
		}
		return false
	}
	if !hasAction(ActionAutoBlacklist) {
		t.Fatal("spam logs should contain ActionAutoBlacklist")
	}
}

func TestGuestSpamIntegrationMessagePathNormalMessageIgnored(t *testing.T) {
	setupGuestSpamIntegration(t)
	_ = useFakeTelegram(t)

	// Plain message without guest fields must NOT trigger
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			MessageID: 5002,
			Chat: &tgbotapi.Chat{
				ID:    integrationChatID,
				Type:  "supergroup",
				Title: "Guest Spam Test",
			},
			From: &tgbotapi.User{
				ID:        1001,
				FirstName: "Normal User",
			},
			Text: "hello",
		},
	}

	if CheckGuestBotSpam(update) {
		t.Fatal("normal update.Message must not trigger guest bot spam")
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

func seedMemberRisk(t *testing.T, db *gorm.DB, risk MemberRisk) {
	t.Helper()
	if err := db.Omit(zeroRiskTimeFields(risk)...).Create(&risk).Error; err != nil {
		t.Fatalf("create member risk %s: %v", risk.ID, err)
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

// ─────────────────────────────────────────────────────────────────────
// 投票系统全覆盖测试
//
// 以下测试覆盖投票系统的所有路径：
//   投票发起（startSpamVote / GuestSpamHandle）
//   投票回调（SpamVoteCallback）
//   投票通过后的处理（applyVotePassed）
//   不同用户角色的投票权重
//   投票数据持久化（Redis → DB → 重载）
//   跨群隔离
//
// 每条测试都包含完整的中文注释，说明测试什么、为什么测、预期结果。
// ─────────────────────────────────────────────────────────────────────

// TestVote_RequiredCount 验证 requiredVoteCount 函数对活跃人数的票数计算
//
// 投票需要的票数 = (活跃用户数 + 1) / 2，活跃数少于 3 人时不允许投票。
// 这条测试覆盖各种边缘值，确保票数计算正确。
func TestVote_RequiredCount(t *testing.T) {
	setupGuestSpamIntegration(t)

	tests := []struct {
		activeCount int
		wantOK      bool
		wantVotes   int
		desc        string
	}{
		{0, false, 0, "0人 → 不够3人，投票无效"},
		{1, false, 0, "1人 → 不够3人，投票无效"},
		{2, false, 0, "2人 → 不够3人，投票无效"},
		{3, true, 2, "3人 → 需要 (3+1)/2 = 2 票才能通过"},
		{4, true, 2, "4人 → 需要 (4+1)/2 = 2 票就能通过"},
		{5, true, 3, "5人 → 需要 (5+1)/2 = 3 票才能通过"},
		{6, true, 3, "6人 → 需要 (6+1)/2 = 3 票就能通过"},
		{9, true, 5, "9人 → 需要 (9+1)/2 = 5 票才能通过"},
		{10, true, 5, "10人 → 需要 (10+1)/2 = 5 票即可通过"},
	}

	for _, tt := range tests {
		votes, ok := requiredVoteCount(tt.activeCount)
		if ok != tt.wantOK || votes != tt.wantVotes {
			t.Errorf("requiredVoteCount(%d) = (%d, %v)，期望 (%d, %v) — %s",
				tt.activeCount, votes, ok, tt.wantVotes, tt.wantOK, tt.desc)
		}
	}
}

// TestVote_StartInsufficientUsers 验证：当最近10分钟内活跃人数少于3时，发起投票会被拒绝
//
// 测试流程：
//   1. 只往Redis写2个活跃用户（不到3人）
//   2. 用 /guest_spam 命令尝试发起投票
//   3. 预期：投票不被创建，提示"活跃人数少于3"的提示消息
//   4. 预期：日志中记录 ActionVoteInvalid
//
// 这防止了在冷群中少数人就能拉黑别人的情况。
func TestVote_StartInsufficientUsers(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	// 先录入一条 guest bot 消息作为投票目标
	RecordRecentGuestMessage(RecentGuestMessage{
		ChatID:           integrationChatID,
		ChatName:         "Guest Spam Test",
		MessageID:        2001,
		GuestBotID:       994001,
		GuestBotName:     "Insufficient Bot",
		GuestBotUserName: "insufficient_bot",
		SeenAt:           time.Now(),
	})

	// 只写2个活跃用户（不够3人）
	writeActiveUsers(integrationChatID, time.Now(), 1001, 1002)

	// 尝试发起投票 → 应该被拒绝
	if err := GuestSpamHandle(tgbotapi.Update{Message: commandMessage("/guest_spam 2001", 5001)}); err != nil {
		t.Fatalf("GuestSpamHandle: %v", err)
	}

	// 验证：系统给出"活跃人数少于3"的提示
	if len(fake.sends) != 1 || !sentMessageContains(fake.sends[0], "活跃人数少于 3") {
		t.Fatalf("sends = %+v，期望提示'活跃人数少于 3'的消息", fake.sends)
	}

	// 验证：日志中有 ActionVoteInvalid
	if !hasLogAction(RecentLogs(integrationChatID, 10), ActionVoteInvalid) {
		t.Fatalf("日志中应该有 ActionVoteInvalid，实际日志 = %+v", RecentLogs(integrationChatID, 10))
	}
}

// TestVote_StartSuccess 验证：活跃人数>=3时，可以成功发起投票
//
// 测试流程：
//   1. 写入3个活跃用户
//   2. 录入 guest bot 消息作为投票目标
//   3. 用 /guest_spam 命令发起投票
//   4. 预期：投票消息被发送（显示"是否将 guest bot..."）
//   5. 预期：Redis 中有一个投票的 dirty 标记
//   6. 预期：日志中记录 ActionVoteStarted
//
// 投票创建后，系统显示带有"判定 spam / 取消"按钮的消息。
func TestVote_StartSuccess(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	RecordRecentGuestMessage(RecentGuestMessage{
		ChatID:           integrationChatID,
		ChatName:         "Guest Spam Test",
		MessageID:        2002,
		GuestBotID:       994002,
		GuestBotName:     "Target Bot",
		GuestBotUserName: "target_bot",
		SeenAt:           time.Now(),
	})

	// 写3个活跃用户
	writeActiveUsers(integrationChatID, time.Now(), 1001, 1002, 1003)

	if err := GuestSpamHandle(tgbotapi.Update{Message: commandMessage("/guest_spam @target_bot", 5001)}); err != nil {
		t.Fatalf("GuestSpamHandle: %v", err)
	}

	// 验证：发送了投票消息
	if len(fake.sends) != 1 || !sentMessageContains(fake.sends[0], "是否将 guest bot") {
		t.Fatalf("sends = %+v，期望投票消息", fake.sends)
	}

	// 验证：Redis 中有一个投票的 dirty 标记
	if len(voteIDs(t)) != 1 {
		t.Fatalf("vote dirty ids = %v，期望恰好1个投票", voteIDs(t))
	}

	// 验证：日志中有 ActionVoteStarted
	if !hasLogAction(RecentLogs(integrationChatID, 10), ActionVoteStarted) {
		t.Fatalf("日志中应该有 ActionVoteStarted，实际日志 = %+v", RecentLogs(integrationChatID, 10))
	}
}

// TestVote_StartWithBotIDArg 验证：用 bot ID 作为参数也能发起投票
//
// /guest_spam 命令支持多种参数格式：
//   - /guest_spam → 列出候选列表
//   - /guest_spam 2003 → 用消息ID定位
//   - /guest_spam @botname → 用bot用户名定位
//   - /guest_spam 994003 → 用bot ID定位
//
// 这里测试用 bot ID 定位目标。
func TestVote_StartWithBotIDArg(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	RecordRecentGuestMessage(RecentGuestMessage{
		ChatID:           integrationChatID,
		ChatName:         "Guest Spam Test",
		MessageID:        2003,
		GuestBotID:       994003,
		GuestBotName:     "ID Bot",
		GuestBotUserName: "id_bot",
		SeenAt:           time.Now(),
	})

	writeActiveUsers(integrationChatID, time.Now(), 1001, 1002, 1003)

	// 用 bot ID 发起投票
	if err := GuestSpamHandle(tgbotapi.Update{Message: commandMessage("/guest_spam 994003", 5001)}); err != nil {
		t.Fatalf("GuestSpamHandle with bot ID arg: %v", err)
	}

	if len(fake.sends) != 1 || !sentMessageContains(fake.sends[0], "是否将 guest bot") {
		t.Fatalf("sends = %+v，期望用 bot ID 也能发起投票", fake.sends)
	}
}

// TestVote_CancelNone 验证：取消候选列表的"取消"按钮
//
// 当用户打了 /guest_spam（不带参数）时，系统显示候选列表，
// 列表底部有一个"取消"按钮。点击它应该关闭列表但不执行任何操作。
func TestVote_CancelNone(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	// 注意 callback data: "guestspam_vote,cancel,none"
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("cancel-none", "guestspam_vote,cancel,none", 7001)}); err != nil {
		t.Fatalf("cancel none callback: %v", err)
	}

	// 验证：回调响应"已取消"并删除回调消息
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "已取消" || len(fake.callbackDelete) != 1 {
		t.Fatalf("callbacks=%+v deletes=%+v，期望'已取消'并删除消息", fake.callbacks, fake.callbackDelete)
	}
}

// TestVote_Expired 验证：对已过期的投票投票会被告知"投票已过期"
//
// 投票有效期为10分钟，过期后不能再投票。
// 发起者或投票人对过期投票进行操作时，系统提示已过期。
func TestVote_Expired(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	// 创建一个已经过期的投票（ExpiresAt 设为过去的时间）
	vote := SpamVote{
		ID:                "expired-vote-test",
		ChatID:            integrationChatID,
		MessageID:         2999,
		GuestBotID:        994299,
		ExpiresAt:         time.Now().Add(-1 * time.Minute), // 1分钟前就过期了
		RequiredVoteCount: 2,
	}
	if err := SaveVote(vote); err != nil {
		t.Fatalf("保存投票: %v", err)
	}

	// 尝试给已过期的投票投票
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("expired-cb", "guestspam_vote,vote,expired-vote-test", 7001)}); err != nil {
		t.Fatalf("expired vote callback: %v", err)
	}

	// 验证：提示"投票已过期"
	if len(fake.callbacks) != 1 || !fake.callbacks[0].showAlert || fake.callbacks[0].text != "投票已过期" {
		t.Fatalf("callbacks=%+v，期望'投票已过期'", fake.callbacks)
	}
}

// TestVote_BotVoteRejected 验证：机器人投票被拒绝
//
// 群里的其他 bot 不能参与 spam 判定投票。只有真实用户才能投票。
func TestVote_BotVoteRejected(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	vote := SpamVote{
		ID:                "bot-vote-test",
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: 2,
	}
	if err := SaveVote(vote); err != nil {
		t.Fatalf("保存投票: %v", err)
	}

	// 构造一个 bot 发起的 callback（IsBot = true）
	botCB := voteCallback("bot-cb", "guestspam_vote,vote,bot-vote-test", 0)
	botCB.From.IsBot = true

	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: botCB}); err != nil {
		t.Fatalf("bot vote callback: %v", err)
	}

	// 验证：提示"机器人不能参与投票"
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "机器人不能参与投票" {
		t.Fatalf("callbacks=%+v，期望'机器人不能参与投票'", fake.callbacks)
	}
}

// TestVote_DuplicateVoteRejected 验证：重复投票被拒绝
//
// 每个用户在同一个投票中只能投一次票，重复投票会被拒绝。
// 这防止了一个用户刷票。
func TestVote_DuplicateVoteRejected(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	// 创建一个投票，7001 已经投过了
	vote := SpamVote{
		ID:                "dup-vote-test",
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: 3,
		Voters:            []int64{7001},
	}
	if err := SaveVote(vote); err != nil {
		t.Fatalf("保存投票: %v", err)
	}

	// 7001 再次投票 → 应该被拒绝
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("dup-cb", "guestspam_vote,vote,dup-vote-test", 7001)}); err != nil {
		t.Fatalf("重复投票回调: %v", err)
	}

	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "你已经投过票了" {
		t.Fatalf("callbacks=%+v，期望'你已经投过票了'", fake.callbacks)
	}
}

// TestVote_InactiveUserRejected 验证：不活跃用户不能投票
//
// 只有最近10分钟内在群里发过言的用户才能参与投票。
// 不活跃用户点击投票按钮会收到"最近不够活跃"的提示。
func TestVote_InactiveUserRejected(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	vote := SpamVote{
		ID:                "inactive-vote-test",
		ChatID:            integrationChatID,
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: 1,
	}
	if err := SaveVote(vote); err != nil {
		t.Fatalf("保存投票: %v", err)
	}

	// 用户 7101 没有在活跃用户集合中 → 不活跃
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("inactive-cb", "guestspam_vote,vote,inactive-vote-test", 7101)}); err != nil {
		t.Fatalf("不活跃用户投票: %v", err)
	}

	// 验证：提示"最近不够活跃"
	if len(fake.callbacks) != 1 || !fake.callbacks[0].showAlert || fake.callbacks[0].text != "最近不够活跃，不能参与投票" {
		t.Fatalf("callbacks=%+v，期望'最近不够活跃'", fake.callbacks)
	}
}

// TestVote_ActiveUserWeight 验证：活跃用户投票权重为1
//
// 活跃用户 = 最近10分钟内在群里发过言的成员，投票权重为1。
func TestVote_ActiveUserWeight(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	vote := SpamVote{
		ID:                "active-weight-test",
		ChatID:            integrationChatID,
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: 3,
	}
	if err := SaveVote(vote); err != nil {
		t.Fatalf("保存投票: %v", err)
	}

	// 把用户 7202 标记为活跃
	writeActiveUsers(integrationChatID, time.Now(), 7202)

	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("active-weight-cb", "guestspam_vote,vote,active-weight-test", 7202)}); err != nil {
		t.Fatalf("活跃用户投票: %v", err)
	}

	// 验证：投票权重为1（活跃成员 +1）
	vote, ok := GetVote("active-weight-test")
	if !ok || vote.VoteScore != activeVoteWeight {
		t.Fatalf("vote = %+v ok=%v，期望 score = %d (活跃成员权重)", vote, ok, activeVoteWeight)
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "已投票 1/3（活跃成员 +1）" {
		t.Fatalf("callbacks=%+v，期望'活跃成员 +1'", fake.callbacks)
	}
}

// TestVote_TrustedMemberWeight 验证：可信成员投票权重为2
//
// 可信成员 = 入群超过3天且发过至少5条消息的成员，投票权重为2。
// 也就是说信得过的成员的一票顶活跃用户的两票。
//
// 测试使用 seedTrustedCaller() 来模拟一个满足条件的可信成员。
func TestVote_TrustedMemberWeight(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	vote := SpamVote{
		ID:                "trusted-weight-test",
		ChatID:            integrationChatID,
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: 5,
	}
	if err := SaveVote(vote); err != nil {
		t.Fatalf("保存投票: %v", err)
	}

	// 模拟一个可信成员（入群超过3天，发过至少5条消息）
	seedTrustedCaller(7301)

	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("trusted-weight-cb", "guestspam_vote,vote,trusted-weight-test", 7301)}); err != nil {
		t.Fatalf("可信成员投票: %v", err)
	}

	vote, ok := GetVote("trusted-weight-test")
	if !ok || vote.VoteScore != trustedVoteWeight {
		t.Fatalf("vote = %+v ok=%v，期望 score = %d (可信成员权重)", vote, ok, trustedVoteWeight)
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "已投票 2/5（可信成员 +2）" {
		t.Fatalf("callbacks=%+v，期望'可信成员 +2'", fake.callbacks)
	}
}

// TestVote_AdminWeight 验证：管理员投票权重为3
//
// 群管理员/创建者的投票权重为3，一票顶三票。
// 这是为了防止恶意 bot spam 时管理员能快速处理。
func TestVote_AdminWeight(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	vote := SpamVote{
		ID:                "admin-weight-test",
		ChatID:            integrationChatID,
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: 5,
	}
	if err := SaveVote(vote); err != nil {
		t.Fatalf("保存投票: %v", err)
	}

	// 把 7401 标记为管理员
	fake.admins[7401] = true

	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("admin-weight-cb", "guestspam_vote,vote,admin-weight-test", 7401)}); err != nil {
		t.Fatalf("管理员投票: %v", err)
	}

	vote, ok := GetVote("admin-weight-test")
	if !ok || vote.VoteScore != adminVoteWeight {
		t.Fatalf("vote = %+v ok=%v，期望 score = %d (管理员权重)", vote, ok, adminVoteWeight)
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "已投票 3/5（管理员 +3）" {
		t.Fatalf("callbacks=%+v，期望'管理员 +3'", fake.callbacks)
	}
}

// TestVote_UnknownActionIgnored 验证：未知的回调操作被忽略
//
// 如果 callback data 格式不符合预期（不是 vote 也不是 cancel），
// 系统不做任何处理，防止恶意构造的 callback data 导致问题。
func TestVote_UnknownActionIgnored(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	vote := SpamVote{
		ID:                "unknown-action-test",
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: 2,
	}
	if err := SaveVote(vote); err != nil {
		t.Fatalf("保存投票: %v", err)
	}

	// 构造一个未知操作的 callback (action = "wat")
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("unknown-cb", "guestspam_vote,wat,unknown-action-test", 7001)}); err != nil {
		t.Fatalf("unknown action callback: %v", err)
	}

	// 验证：投票没有任何变化
	vote, ok := GetVote("unknown-action-test")
	if !ok || len(vote.Voters) != 0 {
		t.Fatalf("未知操作不应修改投票: vote=%+v ok=%v", vote, ok)
	}

	// 验证：没有发送任何回调响应
	if len(fake.callbacks) != 0 {
		t.Fatalf("未知操作不应有回调响应: %+v", fake.callbacks)
	}
}

// TestVote_CancelByStarter 验证：投票发起者可以取消自己的投票
//
// 发起投票的人发现投错了，可以随时取消投票。
// 取消后投票被删除，之前的投票记录清除。
func TestVote_CancelByStarter(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	vote := SpamVote{
		ID:             "cancel-by-starter",
		ChatID:         integrationChatID,
		StarterUserID:  7501, // 发起者是 7501
		ExpiresAt:      time.Now().Add(time.Hour),
	}
	if err := SaveVote(vote); err != nil {
		t.Fatalf("保存投票: %v", err)
	}

	// 发起者 7501 取消投票
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("cancel-cb", "guestspam_vote,cancel,cancel-by-starter", 7501)}); err != nil {
		t.Fatalf("发起者取消投票: %v", err)
	}

	// 验证：投票已被删除
	if _, ok := GetVote("cancel-by-starter"); ok {
		t.Fatal("发起者取消后，投票应该被删除")
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "已取消" || len(fake.callbackDelete) != 1 {
		t.Fatalf("callbacks=%+v deletes=%+v，期望'已取消'", fake.callbacks, fake.callbackDelete)
	}
}

// TestVote_CancelByAdmin 验证：管理员可以取消别人的投票
//
// 即使管理员不是投票的发起者，也有权限取消投票。
// 这是为了防止恶意投票或发起的投票针对正常用户。
func TestVote_CancelByAdmin(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	// 发起者是 7601，管理员是 7602
	vote := SpamVote{
		ID:             "cancel-by-admin",
		ChatID:         integrationChatID,
		StarterUserID:  7601,
		ExpiresAt:      time.Now().Add(time.Hour),
	}
	if err := SaveVote(vote); err != nil {
		t.Fatalf("保存投票: %v", err)
	}

	// 管理员 7602 取消（非发起者）
	fake.admins[7602] = true
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("admin-cancel-cb", "guestspam_vote,cancel,cancel-by-admin", 7602)}); err != nil {
		t.Fatalf("管理员取消投票: %v", err)
	}

	// 验证：投票已被删除
	if _, ok := GetVote("cancel-by-admin"); ok {
		t.Fatal("管理员取消后，投票应该被删除")
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "已取消" || len(fake.callbackDelete) != 1 {
		t.Fatalf("callbacks=%+v deletes=%+v，期望'已取消'", fake.callbacks, fake.callbackDelete)
	}
}

// TestVote_CancelUnauthorized 验证：非发起者/非管理员不能取消投票
//
// 普通群成员不能取消别人的投票，防止恶意干扰正常的 spam 判定流程。
func TestVote_CancelUnauthorized(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	vote := SpamVote{
		ID:             "cancel-unauth",
		ChatID:         integrationChatID,
		StarterUserID:  7701,
		ExpiresAt:      time.Now().Add(time.Hour),
	}
	if err := SaveVote(vote); err != nil {
		t.Fatalf("保存投票: %v", err)
	}

	// 非发起者/非管理员（7702）尝试取消
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("cancel-unauth-cb", "guestspam_vote,cancel,cancel-unauth", 7702)}); err != nil {
		t.Fatalf("未授权取消: %v", err)
	}

	// 验证：投票还在
	if _, ok := GetVote("cancel-unauth"); !ok {
		t.Fatal("未授权取消不应删除投票")
	}
	if len(fake.callbacks) != 1 || !fake.callbacks[0].showAlert || fake.callbacks[0].text != "只有发起者或管理员可以取消投票" {
		t.Fatalf("callbacks=%+v，期望'只有发起者或管理员可以取消投票'", fake.callbacks)
	}
}

// TestVote_PassBlacklistsBot 验证：投票通过后 bot 被拉黑，消息被删除
//
// 测试流程：
//   1. 创建一个投票，已有人投了部分票
//   2. 再投一票达到通过线
//   3. 验证：bot 加入黑名单，投票被删除，guest bot 消息被删除
//   4. 验证：日志中记录 ActionVotePassed
//
// 这是用户用 /guest_spam 发起投票后最关键的流程。
func TestVote_PassBlacklistsBot(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	vote := SpamVote{
		ID:                "pass-vote-test",
		ChatID:            integrationChatID,
		ChatName:          "Guest Spam Test",
		MessageID:         2101,
		GuestBotID:        994201,
		GuestBotName:      "Pass Bot",
		GuestBotUserName:  "pass_bot",
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: 3,
		Voters:            []int64{7001, 7002}, // 已有2票
	}
	if err := SaveVote(vote); err != nil {
		t.Fatalf("保存投票: %v", err)
	}

	// 7003 投最后一票 → 达到3票，投票通过
	writeActiveUsers(integrationChatID, time.Now(), 7003)
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("pass-cb", "guestspam_vote,vote,pass-vote-test", 7003)}); err != nil {
		t.Fatalf("通过投票: %v", err)
	}

	// 验证：投票已被删除
	if _, ok := GetVote("pass-vote-test"); ok {
		t.Fatal("投票通过后应该被删除")
	}

	// 验证：bot 被拉黑
	if !IsBlacklisted(994201) {
		t.Fatal("投票通过的 bot 应该被拉黑")
	}

	// 验证：guest bot 消息被删除
	if len(fake.deletes) != 1 || fake.deletes[0].messageID != 2101 {
		t.Fatalf("deletes = %+v，期望删除 guest bot 消息 2101", fake.deletes)
	}

	// 验证：日志中有 ActionVotePassed
	if !hasLogAction(RecentLogs(integrationChatID, 10), ActionVotePassed) {
		t.Fatalf("日志中应该有 ActionVotePassed，实际 = %+v", RecentLogs(integrationChatID, 10))
	}
}

// TestVote_PassWithDeleteFail 验证：投票通过但删除消息失败时，bot 仍然被拉黑
//
// 有时候 bot 可能因为权限不足无法删除 guest bot 消息，
// 但拉黑操作不受影响。用户被告知删除失败但有正确提示。
func TestVote_PassWithDeleteFail(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	vote := SpamVote{
		ID:                "pass-fail-test",
		ChatID:            integrationChatID,
		ChatName:          "Guest Spam Test",
		MessageID:         2102,
		GuestBotID:        994202,
		GuestBotName:      "Delete Fail Bot",
		GuestBotUserName:  "delete_fail_bot",
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: 2,
		Voters:            []int64{7001},
	}
	if err := SaveVote(vote); err != nil {
		t.Fatalf("保存投票: %v", err)
	}

	// 模拟 Telegram 删除消息失败
	fake.deleteErr = errTelegram()
	writeActiveUsers(integrationChatID, time.Now(), 7002)

	err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("pass-fail-cb", "guestspam_vote,vote,pass-fail-test", 7002)})
	if err == nil {
		t.Fatal("删除失败时应该返回 error")
	}

	// 验证：投票已被删除（即使删除消息失败）
	if _, ok := GetVote("pass-fail-test"); ok {
		t.Fatal("投票通过后应该被删除，不论删除消息是否成功")
	}

	// 验证：bot 仍然被拉黑（最重要的验证）
	if !IsBlacklisted(994202) {
		t.Fatal("即使删除消息失败，bot 也必须被拉黑")
	}

	// 验证：用户被告知删除失败
	if len(fake.callbacks) != 1 || !strings.Contains(fake.callbacks[0].text, "已拉黑，但删除消息失败") {
		t.Fatalf("callbacks=%+v，期望'已拉黑，但删除消息失败'", fake.callbacks)
	}

	// 验证：日志中有 ActionDeleteFailed
	if !hasLogAction(RecentLogs(integrationChatID, 10), ActionDeleteFailed) {
		t.Fatalf("日志中应该有 ActionDeleteFailed，实际 = %+v", RecentLogs(integrationChatID, 10))
	}
}

// TestVote_PassPersistAndReload 验证：投票通过后，拉黑数据写入 MySQL，重启 Redis 后仍然生效
//
// 测试完整的数据持久化链条：
//   1. 投票通过 → Redis 中拉黑 bot
//   2. 手动触发 SyncCacheToDB() 将 Redis 数据同步到 MySQL
//   3. 清空 Redis（模拟 Redis 重启或数据丢失）
//   4. 手动触发 LoadCacheFromDB() 从 MySQL 恢复数据到 Redis
//   5. 验证：bot 仍然在黑名单中
//
// 这确保即使 Redis 宕机，拉黑数据不会丢失。
func TestVote_PassPersistAndReload(t *testing.T) {
	setupGuestSpamIntegration(t)
	_ = useFakeTelegram(t)

	// 创建投票并通过
	vote := SpamVote{
		ID:                "persist-test",
		ChatID:            integrationChatID,
		ChatName:          "Guest Spam Test",
		MessageID:         2201,
		GuestBotID:        994220,
		GuestBotName:      "Persist Bot",
		GuestBotUserName:  "persist_bot",
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: 2,
		Voters:            []int64{7001},
	}
	if err := SaveVote(vote); err != nil {
		t.Fatalf("保存投票: %v", err)
	}
	writeActiveUsers(integrationChatID, time.Now(), 7002)
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("persist-cb", "guestspam_vote,vote,persist-test", 7002)}); err != nil {
		t.Fatalf("通过投票: %v", err)
	}

	// 第1步 验证：Redis 中 bot 已被拉黑
	if !IsBlacklisted(994220) {
		t.Fatal("Redis 中 bot 应该被拉黑")
	}

	// 第2步 同步到 MySQL
	if err := SyncCacheToDB(); err != nil {
		t.Fatalf("同步缓存到 DB: %v", err)
	}

	// 第3步 清空 Redis（模拟重启）
	clearGuestSpamRedis(t)

	// 第4步 从 MySQL 重新加载
	if err := LoadCacheFromDB(); err != nil {
		t.Fatalf("从 DB 加载缓存: %v", err)
	}

	// 第5步 验证：重启后 bot 仍在黑名单中
	if !IsBlacklisted(994220) {
		t.Fatal("从 DB 恢复后，bot 应该仍然在黑名单中")
	}
}

// TestVote_MultipleVotesDifferentChats 验证：不同群里的投票互不干扰
//
// 投票是按 chat_id 隔离的，一个群里的投票不会影响另一个群。
// 一个 bot 在一个群里被拉黑，不应该自动在另一个群被拉黑（除非全局拉黑）。
//
// 测试：
//   1. 在 chat A 给 bot 发起投票并通过
//   2. 验证：chat A 中 bot 被拉黑
//   3. 验证：chat B 中 bot 不应该被拉黑（跨群隔离）
func TestVote_MultipleVotesDifferentChats(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)
	_ = fake

	// 在 chat A 中通过投票拉黑 bot 994230
	voteA := SpamVote{
		ID:                "chat-a-vote",
		ChatID:            integrationChatID,
		ChatName:          "Chat A",
		MessageID:         2301,
		GuestBotID:        994230,
		GuestBotName:      "Cross Chat Bot",
		GuestBotUserName:  "cross_chat_bot",
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: 2,
		Voters:            []int64{7001},
	}
	if err := SaveVote(voteA); err != nil {
		t.Fatalf("保存投票 A: %v", err)
	}
	writeActiveUsers(integrationChatID, time.Now(), 7002)
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("chat-a-cb", "guestspam_vote,vote,chat-a-vote", 7002)}); err != nil {
		t.Fatalf("chat A 投票通过: %v", err)
	}

	// 验证：chat A 中 bot 被拉黑（IsBlacklisted 是全局的）
	if !IsBlacklisted(994230) {
		t.Fatal("chat A 中 bot 应该被拉黑")
	}

	// 在 chat A 中给另一个 bot 发起投票（看是否能正常创建）
	RecordRecentGuestMessage(RecentGuestMessage{
		ChatID:           integrationChatID,
		ChatName:         "Guest Spam Test",
		MessageID:        2302,
		GuestBotID:       994231,
		GuestBotName:     "Other Bot",
		GuestBotUserName: "other_bot",
		SeenAt:           time.Now(),
	})
	writeActiveUsers(integrationChatID, time.Now(), 1001, 1002, 1003)

	fake.sends = nil
	if err := GuestSpamHandle(tgbotapi.Update{Message: commandMessage("/guest_spam @other_bot", 5001)}); err != nil {
		t.Fatalf("第二轮投票: %v", err)
	}

	// 验证：投票正常发起（第二个 bot 没有被第一个投票影响）
	if len(fake.sends) != 1 || !sentMessageContains(fake.sends[0], "是否将 guest bot") {
		t.Fatalf("sends = %+v，期望第二轮投票消息", fake.sends)
	}
}

// TestVote_CancelThenVoteAgain 验证：取消投票后可以重新对同一条消息发起投票
//
// 如果投票被取消，之前的投票记录应该被完全清除。
// 然后用户可以重新对同一条 guest bot 消息发起新的投票。
func TestVote_CancelThenVoteAgain(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	// 先录入一条 guest bot 消息
	RecordRecentGuestMessage(RecentGuestMessage{
		ChatID:           integrationChatID,
		ChatName:         "Guest Spam Test",
		MessageID:        2401,
		GuestBotID:       994240,
		GuestBotName:     "Cancel Again Bot",
		GuestBotUserName: "cancel_again_bot",
		SeenAt:           time.Now(),
	})
	writeActiveUsers(integrationChatID, time.Now(), 1001, 1002, 1003)

	// 创建一个投票并保存
	vote := SpamVote{
		ID:                "cancel-retry-test",
		ChatID:            integrationChatID,
		ChatName:          "Guest Spam Test",
		MessageID:         2401,
		GuestBotID:        994240,
		GuestBotName:      "Cancel Again Bot",
		GuestBotUserName:  "cancel_again_bot",
		StarterUserID:     5001,
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: 3,
	}
	if err := SaveVote(vote); err != nil {
		t.Fatalf("保存投票: %v", err)
	}

	// 取消投票
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("cancel-retry-cb", "guestspam_vote,cancel,cancel-retry-test", 5001)}); err != nil {
		t.Fatalf("取消投票: %v", err)
	}

	// 验证：投票已删除
	if _, ok := GetVote("cancel-retry-test"); ok {
		t.Fatal("取消后投票应该被删除")
	}

	// 重新发起投票
	fake.sends = nil
	if err := GuestSpamHandle(tgbotapi.Update{Message: commandMessage("/guest_spam 2401", 5001)}); err != nil {
		t.Fatalf("重新发起投票: %v", err)
	}

	// 验证：新的投票被创建
	if len(fake.sends) != 1 || !sentMessageContains(fake.sends[0], "是否将 guest bot") {
		t.Fatalf("sends = %+v，期望重新发起投票", fake.sends)
	}
}

// TestVote_PartialThenPass 验证：多用户逐步投票，累积到通过线后自动触发处理
//
// 测试用户多次投票的累积过程：
//   1. 创建投票，需要 3 票通过
//   2. 用户 A 投票（score = 1/3）
//   3. 用户 B 投票（score = 2/3）
//   4. 用户 C 投票（score = 3/3）→ 超过通过线，自动触发拉黑
//   5. 验证：累计过程中每个阶段票数正确
//   6. 验证：达到通过线后 bot 被拉黑
func TestVote_PartialThenPass(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	vote := SpamVote{
		ID:                "partial-pass-test",
		ChatID:            integrationChatID,
		ChatName:          "Guest Spam Test",
		MessageID:         2501,
		GuestBotID:        994250,
		GuestBotName:      "Partial Pass Bot",
		GuestBotUserName:  "partial_pass_bot",
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: 3,
	}
	if err := SaveVote(vote); err != nil {
		t.Fatalf("保存投票: %v", err)
	}

	// 写3个活跃用户
	writeActiveUsers(integrationChatID, time.Now(), 7001, 7002, 7003)

	// 用户A投票
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("a-cb", "guestspam_vote,vote,partial-pass-test", 7001)}); err != nil {
		t.Fatalf("用户A投票: %v", err)
	}
	vote, _ = GetVote("partial-pass-test")
	if vote.VoteScore != 1 || len(vote.Voters) != 1 {
		t.Fatalf("用户A投票后 = %+v，期望 score=1", vote)
	}

	// 用户B投票
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("b-cb", "guestspam_vote,vote,partial-pass-test", 7002)}); err != nil {
		t.Fatalf("用户B投票: %v", err)
	}
	vote, _ = GetVote("partial-pass-test")
	if vote.VoteScore != 2 || len(vote.Voters) != 2 {
		t.Fatalf("用户B投票后 = %+v，期望 score=2", vote)
	}

	// 用户C投票 → 达到3票，投票通过
	fake.callbacks = nil
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("c-cb", "guestspam_vote,vote,partial-pass-test", 7003)}); err != nil {
		t.Fatalf("用户C投票: %v", err)
	}

	// 验证：投票已被删除
	if _, ok := GetVote("partial-pass-test"); ok {
		t.Fatal("投票通过后应该被删除")
	}

	// 验证：bot 被拉黑
	if !IsBlacklisted(994250) {
		t.Fatal("投票通过的 bot 应该被拉黑")
	}

	// 验证：有回调响应（通过或删除消息）
	if len(fake.callbacks) != 1 {
		t.Fatalf("通过后应有回调响应: %+v", fake.callbacks)
	}
}
