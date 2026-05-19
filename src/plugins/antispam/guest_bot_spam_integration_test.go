//go:build integration

package antispam

import (
	bot "arknights_bot/config"
	"context"
	"fmt"
	"os"
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
