package antispam

import (
	bot "arknights_bot/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm/clause"
)

var redisCtx = context.Background()

func Init() {
	if err := LoadCacheFromDB(); err != nil {
		log.Printf("guest spam: loading cache from db failed: %v", err)
	}
}

func LoadCacheFromDB() error {
	if bot.DBEngine == nil || bot.GoRedis == nil {
		return nil
	}

	var blacklists []GuestBotBlacklist
	if err := bot.DBEngine.Find(&blacklists).Error; err != nil {
		return err
	}
	for _, item := range blacklists {
		cacheBlacklist(item, false)
	}

	var risks []MemberRisk
	if err := bot.DBEngine.Find(&risks).Error; err != nil {
		return err
	}
	for _, risk := range risks {
		setJSON(memberRiskKey(risk.ChatID, risk.UserID), risk, riskTTL)
		if !risk.LastMessageAt.IsZero() && time.Since(risk.LastMessageAt) <= activeWindowTTL {
			bot.GoRedis.SAdd(redisCtx, activeUsersKey(risk.ChatID), risk.UserID)
			bot.GoRedis.Expire(redisCtx, activeUsersKey(risk.ChatID), activeWindowTTL)
		}
	}

	var activities []MemberActivity
	cutoffDay := startOfDay(time.Now().AddDate(0, 0, -29))
	if err := bot.DBEngine.Where("activity_day >= ?", cutoffDay).Find(&activities).Error; err != nil {
		return err
	}
	for _, activity := range activities {
		key := memberActivityDayKey(activity.ChatID, activity.UserID, dayKey(activity.ActivityDay))
		bot.GoRedis.Set(redisCtx, key, activity.MessageCount, activityDayTTL)
	}

	var logs []SpamLog
	cutoff := time.Now().Add(-logRetention)
	if err := bot.DBEngine.Where("create_time >= ?", cutoff).Order("create_time desc").Limit(1000).Find(&logs).Error; err != nil {
		return err
	}
	for i := len(logs) - 1; i >= 0; i-- {
		cacheLog(logs[i], false)
	}
	recentSeen := make(map[string]struct{})
	recentItems := make([]RecentGuestMessage, 0, len(logs))
	for _, item := range logs {
		if item.MessageID == 0 || item.GuestBotID == 0 {
			continue
		}
		key := fmt.Sprintf("%d:%d", item.ChatID, item.MessageID)
		if _, ok := recentSeen[key]; ok {
			continue
		}
		recentSeen[key] = struct{}{}
		recentItems = append(recentItems, RecentGuestMessage{
			ChatID:           item.ChatID,
			ChatName:         item.ChatName,
			MessageID:        item.MessageID,
			GuestBotID:       item.GuestBotID,
			GuestBotName:     item.GuestBotName,
			GuestBotUserName: item.GuestBotUser,
			CallerUserID:     item.CallerUserID,
			CallerUserName:   item.CallerUserName,
			CallerChatID:     item.CallerChatID,
			CallerChatName:   item.CallerChatName,
			SeenAt:           item.CreateTime,
		})
	}
	for i := len(recentItems) - 1; i >= 0; i-- {
		cacheRecentGuestMessage(recentItems[i])
	}
	return nil
}

func SyncCacheToDB() error {
	if bot.DBEngine == nil || bot.GoRedis == nil {
		return nil
	}
	if err := syncRisksToDB(); err != nil {
		return err
	}
	if err := syncActivitiesToDB(); err != nil {
		return err
	}
	if err := syncBlacklistsToDB(); err != nil {
		return err
	}
	if err := syncLogsToDB(); err != nil {
		return err
	}
	if err := cleanupOldLogs(); err != nil {
		return err
	}
	return nil
}

func syncRisksToDB() error {
	ids, err := bot.GoRedis.SMembers(redisCtx, memberDirtySetKey()).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	for _, id := range ids {
		parts := splitDirtyID(id)
		if len(parts) != 2 {
			continue
		}
		chatID, _ := strconv.ParseInt(parts[0], 10, 64)
		userID, _ := strconv.ParseInt(parts[1], 10, 64)
		risk, ok := getMemberRisk(chatID, userID)
		if !ok {
			continue
		}
		db := bot.DBEngine.Omit(zeroRiskTimeFields(risk)...)
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(memberRiskUpdates(risk)),
		}).Create(&risk).Error; err != nil {
			return err
		}
		bot.GoRedis.SRem(redisCtx, memberDirtySetKey(), id)
	}
	return nil
}

