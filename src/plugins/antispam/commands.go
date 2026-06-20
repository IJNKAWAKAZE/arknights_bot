package antispam

import (
	bot "arknights_bot/config"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"strconv"
	"strings"
	"time"
)

const (
	activeVoteWeight  = 1
	trustedVoteWeight = 2
	adminVoteWeight   = 3
)

func GuestSpamHandle(update tgbotapi.Update) error {
	message := update.Message
	if message == nil || message.Chat == nil || message.From == nil {
		return nil
	}
	chatID := message.Chat.ID
	guestSpamTelegram.QueueDelete(chatID, message.MessageID, 5)

	args := strings.TrimSpace(message.CommandArguments())
	if args == "" {
		return sendRecentGuestCandidates(message)
	}
	selected, ok := findRecentMessage(chatID, args)
	if !ok {
		return sendRecentGuestCandidates(message)
	}
	_, err := startSpamVote(message, selected)
	return err
}

func GuestSpamLogHandle(update tgbotapi.Update) error {
	message := update.Message
	if message == nil || message.Chat == nil || message.From == nil {
		return nil
	}
	chatID := message.Chat.ID
	userID := message.From.ID
	guestSpamTelegram.QueueDelete(chatID, message.MessageID, 5)
	if !guestSpamTelegram.IsAdmin(chatID, userID) {
		send := tgbotapi.NewMessage(chatID, "无使用权限！")
		send.ReplyToMessageID = message.MessageID
		msg, err := guestSpamTelegram.Send(send)
		if err == nil {
			guestSpamTelegram.QueueDelete(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		}
		return err
	}
	args := strings.Fields(message.CommandArguments())
	if len(args) >= 2 && args[0] == "restore" {
		target, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return sendTempMessage(chatID, message.MessageID, "用户 ID 格式错误")
		}
		return restoreCaller(chatID, target, message)
	}
	return sendLogs(chatID, message.MessageID)
}

func SpamVoteCallback(update tgbotapi.Update) error {
	callback := update.CallbackQuery
	if callback == nil {
		return nil
	}
	data := strings.Split(callback.Data, ",")
	if len(data) < 3 {
		return nil
	}
	action := data[1]
	voteID := data[2]
	if action == "cancel" && voteID == "none" {
		guestSpamTelegram.AnswerCallback(callback.ID, false, "已取消")
		guestSpamTelegram.DeleteCallbackMessage(callback)
		return nil
	}
	vote, ok := GetVote(voteID)
	if !ok {
		guestSpamTelegram.AnswerCallback(callback.ID, true, "投票已过期")
		return nil
	}
	if action == "cancel" {
		if !canCancelVote(vote, callback.From) {
			guestSpamTelegram.AnswerCallback(callback.ID, true, "只有发起者或管理员可以取消投票")
			return nil
		}
		DeleteVote(voteID)
		guestSpamTelegram.AnswerCallback(callback.ID, false, "已取消")
		guestSpamTelegram.DeleteCallbackMessage(callback)
		return nil
	}
	if action != "vote" {
		return nil
	}
	if callback.From == nil || callback.From.IsBot {
		guestSpamTelegram.AnswerCallback(callback.ID, true, "机器人不能参与投票")
		return nil
	}
	if containsUser(vote.Voters, callback.From.ID) {
		guestSpamTelegram.AnswerCallback(callback.ID, true, "你已经投过票了")
		return nil
	}
	weight, voterTier := voteWeight(vote.ChatID, callback.From.ID)
	if weight == 0 {
		guestSpamTelegram.AnswerCallback(callback.ID, true, "最近不够活跃，不能参与投票")
		return nil
	}
	score := currentVoteScore(vote)
	vote.Voters = append(vote.Voters, callback.From.ID)
	vote.VoteScore = score + weight
	SaveVote(vote)
	if vote.VoteScore >= vote.RequiredVoteCount {
		return applyVotePassed(vote, callback)
	}
	guestSpamTelegram.AnswerCallback(callback.ID, false, fmt.Sprintf("已投票 %d/%d（%s +%d）", vote.VoteScore, vote.RequiredVoteCount, voterTier, weight))
	return nil
}

