package antispam

import "time"

const (
	ActionGuestSeen     = "guest_seen"
	ActionBlacklistHit  = "blacklist_hit"
	ActionAutoBlacklist = "auto_blacklist"
	ActionVoteStarted   = "vote_started"
	ActionVotePassed    = "vote_passed"
	ActionVoteInvalid   = "vote_invalid"
	ActionWarnCaller    = "warn_caller"
	ActionMuteCaller    = "mute_caller"
	ActionRestoreCaller = "restore_caller"
	ActionBanCaller     = "ban_caller"
	ActionBanCallerChat = "ban_caller_chat"
	ActionDeleteMessage = "delete_message"
	ActionDeleteFailed  = "delete_failed"
)

const (
	ReasonBlacklist       = "blacklist"
	ReasonLowTrust        = "low_trust"
	ReasonTrusted         = "trusted"
	ReasonVote            = "vote"
	ReasonAdminRestore    = "admin_restore"
	ReasonInsufficientAct = "insufficient_active_users"
)

type MemberRisk struct {
	ID                 string    `json:"id" gorm:"primaryKey;size:64"`
	ChatID             int64     `json:"chatId" gorm:"column:chat_id;index"`
	UserID             int64     `json:"userId" gorm:"column:user_id;index"`
	UserName           string    `json:"userName" gorm:"column:user_name;size:150"`
	FirstSeenAt        time.Time `json:"firstSeenAt" gorm:"column:first_seen_at"`
	LastMessageAt      time.Time `json:"lastMessageAt" gorm:"column:last_message_at"`
	RecentMessageCount int64     `json:"recentMessageCount" gorm:"column:recent_message_count"`
	WarningCount       int64     `json:"warningCount" gorm:"column:warning_count"`
	MuteLevel          int64     `json:"muteLevel" gorm:"column:mute_level"`
	LastPenaltyAt      time.Time `json:"lastPenaltyAt" gorm:"column:last_penalty_at"`
	CreateTime         time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"`
	UpdateTime         time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime"`
}

func (MemberRisk) TableName() string {
	return "guest_spam_member_risk"
}

type MemberActivity struct {
	ID           string    `json:"id" gorm:"primaryKey;size:80"`
	ChatID       int64     `json:"chatId" gorm:"column:chat_id;index"`
	UserID       int64     `json:"userId" gorm:"column:user_id;index"`
	Day          string    `json:"day" gorm:"column:day;size:8;index"`
	MessageCount int64     `json:"messageCount" gorm:"column:message_count"`
	CreateTime   time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"`
	UpdateTime   time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime"`
}

func (MemberActivity) TableName() string {
	return "guest_spam_member_activity"
}

type GuestBotBlacklist struct {
	ID                string    `json:"id" gorm:"primaryKey;size:64"`
	BotID             int64     `json:"botId" gorm:"column:bot_id;uniqueIndex"`
	BotName           string    `json:"botName" gorm:"column:bot_name;size:150"`
	BotUserName       string    `json:"botUserName" gorm:"column:bot_user_name;size:150"`
	Source            string    `json:"source" gorm:"column:source;size:50"`
	FirstChatID       int64     `json:"firstChatId" gorm:"column:first_chat_id"`
	FirstMessageID    int       `json:"firstMessageId" gorm:"column:first_message_id"`
	FirstCallerUserID int64     `json:"firstCallerUserId" gorm:"column:first_caller_user_id"`
	FirstCallerChatID int64     `json:"firstCallerChatId" gorm:"column:first_caller_chat_id"`
	CreateTime        time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"`
	UpdateTime        time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime"`
}

func (GuestBotBlacklist) TableName() string {
	return "guest_spam_bot_blacklist"
}

type SpamLog struct {
	ID             string    `json:"id" gorm:"primaryKey;size:64"`
	ChatID         int64     `json:"chatId" gorm:"column:chat_id;index"`
	ChatName       string    `json:"chatName" gorm:"column:chat_name;size:150"`
	MessageID      int       `json:"messageId" gorm:"column:message_id"`
	GuestBotID     int64     `json:"guestBotId" gorm:"column:guest_bot_id;index"`
	GuestBotName   string    `json:"guestBotName" gorm:"column:guest_bot_name;size:150"`
	GuestBotUser   string    `json:"guestBotUser" gorm:"column:guest_bot_user;size:150"`
	CallerUserID   int64     `json:"callerUserId" gorm:"column:caller_user_id;index"`
	CallerUserName string    `json:"callerUserName" gorm:"column:caller_user_name;size:150"`
	CallerChatID   int64     `json:"callerChatId" gorm:"column:caller_chat_id;index"`
	CallerChatName string    `json:"callerChatName" gorm:"column:caller_chat_name;size:150"`
	Action         string    `json:"action" gorm:"column:action;size:50;index"`
	Reason         string    `json:"reason" gorm:"column:reason;size:50"`
	Detail         string    `json:"detail" gorm:"column:detail;size:500"`
	CreateTime     time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"`
	UpdateTime     time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime"`
}

func (SpamLog) TableName() string {
	return "guest_spam_log"
}

type RecentGuestMessage struct {
	ChatID           int64     `json:"chatId"`
	ChatName         string    `json:"chatName"`
	MessageID        int       `json:"messageId"`
	GuestBotID       int64     `json:"guestBotId"`
	GuestBotName     string    `json:"guestBotName"`
	GuestBotUserName string    `json:"guestBotUserName"`
	CallerUserID     int64     `json:"callerUserId"`
	CallerUserName   string    `json:"callerUserName"`
	CallerChatID     int64     `json:"callerChatId"`
	CallerChatName   string    `json:"callerChatName"`
	SeenAt           time.Time `json:"seenAt"`
}

type SpamVote struct {
	ID                string    `json:"id"`
	ChatID            int64     `json:"chatId"`
	ChatName          string    `json:"chatName"`
	MessageID         int       `json:"messageId"`
	VoteMessageID     int       `json:"voteMessageId"`
	GuestBotID        int64     `json:"guestBotId"`
	GuestBotName      string    `json:"guestBotName"`
	GuestBotUserName  string    `json:"guestBotUserName"`
	StarterUserID     int64     `json:"starterUserId"`
	StarterUserName   string    `json:"starterUserName"`
	ActiveUserCount   int       `json:"activeUserCount"`
	RequiredVoteCount int       `json:"requiredVoteCount"`
	Voters            []int64   `json:"voters"`
	CreatedAt         time.Time `json:"createdAt"`
	ExpiresAt         time.Time `json:"expiresAt"`
}

type MemberTrust struct {
	Trusted            bool
	LowTrust           bool
	FirstSeenAt        time.Time
	LastMessageAt      time.Time
	RecentMessageCount int64
}

func isTrustedRiskAt(risk MemberRisk, now time.Time) bool {
	return !risk.FirstSeenAt.IsZero() && now.Sub(risk.FirstSeenAt) >= trustAge && risk.RecentMessageCount >= trustMessageCount
}