func syncActivitiesToDB() error {
	ids, err := bot.GoRedis.SMembers(redisCtx, memberActivityDirtySetKey()).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	for _, id := range ids {
		parts := splitDirtyID(id)
		if len(parts) != 3 {
			continue
		}
		chatID, _ := strconv.ParseInt(parts[0], 10, 64)
		userID, _ := strconv.ParseInt(parts[1], 10, 64)
		day := parts[2]
		activityDay, err := parseDayKey(day)
		if err != nil {
			log.Printf("guest spam: parse activity day failed: %v", err)
			bot.GoRedis.SRem(redisCtx, memberActivityDirtySetKey(), id)
			continue
		}
		count, err := bot.GoRedis.Get(redisCtx, memberActivityDayKey(chatID, userID, day)).Int64()
		if errors.Is(err, redis.Nil) {
			bot.GoRedis.SRem(redisCtx, memberActivityDirtySetKey(), id)
			continue
		}
		if err != nil {
			return err
		}
		activity := MemberActivity{
			ID:           activityID(chatID, userID, day),
			ChatID:       chatID,
			UserID:       userID,
			ActivityDay:  activityDay,
			MessageCount: count,
		}
		if err := bot.DBEngine.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&activity).Error; err != nil {
			return err
		}
		bot.GoRedis.SRem(redisCtx, memberActivityDirtySetKey(), id)
	}
	return nil
}

func zeroRiskTimeFields(risk MemberRisk) []string {
	fields := make([]string, 0, 3)
	if risk.FirstSeenAt.IsZero() {
		fields = append(fields, "FirstSeenAt")
	}
	if risk.LastMessageAt.IsZero() {
		fields = append(fields, "LastMessageAt")
	}
	if risk.LastPenaltyAt.IsZero() {
		fields = append(fields, "LastPenaltyAt")
	}
	return fields
}

func memberRiskUpdates(risk MemberRisk) map[string]any {
	updates := map[string]any{
		"chat_id":              risk.ChatID,
		"user_id":              risk.UserID,
		"user_name":            risk.UserName,
		"recent_message_count": risk.RecentMessageCount,
		"warning_count":        risk.WarningCount,
		"mute_level":           risk.MuteLevel,
		"update_time":          time.Now(),
	}
	if !risk.FirstSeenAt.IsZero() {
		updates["first_seen_at"] = risk.FirstSeenAt
	}
	if !risk.LastMessageAt.IsZero() {
		updates["last_message_at"] = risk.LastMessageAt
	}
	if !risk.LastPenaltyAt.IsZero() {
		updates["last_penalty_at"] = risk.LastPenaltyAt
	}
	return updates
}

func syncBlacklistsToDB() error {
	ids, err := bot.GoRedis.SMembers(redisCtx, blacklistDirtySetKey()).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	for _, id := range ids {
		botID, _ := strconv.ParseInt(id, 10, 64)
		item, ok := getBlacklistItem(botID)
		if !ok {
			continue
		}
		if err := bot.DBEngine.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "bot_id"}},
			UpdateAll: true,
		}).Create(&item).Error; err != nil {
			return err
		}
		bot.GoRedis.SRem(redisCtx, blacklistDirtySetKey(), id)
	}
	return nil
}

func syncLogsToDB() error {
	raws, err := bot.GoRedis.LRange(redisCtx, logDirtyListKey(), 0, -1).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	for _, raw := range raws {
		var item SpamLog
		if err := json.Unmarshal([]byte(raw), &item); err != nil {
			log.Printf("guest spam: unmarshal dirty log failed: %v", err)
			bot.GoRedis.LRem(redisCtx, logDirtyListKey(), 1, raw)
			continue
		}
		if err := bot.DBEngine.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).Create(&item).Error; err != nil {
			return err
		}
		bot.GoRedis.LRem(redisCtx, logDirtyListKey(), 1, raw)
	}
	return nil
}

func cleanupOldLogs() error {
	return bot.DBEngine.Where("create_time < ?", time.Now().Add(-logRetention)).Delete(&SpamLog{}).Error
}

