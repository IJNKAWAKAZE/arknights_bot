package antispam

import (
	"fmt"

	bot "arknights_bot/config"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"log"
	"strings"
	"time"
)

func configGuestBotSpamEnabled() bool {
	return bot.GuestBotSpamEnabled
}

func CheckGuestBotSpam(update tgbotapi.Update) bool {
	if !configGuestBotSpamEnabled() {
		return false
	}
	if update.GuestMessage != nil && isGuestBotMessage(update.GuestMessage) {
		return true
	}
	if update.Message != nil && isGuestBotMessage(update.Message) {
		return true
	}
	return false
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
	if message == nil {
		message = update.Message
	}
	if message == nil || message.Chat == nil || message.From == nil {
		return nil
	}
	decision := EvaluateGuestMessage(message)
	ApplyGuestSpamState(decision)
	if decision.DeleteMessage {
		if err := deleteGuestMessageWithLog(message, decision.Reason); err != nil {
			log.Printf("guest spam: delete low trust message failed: %v", err)
		}
	}
	if decision.BanCallerUser || decision.BanCallerChat {
		penalizeBlacklistCaller(message)
	}
	if decision.MuteCaller {
		muteCaller(message, message.GuestBotCallerUser, decision.MuteDuration)
	}
	return nil
}

func ApplyGuestSpamState(decision GuestSpamDecision) {
	for _, item := range decision.Logs {
		AddLog(item)
	}
	if decision.BlacklistBot {
		AddBlacklist(decision.GuestBot, true)
	}
}

func EvaluateGuestMessage(message *tgbotapi.Message) GuestSpamDecision {
	decision := GuestSpamDecision{}
	if message == nil || message.Chat == nil || message.From == nil {
		return decision
	}
	RecordRecentGuestMessage(recentFromMessage(message))

	guestBotID := message.From.ID
	if IsBlacklisted(guestBotID) {
		decision.Reason = ReasonBlacklist
		decision.DeleteMessage = true
		decision.BanCallerUser = message.GuestBotCallerUser != nil && !message.GuestBotCallerUser.IsBot
		decision.BanCallerChat = message.GuestBotCallerChat != nil
		decision.Logs = append(decision.Logs, logFromMessage(message, ActionBlacklistHit, ReasonBlacklist, "global blacklist hit"))
		return decision
	}

	if caller := message.GuestBotCallerUser; caller != nil {
		trust := TrustFor(message.Chat.ID, caller.ID)
		decision.Trusted = trust.Trusted
		if trust.LowTrust {
			decision.Reason = ReasonLowTrust
			decision.DeleteMessage = true
			decision.BlacklistBot = true
			decision.WarnCaller = true
			decision.GuestBot = blacklistFromMessage(message, "low_trust")
			risk, shouldMute := AddWarning(message.Chat.ID, caller.ID, caller.FullName())
			decision.CallerRisk = risk
			decision.MuteCaller = shouldMute
			if shouldMute {
				decision.MuteDuration = MuteDuration(risk.MuteLevel)
			}
			decision.Logs = append(decision.Logs, logFromMessage(message, ActionAutoBlacklist, ReasonLowTrust, "blacklisted guest bot"))
			decision.Logs = append(decision.Logs, logFromMessage(message, ActionWarnCaller, ReasonLowTrust, fmt.Sprintf("warning count %d", risk.WarningCount)))
			return decision
		}
		decision.Reason = ReasonTrusted
		decision.Logs = append(decision.Logs, logFromMessage(message, ActionGuestSeen, ReasonTrusted, "trusted caller; message allowed"))
		return decision
	}

	decision.Reason = ReasonTrusted
	if message.GuestBotCallerChat != nil {
		decision.Logs = append(decision.Logs, logFromMessage(message, ActionGuestSeen, ReasonTrusted, "caller chat; message allowed"))
		return decision
	}
	decision.Logs = append(decision.Logs, logFromMessage(message, ActionGuestSeen, ReasonTrusted, "guest message without caller; message allowed"))
	return decision
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

func deleteGuestMessageWithLog(message *tgbotapi.Message, reason string) error {
	return deleteGuestMessageByLog(logFromMessage(message, "", reason, ""))
}

func deleteGuestMessageByLog(item SpamLog) error {
	if _, err := guestSpamTelegram.DeleteMessage(item.ChatID, item.MessageID); err != nil {
		item.Action = ActionDeleteFailed
		item.Detail = err.Error()
		AddLog(item)
		return err
	}
	item.Action = ActionDeleteMessage
	item.Detail = "deleted guest bot message"
	AddLog(item)
	return nil
}

func penalizeBlacklistCaller(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	if caller := message.GuestBotCallerUser; caller != nil && !caller.IsBot {
		if _, err := guestSpamTelegram.BanChatMember(chatID, caller.ID); err != nil {
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
		if _, err := guestSpamTelegram.Request(config); err != nil {
			AddLog(logFromMessage(message, ActionBanCallerChat, ReasonBlacklist, err.Error()))
			log.Printf("guest spam: ban caller chat %d failed: %v", callerChat.ID, err)
			return
		}
		AddLog(logFromMessage(message, ActionBanCallerChat, ReasonBlacklist, "banned caller chat of blacklisted guest bot"))
	}
}

func blacklistFromMessage(message *tgbotapi.Message, source string) GuestBotBlacklist {
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
	return item
}

func muteCaller(message *tgbotapi.Message, caller *tgbotapi.User, duration time.Duration) {
	if caller == nil || duration <= 0 {
		return
	}
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
	if _, err := guestSpamTelegram.Request(config); err != nil {
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
