package antispam

import (
	bot "arknights_bot/config"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"log"
	"strings"
	"time"
)

func CheckGuestBotSpam(update tgbotapi.Update) bool {
	if !bot.GuestBotSpamEnabled {
		return false
	}
	message := update.GuestMessage
	return message != nil && isGuestBotMessage(message)
}

func CheckNormalActivity(update tgbotapi.Update) bool {
	return update.Message != nil && isTrackableMessage(update.Message)
}

func TrackActivityHandle(update tgbotapi.Update) error {
	message := update.Message
	if message == nil || message.From == nil || message.Chat == nil {
		return nil
	}
	RecordMessageActivity(message.Chat.ID, message.From.ID, message.From.FullName())
	return nil
}

func GuestBotSpamHandle(update tgbotapi.Update) error {
	message := update.GuestMessage
	if message == nil || message.Chat == nil || message.From == nil {
		return nil
	}
	recent := recentFromMessage(message)
	RecordRecentGuestMessage(recent)

	guestBotID := message.From.ID
	if IsBlacklisted(guestBotID) {
		AddLog(logFromMessage(message, ActionBlacklistHit, ReasonBlacklist, "global blacklist hit"))
		deleteGuestMessageWithLog(message, ReasonBlacklist)
		penalizeBlacklistCaller(message)
		return nil
	}

	if caller := message.GuestBotCallerUser; caller != nil {
		trust := TrustFor(message.Chat.ID, caller.ID)
		if trust.LowTrust {
			deleteGuestMessageWithLog(message, ReasonLowTrust)
			addGuestBotToBlacklist(message, "low_trust")
			warnCaller(message, caller)
			return nil
		}
		AddLog(logFromMessage(message, ActionGuestSeen, ReasonTrusted, "trusted caller; message allowed"))
		return nil
	}

	if message.GuestBotCallerChat != nil {
		AddLog(logFromMessage(message, ActionGuestSeen, ReasonTrusted, "caller chat; message allowed"))
		return nil
	}

	AddLog(logFromMessage(message, ActionGuestSeen, ReasonTrusted, "guest message without caller; message allowed"))
	return nil
}

func isGuestBotMessage(message *tgbotapi.Message) bool {
	return message.GuestBotCallerUser != nil || message.GuestBotCallerChat != nil || message.GuestQueryID != ""
}

func isTrackableMessage(message *tgbotapi.Message) bool {
	if message == nil || message.From == nil || message.From.IsBot || message.Chat == nil {
		return false
	}
	if message.Chat.IsPrivate() || message.IsCommand() {
		return false
	}
	if message.Text == "" && message.Caption == "" {
		return false
	}
	return true
}

func deleteGuestMessageWithLog(message *tgbotapi.Message, reason string) {
	if _, err := message.Delete(); err != nil {
		AddLog(logFromMessage(message, ActionDeleteFailed, reason, err.Error()))
		log.Printf("guest spam: delete message failed: %v", err)
		return
	}
	AddLog(logFromMessage(message, ActionDeleteMessage, reason, "deleted guest bot message"))
}

func penalizeBlacklistCaller(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	if caller := message.GuestBotCallerUser; caller != nil && !caller.IsBot {
		if _, err := bot.Arknights.BanChatMember(chatID, caller.ID); err != nil {
			AddLog(logFromMessage(message, ActionBanCaller, ReasonBlacklist, err.Error()))
			log.Printf("guest spam: ban caller %s failed: %v", userLogName(caller), err)
			return
		}
		AddLog(logFromMessage(message, ActionBanCaller, ReasonBlacklist, "banned caller of blacklisted guest bot"))
		return
	}
	if callerChat := message.GuestBotCallerChat; callerChat != nil {
		config := tgbotapi.BanChatSenderChatConfig{
			ChatID:       chatID,
			SenderChatID: callerChat.ID,
		}
		if _, err := bot.Arknights.Request(config); err != nil {
			AddLog(logFromMessage(message, ActionBanCallerChat, ReasonBlacklist, err.Error()))
			log.Printf("guest spam: ban caller chat %d failed: %v", callerChat.ID, err)
			return
		}
		AddLog(logFromMessage(message, ActionBanCallerChat, ReasonBlacklist, "banned caller chat of blacklisted guest bot"))
	}
}