func RecordMessageActivity(chatID int64, userID int64, userName string) {
	if chatID == 0 || userID == 0 || bot.GoRedis == nil {
		return
	}
	now := time.Now()
	risk, ok := getMemberRisk(chatID, userID)
	if !ok || risk.ID == "" {
		risk = MemberRisk{
			ID:          riskID(chatID, userID),
			ChatID:      chatID,
			UserID:      userID,
			FirstSeenAt: now,
		}
	}
	if risk.FirstSeenAt.IsZero() {
		risk.FirstSeenAt = now
	}
	risk.UserName = userName
	risk.LastMessageAt = now
	currentDay := dayKey(now)
	activityKey := memberActivityDayKey(chatID, userID, currentDay)
	bot.GoRedis.Incr(redisCtx, activityKey)
	bot.GoRedis.Expire(redisCtx, activityKey, activityDayTTL)
	bot.GoRedis.SAdd(redisCtx, memberActivityDirtySetKey(), fmt.Sprintf("%d:%d:%s", chatID, userID, currentDay))
	risk.RecentMessageCount = recentActivityCount(chatID, userID, now)
	setMemberRisk(risk, true)
	bot.GoRedis.SAdd(redisCtx, activeUsersKey(chatID), userID)
	bot.GoRedis.Expire(redisCtx, activeUsersKey(chatID), activeWindowTTL)
}

func TrustFor(chatID, userID int64) MemberTrust {
	risk, ok := getMemberRisk(chatID, userID)
	if !ok {
		now := time.Now()
		risk = MemberRisk{
			ID:          riskID(chatID, userID),
			ChatID:      chatID,
			UserID:      userID,
			FirstSeenAt: now,
		}
		setMemberRisk(risk, true)
	}
	recentCount := risk.RecentMessageCount
	if bot.GoRedis != nil {
		recentCount = recentActivityCount(chatID, userID, time.Now())
	}
	risk.RecentMessageCount = recentCount
	trusted := isTrustedRisk(risk)
	return MemberTrust{
		Trusted:            trusted,
		LowTrust:           !trusted,
		FirstSeenAt:        risk.FirstSeenAt,
		LastMessageAt:      risk.LastMessageAt,
		RecentMessageCount: recentCount,
	}
}

func AddWarning(chatID, userID int64, userName string) (MemberRisk, bool) {
	now := time.Now()
	risk, ok := getMemberRisk(chatID, userID)
	if !ok || risk.ID == "" {
		risk = MemberRisk{
			ID:          riskID(chatID, userID),
			ChatID:      chatID,
			UserID:      userID,
			FirstSeenAt: now,
		}
	}
	if risk.LastPenaltyAt.IsZero() || now.Sub(risk.LastPenaltyAt) > warningWindow {
		risk.WarningCount = 0
	}
	risk.UserName = userName
	risk.WarningCount++
	risk.LastPenaltyAt = now
	shouldMute := risk.WarningCount >= warningMuteThreshold
	if shouldMute {
		risk.MuteLevel++
		risk.WarningCount = 0
	}
	setMemberRisk(risk, true)
	return risk, shouldMute
}

func ClearWarnings(chatID, userID int64) {
	risk, ok := getMemberRisk(chatID, userID)
	if !ok {
		return
	}
	risk.WarningCount = 0
	risk.MuteLevel = 0
	setMemberRisk(risk, true)
}

func MuteDuration(level int64) time.Duration {
	if level <= 0 {
		return muteDurations[0]
	}
	index := int(level) - 1
	if index >= len(muteDurations) {
		index = len(muteDurations) - 1
	}
	return muteDurations[index]
}

func IsBlacklisted(botID int64) bool {
	if botID == 0 || bot.GoRedis == nil {
		return false
	}
	ok, err := bot.GoRedis.SIsMember(redisCtx, blacklistSetKey(), botID).Result()
	if err == nil && ok {
		return true
	}
	if bot.DBEngine == nil {
		return false
	}
	var count int64
	if err := bot.DBEngine.Model(&GuestBotBlacklist{}).Where("bot_id = ?", botID).Count(&count).Error; err != nil {
		log.Printf("guest spam: lookup blacklist db failed: %v", err)
		return false
	}
	if count > 0 {
		bot.GoRedis.SAdd(redisCtx, blacklistSetKey(), botID)
		return true
	}
	return false
}

