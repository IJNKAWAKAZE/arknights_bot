package antispam

import (
	"fmt"
	"time"
)

const (
	redisPrefix          = "guest_spam"
	recentMessageLimit   = 10
	recentLogLimit       = 100
	recentMessageTTL     = 24 * time.Hour
	activityDayTTL       = 31 * 24 * time.Hour
	activeWindowTTL      = 11 * time.Minute
	voteTTL              = 11 * time.Minute
	riskTTL              = 60 * 24 * time.Hour
	logTTL               = 7 * 24 * time.Hour
	blacklistCacheTTL    = 180 * 24 * time.Hour
	warningWindow        = 30 * 24 * time.Hour
	trustAge             = 3 * 24 * time.Hour
	trustMessageCount    = int64(5)
	warningMuteThreshold = int64(2)
	logRetention         = 90 * 24 * time.Hour
)

var muteDurations = []time.Duration{
	24 * time.Hour,
	7 * 24 * time.Hour,
	30 * 24 * time.Hour,
}

func memberRiskKey(chatID, userID int64) string {
	return fmt.Sprintf("%s:risk:%d:%d", redisPrefix, chatID, userID)
}

func memberActivityDayKey(chatID, userID int64, day string) string {
	return fmt.Sprintf("%s:activity:%d:%d:%s", redisPrefix, chatID, userID, day)
}

func memberActivityDirtySetKey() string {
	return fmt.Sprintf("%s:activity:dirty", redisPrefix)
}

func memberDirtySetKey() string {
	return fmt.Sprintf("%s:risk:dirty", redisPrefix)
}

func blacklistSetKey() string {
	return fmt.Sprintf("%s:blacklist", redisPrefix)
}

func blacklistDirtySetKey() string {
	return fmt.Sprintf("%s:blacklist:dirty", redisPrefix)
}

func blacklistBotKey(botID int64) string {
	return fmt.Sprintf("%s:blacklist:bot:%d", redisPrefix, botID)
}

func recentMessagesKey(chatID int64) string {
	return fmt.Sprintf("%s:recent:%d", redisPrefix, chatID)
}

func activeUsersKey(chatID int64) string {
	return fmt.Sprintf("%s:active:%d", redisPrefix, chatID)
}

func logsKey(chatID int64) string {
	return fmt.Sprintf("%s:logs:%d", redisPrefix, chatID)
}

func logDirtyListKey() string {
	return fmt.Sprintf("%s:logs:dirty", redisPrefix)
}

func voteKey(voteID string) string {
	return fmt.Sprintf("%s:vote:%s", redisPrefix, voteID)
}

func voteDirtySetKey() string {
	return fmt.Sprintf("%s:vote:dirty", redisPrefix)
}