func sendRecentGuestCandidates(message *tgbotapi.Message) error {
	recents := RecentGuestMessages(message.Chat.ID)
	if len(recents) == 0 {
		return sendTempMessage(message.Chat.ID, message.MessageID, "最近没有可判定的 guest bot 消息。")
	}
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(recents)+1)
	text := "请选择要发起 spam 判定的 guest bot 消息："
	for i, item := range recents {
		label := fmt.Sprintf("%d. %s", i+1, displayBot(item))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, fmt.Sprintf("guestspam_select,%d", item.MessageID)),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("取消", "guestspam_vote,cancel,none"),
	))
	send := tgbotapi.NewMessage(message.Chat.ID, text)
	send.ReplyToMessageID = message.MessageID
	send.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg, err := guestSpamTelegram.Send(send)
	if err == nil {
		guestSpamTelegram.QueueDelete(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
	}
	return err
}

func SelectRecentGuestCallback(update tgbotapi.Update) error {
	callback := update.CallbackQuery
	if callback == nil || callback.Message == nil {
		return nil
	}
	data := strings.Split(callback.Data, ",")
	if len(data) < 2 {
		return nil
	}
	messageID, _ := strconv.Atoi(data[1])
	var selected RecentGuestMessage
	for _, item := range RecentGuestMessages(callback.Message.Chat.ID) {
		if item.MessageID == messageID {
			selected = item
			break
		}
	}
	if selected.MessageID == 0 {
		guestSpamTelegram.AnswerCallback(callback.ID, true, "候选消息已过期")
		return nil
	}
	msg := &tgbotapi.Message{
		MessageID: callback.Message.MessageID,
		Chat:      callback.Message.Chat,
		From:      callback.From,
	}
	started, err := startSpamVote(msg, selected)
	if err == nil && started {
		guestSpamTelegram.DeleteCallbackMessage(callback)
		guestSpamTelegram.AnswerCallback(callback.ID, false, "已发起投票")
	}
	return err
}

func startSpamVote(message *tgbotapi.Message, selected RecentGuestMessage) (bool, error) {
	if message == nil || message.Chat == nil || message.From == nil {
		return false, nil
	}
	activeCount := ActiveUserCount(message.Chat.ID)
	required, ok := requiredVoteCount(activeCount)
	if !ok {
		AddLog(logFromRecent(selected, ActionVoteInvalid, ReasonInsufficientAct, fmt.Sprintf("active users: %d", activeCount)))
		return false, sendTempMessage(message.Chat.ID, message.MessageID, "最近 10 分钟活跃人数少于 3，投票无效，请管理员检查。")
	}
	voteID, _ := gonanoid.New(16)
	vote := SpamVote{
		ID:                voteID,
		ChatID:            selected.ChatID,
		ChatName:          selected.ChatName,
		MessageID:         selected.MessageID,
		GuestBotID:        selected.GuestBotID,
		GuestBotName:      selected.GuestBotName,
		GuestBotUserName:  selected.GuestBotUserName,
		StarterUserID:     message.From.ID,
		StarterUserName:   message.From.FullName(),
		ActiveUserCount:   activeCount,
		RequiredVoteCount: required,
		Voters:            []int64{},
		CreatedAt:         time.Now(),
		ExpiresAt:         time.Now().Add(10 * time.Minute),
	}
	text := fmt.Sprintf("是否将 guest bot %s 判定为 spam？\n10 分钟内需要 %d 票，当前 0/%d。", displayBot(selected), required, required)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("判定 spam", fmt.Sprintf("guestspam_vote,vote,%s", voteID)),
			tgbotapi.NewInlineKeyboardButtonData("取消", fmt.Sprintf("guestspam_vote,cancel,%s", voteID)),
		),
	)
	send := tgbotapi.NewMessage(message.Chat.ID, text)
	send.ReplyToMessageID = message.MessageID
	send.ReplyMarkup = keyboard
	msg, err := guestSpamTelegram.Send(send)
	if err != nil {
		return false, err
	}
	vote.VoteMessageID = msg.MessageID
	SaveVote(vote)
	AddLog(logFromRecent(selected, ActionVoteStarted, ReasonVote, fmt.Sprintf("required votes: %d", required)))
	return true, nil
}

