package utils

import (
	initBot "arknights_bot/bot/init"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SendMessage 发送文本
func SendMessage(message tgbotapi.MessageConfig) (tgbotapi.Message, error) {
	return initBot.Kawakaze.Send(message)
}

// SendSticker 发送贴纸
func SendSticker(sticker tgbotapi.StickerConfig) (tgbotapi.Message, error) {
	return initBot.Kawakaze.Send(sticker)
}

// SendPhoto 发送图片
func SendPhoto(photo tgbotapi.PhotoConfig) (tgbotapi.Message, error) {
	return initBot.Kawakaze.Send(photo)
}

// SendAnimation 发送动图
func SendAnimation(animation tgbotapi.AnimationConfig) (tgbotapi.Message, error) {
	return initBot.Kawakaze.Send(animation)
}

// SendMediaGroup 发送媒体组
func SendMediaGroup(mediaGroup tgbotapi.MediaGroupConfig) ([]tgbotapi.Message, error) {
	return initBot.Kawakaze.SendMediaGroup(mediaGroup)
}

// GetChatMemberInfo 获取成员信息
func GetChatMemberInfo(member tgbotapi.GetChatMemberConfig) (tgbotapi.ChatMember, error) {
	return initBot.Kawakaze.GetChatMember(member)
}

// SetMemberPermissions 设置群员权限
func SetMemberPermissions(permissions tgbotapi.RestrictChatMemberConfig) (tgbotapi.Message, error) {
	return initBot.Kawakaze.Send(permissions)
}

// KickChatMember 踢出群员
func KickChatMember(kickMember tgbotapi.KickChatMemberConfig) (tgbotapi.Message, error) {
	return initBot.Kawakaze.Send(kickMember)
}

// UnbanChatMember 解除成员黑名单
func UnbanChatMember(unbanMember tgbotapi.UnbanChatMemberConfig) (tgbotapi.Message, error) {
	return initBot.Kawakaze.Send(unbanMember)
}

// SendCallbackAnswer 发送回调响应
func SendCallbackAnswer(answer tgbotapi.CallbackConfig) (tgbotapi.Message, error) {
	return initBot.Kawakaze.Send(answer)
}

// DeleteMessage 删除消息
func DeleteMessage(delMsg tgbotapi.DeleteMessageConfig) (tgbotapi.Message, error) {
	return initBot.Kawakaze.Send(delMsg)
}

// SendPoll 发送投票
func SendPoll(poll tgbotapi.SendPollConfig) (tgbotapi.Message, error) {
	return initBot.Kawakaze.Send(poll)
}

// SendAction 发送动作
func SendAction(action tgbotapi.ChatActionConfig) (tgbotapi.Message, error) {
	return initBot.Kawakaze.Send(action)
}
