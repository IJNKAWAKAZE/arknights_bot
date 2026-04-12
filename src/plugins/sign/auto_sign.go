package sign

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"crypto/rand"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"log"
	"math/big"
	"strconv"
	"time"
)

// AutoSign 森空岛自动签到
func AutoSign() {
	var users []UserSign
	res := utils.GetAutoSign().Scan(&users)
	if res.RowsAffected > 0 {
		go func() {
			log.Println("开始执行自动签到...")
			// 遍历所有自动签到用户
			for _, user := range users {
				r, _ := rand.Int(rand.Reader, big.NewInt(60))
				random, _ := strconv.Atoi(r.String())
				time.Sleep(time.Second * time.Duration(random))
				sign(user)
			}
			log.Println("自动签到执行完毕...")
		}()
	}
}

func sign(user UserSign) {
	var players []account.UserPlayer
	res := utils.GetPlayersByUserId(user.UserNumber).Scan(&players)
	if res.RowsAffected > 0 {
		// 对所有绑定角色执行签到
		for _, player := range players {
			var skAccount skland.Account
			var userAccount account.UserAccount
			// 获取用户账号信息
			res := utils.GetAccountByUid(user.UserNumber, player.Uid).Scan(&userAccount)
			if res.RowsAffected > 0 {
				skAccount.Hypergryph.Token = userAccount.HypergryphToken
				skAccount.Skland.Token = userAccount.SklandToken
				skAccount.Skland.Cred = userAccount.SklandCred

				// 执行签到
				award, hasSigned, err := skland.SignGamePlayer(player.Uid, skAccount, userAccount.ServerName)
				if err != nil {
					// 签到失败 - notify_mode: 0(全部通知) 或 1(仅失败通知) 时发送
					if user.NotifyMode == 0 || user.NotifyMode == 1 {
						sendMessage := tgbotapi.NewMessage(user.UserNumber, fmt.Sprintf("角色 %s 签到失败!\n失败原因:%s", player.PlayerName, err.Error()))
						bot.Arknights.Send(sendMessage)
					}
					log.Println(player.PlayerName, err)
					return
				}
				// 今日已完成签到 - notify_mode: 0(全部通知) 时发送
				if hasSigned {
					if user.NotifyMode == 0 {
						sendMessage := tgbotapi.NewMessage(user.UserNumber, fmt.Sprintf("角色 %s 今天已经签到过了", player.PlayerName))
						bot.Arknights.Send(sendMessage)
					}
					return
				}
				// 签到成功 - notify_mode: 0(全部通知) 或 2(仅成功通知) 时发送
				if user.NotifyMode == 0 || user.NotifyMode == 2 {
					sendMessage := tgbotapi.NewMessage(user.UserNumber, fmt.Sprintf("角色 %s 签到成功!\n今日奖励：%s", player.PlayerName, award))
					bot.Arknights.Send(sendMessage)
				}
			}
		}
	}
}
