package sign

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"crypto/rand"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/big"
	"strconv"
	"time"
)

// AutoSign 森空岛自动签到
func AutoSign() func() {
	return func() {
		var users []UserSign
		res := utils.GetAutoSign().Scan(&users)
		if res.RowsAffected > 0 {
			log.Println("开始执行自动签到...")
			// 遍历所有自动签到用户
			for _, user := range users {
				r, _ := rand.Int(rand.Reader, big.NewInt(60))
				random, _ := strconv.Atoi(r.String())
				time.Sleep(time.Second * time.Duration(random))
				sign(user)
			}
			log.Println("自动签到执行完毕...")
		}
	}
}

func sign(user UserSign) {
	var players []account.UserPlayer
	res := utils.GetPlayersByUserId(user.UserNumber).Scan(&players)
	if res.RowsAffected > 0 {

		// 用户绑定了角色
		var skAccount skland.Account
		var userAccount account.UserAccount
		res := utils.GetAccountByUserId(user.UserNumber).Scan(&userAccount)
		if res.RowsAffected > 0 {

			// 获取用户账号信息
			skAccount.Hypergryph.Token = userAccount.HypergryphToken
			skAccount.Skland.Token = userAccount.SklandToken
			skAccount.Skland.Cred = userAccount.SklandCred

			// 对所有绑定角色执行签到
			for _, player := range players {
				var skPlayer skland.Player
				skPlayer.NickName = player.PlayerName
				skPlayer.ChannelName = player.ServerName
				skPlayer.Uid = player.Uid

				// 执行签到
				record, err := skland.SignGamePlayer(&skPlayer, skAccount)
				if err != nil {
					// 签到失败
					sendMessage := tgbotapi.NewMessage(user.UserNumber, fmt.Sprintf("角色 %s 签到失败!\nmsg:%s", player.PlayerName, err.Error()))
					bot.Arknights.Send(sendMessage)
					log.Println(player.PlayerName, err)
					return
				}
				// 今日已完成签到
				if record.HasSigned {
					sendMessage := tgbotapi.NewMessage(user.UserNumber, fmt.Sprintf("角色 %s 今天已经签到过了", player.PlayerName))
					bot.Arknights.Send(sendMessage)
					return
				}
				// 签到成功
				sendMessage := tgbotapi.NewMessage(user.UserNumber, fmt.Sprintf("角色 %s 签到成功!\n今日奖励：%s", player.PlayerName, record.Award))
				bot.Arknights.Send(sendMessage)
			}
		}
	}
}
