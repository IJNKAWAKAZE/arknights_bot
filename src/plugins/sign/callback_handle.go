package sign

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

// SignPlayer 选择签到角色
func SignPlayer(callBack tgbotapi.Update) (bool, error) {
	callbackQuery := callBack.CallbackQuery
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 3 {
		return true, nil
	}

	userId := callbackQuery.From.ID
	chatId := callbackQuery.Message.Chat.ID
	clickUserId, _ := strconv.ParseInt(d[1], 10, 64)
	uid := d[2]

	if userId != clickUserId {
		answer := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "这不是你的角色！")
		bot.Arknights.Send(answer)
		return true, nil
	}

	var userAccount account.UserAccount
	var skPlayer *skland.Player
	var player account.UserPlayer
	var skAccount skland.Account

	utils.GetAccountByUserId(userId).Scan(&userAccount)
	utils.GetPlayerByUserId(userId, uid).Scan(&player)

	skPlayer.NickName = player.PlayerName
	skPlayer.ChannelName = player.ServerName
	skPlayer.Uid = player.Uid
	skAccount.Hypergryph.Token = userAccount.HypergryphToken
	skAccount.Skland.Token = userAccount.SklandToken
	skAccount.Skland.Cred = userAccount.SklandCred
	record, err := skland.SignGamePlayer(skPlayer, skAccount)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("角色 %s 今天已经签到过了", player.PlayerName))
		msg, _ := bot.Arknights.Send(sendMessage)
		go utils.DelayDelMsg(msg.Chat.ID, msg.MessageID, viper.GetDuration("bot.msg_del_delay"))
		return true, nil
	}

	sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("角色 %s 签到成功!\n今日奖励：%s", player.PlayerName, record.Award))
	msg, _ := bot.Arknights.Send(sendMessage)
	go utils.DelayDelMsg(msg.Chat.ID, msg.MessageID, viper.GetDuration("bot.msg_del_delay"))
	return true, nil
}