func AddBlacklist(item GuestBotBlacklist, dirty bool) {
	if item.ID == "" {
		item.ID = blacklistID(item.BotID)
	}
	if item.CreateTime.IsZero() {
		item.CreateTime = time.Now()
	}
	item.UpdateTime = time.Now()
	cacheBlacklist(item, dirty)
}

func RecordRecentGuestMessage(message RecentGuestMessage) {
	if bot.GoRedis == nil {
		return
	}
	cacheRecentGuestMessage(message)
}

func RecentGuestMessages(chatID int64) []RecentGuestMessage {
	if bot.GoRedis == nil {
		return nil
	}
	raws, err := bot.GoRedis.LRange(redisCtx, recentMessagesKey(chatID), 0, recentMessageLimit-1).Result()
	if err != nil && err != redis.Nil {
		log.Printf("guest spam: read recent messages failed: %v", err)
		return nil
	}
	result := make([]RecentGuestMessage, 0, len(raws))
	for _, raw := range raws {
		var item RecentGuestMessage
		if err := json.Unmarshal([]byte(raw), &item); err == nil {
			result = append(result, item)
		}
	}
	return result
}

func ActiveUserCount(chatID int64) int {
	if bot.GoRedis == nil {
		return 0
	}
	count, err := bot.GoRedis.SCard(redisCtx, activeUsersKey(chatID)).Result()
	if err != nil {
		return 0
	}
	return int(count)
}

func SaveVote(vote SpamVote) {
	if bot.GoRedis == nil {
		return
	}
	if vote.ID == "" {
		vote.ID, _ = gonanoid.New(16)
	}
	setJSON(voteKey(vote.ID), vote, voteTTL)
	bot.GoRedis.SAdd(redisCtx, voteDirtySetKey(), vote.ID)
}

func GetVote(voteID string) (SpamVote, bool) {
	var vote SpamVote
	if bot.GoRedis == nil {
		return vote, false
	}
	if !getJSON(voteKey(voteID), &vote) {
		return vote, false
	}
	if time.Now().After(vote.ExpiresAt) {
		bot.GoRedis.Del(redisCtx, voteKey(voteID))
		return vote, false
	}
	return vote, true
}

func DeleteVote(voteID string) {
	if bot.GoRedis == nil {
		return
	}
	bot.GoRedis.Del(redisCtx, voteKey(voteID))
	bot.GoRedis.SRem(redisCtx, voteDirtySetKey(), voteID)
}

func AddLog(item SpamLog) {
	if bot.GoRedis == nil && bot.DBEngine == nil {
		return
	}
	if item.ID == "" {
		item.ID, _ = gonanoid.New(32)
	}
	if item.CreateTime.IsZero() {
		item.CreateTime = time.Now()
	}
	item.UpdateTime = time.Now()
	cacheLog(item, true)
}

func RecentLogs(chatID int64, limit int64) []SpamLog {
	if limit <= 0 {
		limit = 10
	}
	logs := make([]SpamLog, 0, limit)
	if bot.GoRedis != nil {
		raws, err := bot.GoRedis.LRange(redisCtx, logsKey(chatID), 0, limit-1).Result()
		if err != nil && err != redis.Nil {
			log.Printf("guest spam: read logs failed: %v", err)
			return nil
		}
		for _, raw := range raws {
			var item SpamLog
			if err := json.Unmarshal([]byte(raw), &item); err == nil {
				logs = append(logs, item)
			}
		}
	}
	if len(logs) == 0 && bot.DBEngine != nil {
		bot.DBEngine.Where("chat_id = ?", chatID).Order("create_time desc").Limit(int(limit)).Find(&logs)
		for i := len(logs) - 1; i >= 0; i-- {
			cacheLog(logs[i], false)
		}
	}
	return logs
}

func cacheBlacklist(item GuestBotBlacklist, dirty bool) {
	if bot.GoRedis == nil || item.BotID == 0 {
		return
	}
	setJSON(blacklistBotKey(item.BotID), item, blacklistCacheTTL)
	bot.GoRedis.SAdd(redisCtx, blacklistSetKey(), item.BotID)
	if dirty {
		bot.GoRedis.SAdd(redisCtx, blacklistDirtySetKey(), item.BotID)
	}
}