func applyVotePassed(vote SpamVote, callback *tgbotapi.CallbackQuery) error {
	ApplyVotePassedState(vote)
	if err := deleteVoteGuestMessageWithLog(vote); err != nil {
		guestSpamTelegram.AnswerCallback(callback.ID, false, "投票通过，已拉黑，但删除消息失败，请管理员检查权限")
		guestSpamTelegram.DeleteCallbackMessage(callback)
		return err
	}
	guestSpamTelegram.AnswerCallback(callback.ID, false, "投票通过，已拉黑并删除消息")
	guestSpamTelegram.DeleteCallbackMessage(callback)
	return nil
}

func deleteVoteGuestMessageWithLog(vote SpamVote) error {
	message := &tgbotapi.Message{
		MessageID: vote.MessageID,
		Chat: &tgbotapi.Chat{
			ID:    vote.ChatID,
			Title: vote.ChatName,
		},
		From: &tgbotapi.User{
			ID:       vote.GuestBotID,
			IsBot:    true,
			UserName: vote.GuestBotUserName,
		},
	}
	if vote.GuestBotName != "" {
		message.From.FirstName = vote.GuestBotName
	}
	return deleteGuestMessageWithLog(message, ReasonVote)
}

func ApplyVotePassedState(vote SpamVote) {
	item := GuestBotBlacklist{
		BotID:          vote.GuestBotID,
		BotName:        vote.GuestBotName,
		BotUserName:    vote.GuestBotUserName,
		Source:         "vote",
		FirstChatID:    vote.ChatID,
		FirstMessageID: vote.MessageID,
	}
	AddBlacklist(item, true)
	AddLog(SpamLog{
		ChatID:       vote.ChatID,
		ChatName:     vote.ChatName,
		MessageID:    vote.MessageID,
		GuestBotID:   vote.GuestBotID,
		GuestBotName: vote.GuestBotName,
		GuestBotUser: vote.GuestBotUserName,
		Action:       ActionVotePassed,
		Reason:       ReasonVote,
		Detail:       fmt.Sprintf("votes: %d/%d", currentVoteScore(vote), vote.RequiredVoteCount),
	})
	DeleteVote(vote.ID)
}

func restoreCaller(chatID, userID int64, message *tgbotapi.Message) error {
	if _, err := guestSpamTelegram.UnbanChatMember(chatID, userID); err != nil {
		item := SpamLog{
			ChatID:       chatID,
			CallerUserID: userID,
			Action:       ActionRestoreCaller,
			Reason:       ReasonAdminRestore,
			Detail:       fmt.Sprintf("unban failed: %v", err),
		}
		if message != nil && message.Chat != nil {
			item.ChatName = message.Chat.Title
		}
		AddLog(item)
		return err
	}
	if _, err := guestSpamTelegram.RestrictChatMember(chatID, userID, tgbotapi.AllPermissions); err != nil {
		item := SpamLog{
			ChatID:       chatID,
			CallerUserID: userID,
			Action:       ActionRestoreCaller,
			Reason:       ReasonAdminRestore,
			Detail:       fmt.Sprintf("unrestrict failed: %v", err),
		}
		if message != nil && message.Chat != nil {
			item.ChatName = message.Chat.Title
		}
		AddLog(item)
		return err
	}
	RestoreCallerState(chatID, userID, message)
	return sendTempMessage(chatID, message.MessageID, "已恢复该用户并清除 guest spam 警告。")
}