func addGuestBotToBlacklist(message *tgbotapi.Message, source string) {
	if message.From == nil {
		return
	}
	item := GuestBotBlacklist{
		BotID:          message.From.ID,
		BotName:        message.From.FullName(),
		BotUserName:    message.From.UserName,
		Source:         source,
		FirstChatID:    message.Chat.ID,
		FirstMessageID: message.MessageID,
	}
	if message.GuestBotCallerUser != nil {
		item.FirstCallerUserID = message.GuestBotCallerUser.ID
	}
	if message.GuestBotCallerChat != nil {
		item.FirstCallerChatID = message.GuestBotCallerChat.ID
	}
	AddBlacklist(item, true)
	AddLog(logFromMessage(message, ActionAutoBlacklist, ReasonLowTrust, "blacklisted guest bot"))
}

func warnCaller(message *tgbotapi.Message, caller *tgbotapi.User) {
	risk, shouldMute := AddWarning(message.Chat.ID, caller.ID, caller.FullName())
	AddLog(logFromMessage(message, ActionWarnCaller, ReasonLowTrust, fmt.Sprintf("warning count %d", risk.WarningCount)))
	if !shouldMute {
		return
	}
	duration := MuteDuration(risk.MuteLevel)
	until := time.Now().Add(duration)
	config := tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: message.Chat.ID,
			UserID: caller.ID,
		},
		UntilDate: until.Unix(),
		Permissions: &tgbotapi.ChatPermissions{
			CanSendMessages: false,
		},
	}
	if _, err := bot.Arknights.Request(config); err != nil {
		AddLog(logFromMessage(message, ActionMuteCaller, ReasonLowTrust, err.Error()))
		log.Printf("guest spam: mute caller %s failed: %v", userLogName(caller), err)
		return
	}
	AddLog(logFromMessage(message, ActionMuteCaller, ReasonLowTrust, fmt.Sprintf("muted for %s", duration)))
}

func recentFromMessage(message *tgbotapi.Message) RecentGuestMessage {
	item := RecentGuestMessage{
		ChatID:           message.Chat.ID,
		ChatName:         message.Chat.Title,
		MessageID:        message.MessageID,
		GuestBotID:       message.From.ID,
		GuestBotName:     message.From.FullName(),
		GuestBotUserName: message.From.UserName,
		SeenAt:           time.Now(),
	}
	if caller := message.GuestBotCallerUser; caller != nil {
		item.CallerUserID = caller.ID
		item.CallerUserName = caller.FullName()
	}
	if callerChat := message.GuestBotCallerChat; callerChat != nil {
		item.CallerChatID = callerChat.ID
		item.CallerChatName = callerChat.Title
	}
	return item
}

func logFromMessage(message *tgbotapi.Message, action string, reason string, detail string) SpamLog {
	item := SpamLog{
		ChatID:         message.Chat.ID,
		ChatName:       message.Chat.Title,
		MessageID:      message.MessageID,
		GuestBotID:     message.From.ID,
		GuestBotName:   message.From.FullName(),
		GuestBotUser:   message.From.UserName,
		Action:         action,
		Reason:         reason,
		Detail:         detail,
		CallerUserID:   0,
		CallerUserName: "",
		CallerChatID:   0,
		CallerChatName: "",
	}
	if caller := message.GuestBotCallerUser; caller != nil {
		item.CallerUserID = caller.ID
		item.CallerUserName = caller.FullName()
	}
	if callerChat := message.GuestBotCallerChat; callerChat != nil {
		item.CallerChatID = callerChat.ID
		item.CallerChatName = callerChat.Title
	}
	return item
}

func userLogName(user *tgbotapi.User) string {
	if user == nil {
		return "<nil>"
	}
	if user.UserName != "" {
		return fmt.Sprintf("@%s(%d)", user.UserName, user.ID)
	}
	return fmt.Sprintf("%s(%d)", user.FullName(), user.ID)
}

func findRecentMessage(chatID int64, selector string) (RecentGuestMessage, bool) {
	recents := RecentGuestMessages(chatID)
	if len(recents) == 0 {
		return RecentGuestMessage{}, false
	}
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return RecentGuestMessage{}, false
	}
	if selector == "cancel" || selector == "取消" {
		return RecentGuestMessage{}, false
	}
	for _, item := range recents {
		if fmt.Sprint(item.MessageID) == selector || fmt.Sprint(item.GuestBotID) == selector || item.GuestBotUserName == strings.TrimPrefix(selector, "@") {
			return item, true
		}
	}
	return RecentGuestMessage{}, false
}