func getBlacklistItem(botID int64) (GuestBotBlacklist, bool) {
	var item GuestBotBlacklist
	if getJSON(blacklistBotKey(botID), &item) {
		return item, true
	}
	return item, false
}

func cacheLog(item SpamLog, dirty bool) {
	if bot.GoRedis == nil {
		if dirty && bot.DBEngine != nil {
			if err := bot.DBEngine.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoNothing: true,
			}).Create(&item).Error; err != nil {
				log.Printf("guest spam: write log db failed: %v", err)
			}
		}
		return
	}
	raw := marshalString(item)
	key := logsKey(item.ChatID)
	bot.GoRedis.LPush(redisCtx, key, raw)
	bot.GoRedis.LTrim(redisCtx, key, 0, recentLogLimit-1)
	bot.GoRedis.Expire(redisCtx, key, logTTL)
	if dirty {
		bot.GoRedis.RPush(redisCtx, logDirtyListKey(), raw)
	}
}

func cacheRecentGuestMessage(message RecentGuestMessage) {
	if bot.GoRedis == nil {
		return
	}
	if message.SeenAt.IsZero() {
		message.SeenAt = time.Now()
	}
	mustJSON := marshalString(message)
	key := recentMessagesKey(message.ChatID)
	bot.GoRedis.LPush(redisCtx, key, mustJSON)
	bot.GoRedis.LTrim(redisCtx, key, 0, recentMessageLimit-1)
	bot.GoRedis.Expire(redisCtx, key, recentMessageTTL)
}

func getMemberRisk(chatID, userID int64) (MemberRisk, bool) {
	var risk MemberRisk
	if getJSON(memberRiskKey(chatID, userID), &risk) {
		return risk, true
	}
	if bot.DBEngine == nil {
		return risk, false
	}
	if err := bot.DBEngine.Where("chat_id = ? and user_id = ?", chatID, userID).First(&risk).Error; err != nil {
		return risk, false
	}
	setJSON(memberRiskKey(chatID, userID), risk, riskTTL)
	return risk, true
}

func setMemberRisk(risk MemberRisk, dirty bool) {
	if bot.GoRedis == nil {
		return
	}
	setJSON(memberRiskKey(risk.ChatID, risk.UserID), risk, riskTTL)
	if dirty {
		bot.GoRedis.SAdd(redisCtx, memberDirtySetKey(), fmt.Sprintf("%d:%d", risk.ChatID, risk.UserID))
	}
}

func recentActivityCount(chatID, userID int64, now time.Time) int64 {
	if bot.GoRedis == nil {
		return 0
	}
	var total int64
	for i := 0; i < 30; i++ {
		day := dayKey(now.AddDate(0, 0, -i))
		count, err := bot.GoRedis.Get(redisCtx, memberActivityDayKey(chatID, userID, day)).Int64()
		if err != nil && err != redis.Nil {
			log.Printf("guest spam: read activity bucket failed: %v", err)
		}
		total += count
	}
	return total
}

func isTrustedRisk(risk MemberRisk) bool {
	return isTrustedRiskAt(risk, time.Now())
}

func setJSON(key string, val any, expiration time.Duration) {
	if bot.GoRedis == nil {
		return
	}
	bot.GoRedis.Set(redisCtx, key, marshalString(val), expiration)
}

func getJSON(key string, val any) bool {
	if bot.GoRedis == nil {
		return false
	}
	raw, err := bot.GoRedis.Get(redisCtx, key).Result()
	if err != nil {
		return false
	}
	return json.Unmarshal([]byte(raw), val) == nil
}

func marshalString(val any) string {
	data, err := json.Marshal(val)
	if err != nil {
		log.Printf("guest spam: marshal failed: %v", err)
		return "{}"
	}
	return string(data)
}

func riskID(chatID, userID int64) string {
	return fmt.Sprintf("%d:%d", chatID, userID)
}

func blacklistID(botID int64) string {
	return fmt.Sprintf("%d", botID)
}

func activityID(chatID, userID int64, day string) string {
	return fmt.Sprintf("%d:%d:%s", chatID, userID, day)
}

func dayKey(t time.Time) string {
	return startOfDay(t).Format("20060102")
}

func parseDayKey(day string) (time.Time, error) {
	parsed, err := time.ParseInLocation("20060102", day, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return startOfDay(parsed), nil
}

func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func splitDirtyID(id string) []string {
	return strings.Split(id, ":")
}