func RestoreCallerState(chatID, userID int64, message *tgbotapi.Message) {
	ClearWarnings(chatID, userID)
	item := SpamLog{
		ChatID:       chatID,
		CallerUserID: userID,
		Action:       ActionRestoreCaller,
		Reason:       ReasonAdminRestore,
		Detail:       "restored",
	}
	if message != nil {
		if message.Chat != nil {
			item.ChatName = message.Chat.Title
		}
		if message.From != nil {
			item.Detail = fmt.Sprintf("restored by %s", message.From.FullName())
		}
	}
	AddLog(item)
}

func sendLogs(chatID int64, replyTo int) error {
	logs := RecentLogs(chatID, 10)
	if len(logs) == 0 {
		return sendTempMessage(chatID, replyTo, "暂无 guest spam 日志。")
	}
	var builder strings.Builder
	builder.WriteString("最近 guest spam 日志：\n")
	for _, item := range logs {
		builder.WriteString(fmt.Sprintf("%s %s bot=%s caller=%s %s\n",
			item.CreateTime.Format("01-02 15:04"),
			item.Action,
			item.GuestBotUser,
			callerDisplay(item),
			item.Detail,
		))
	}
	return sendTempMessage(chatID, replyTo, builder.String())
}

func sendTempMessage(chatID int64, replyTo int, text string) error {
	send := tgbotapi.NewMessage(chatID, text)
	send.ReplyToMessageID = replyTo
	msg, err := guestSpamTelegram.Send(send)
	if err == nil {
		guestSpamTelegram.QueueDelete(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
	}
	return err
}

func containsUser(users []int64, userID int64) bool {
	for _, id := range users {
		if id == userID {
			return true
		}
	}
	return false
}

func canCancelVote(vote SpamVote, user *tgbotapi.User) bool {
	if user == nil || user.IsBot {
		return false
	}
	if vote.StarterUserID != 0 && vote.StarterUserID == user.ID {
		return true
	}
	return guestSpamTelegram.IsAdmin(vote.ChatID, user.ID)
}

func voteWeight(chatID, userID int64) (int, string) {
	if guestSpamTelegram.IsAdmin(chatID, userID) {
		return adminVoteWeight, "管理员"
	}
	if IsTrustedMember(chatID, userID) {
		return trustedVoteWeight, "可信成员"
	}
	if IsActiveUser(chatID, userID) {
		return activeVoteWeight, "活跃成员"
	}
	return 0, ""
}

func currentVoteScore(vote SpamVote) int {
	if vote.VoteScore > 0 {
		return vote.VoteScore
	}
	return len(vote.Voters)
}

func requiredVoteCount(activeCount int) (int, bool) {
	if activeCount < 3 {
		return 0, false
	}
	return (activeCount + 1) / 2, true
}

func displayBot(item RecentGuestMessage) string {
	if item.GuestBotUserName != "" {
		return "@" + item.GuestBotUserName
	}
	if item.GuestBotName != "" {
		return item.GuestBotName
	}
	return fmt.Sprintf("%d", item.GuestBotID)
}

func callerDisplay(item SpamLog) string {
	if item.CallerUserID != 0 {
		return fmt.Sprintf("%s(%d)", item.CallerUserName, item.CallerUserID)
	}
	if item.CallerChatID != 0 {
		return fmt.Sprintf("%s(%d)", item.CallerChatName, item.CallerChatID)
	}
	return "-"
}

func logFromRecent(item RecentGuestMessage, action string, reason string, detail string) SpamLog {
	return SpamLog{
		ChatID:         item.ChatID,
		ChatName:       item.ChatName,
		MessageID:      item.MessageID,
		GuestBotID:     item.GuestBotID,
		GuestBotName:   item.GuestBotName,
		GuestBotUser:   item.GuestBotUserName,
		CallerUserID:   item.CallerUserID,
		CallerUserName: item.CallerUserName,
		CallerChatID:   item.CallerChatID,
		CallerChatName: item.CallerChatName,
		Action:         action,
		Reason:         reason,
		Detail:         detail,
	}
}
