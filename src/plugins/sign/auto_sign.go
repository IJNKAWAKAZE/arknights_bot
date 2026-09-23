package sign

import (
	"arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils/repo"
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	bulkMu                  sync.Mutex
	bulkContext, cancelBulk = context.WithCancel(context.Background())
	ErrAlreadyRunning       = errors.New("已有批量签到任务正在执行")
)

// AutoSign 同步执行签到，便于定时任务调度器跟踪执行状态并等待结束。
func AutoSign() {
	if err := RunAutoSign(); err != nil {
		log.Println("自动签到:", err)
	}
}
func Stop() { cancelBulk() }

func RunAutoSign() error {
	if !bulkMu.TryLock() {
		return ErrAlreadyRunning
	}
	defer bulkMu.Unlock()
	if err := bulkContext.Err(); err != nil {
		return err
	}
	var users []UserSign
	if err := repo.GetAutoSign().Scan(&users).Error; err != nil {
		return err
	}
	log.Println("开始执行自动签到...")
	for _, user := range users {
		if err := waitSignDelay(bulkContext, time.Duration(rand.Intn(60))*time.Second); err != nil {
			return err
		}
		if err := signUser(bulkContext, user, skland.SignGamePlayer, func(id int64, text string) {
			if _, err := config.Arknights.SendText(id, text); err != nil {
				log.Println("签到通知失败:", err)
			}
		}); err != nil {
			if bulkContext.Err() != nil {
				return bulkContext.Err()
			}
			log.Printf("用户 %d 签到失败: %v", user.UserNumber, err)
		}
	}
	log.Println("自动签到执行完毕...")
	return nil
}

func waitSignDelay(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return ctx.Err()
	}
}

func signUser(ctx context.Context, user UserSign,
	signPlayer func(string, skland.Account, string) (string, bool, error),
	notify func(int64, string)) error {
	var players []account.UserPlayer
	if err := repo.GetPlayersByUserId(user.UserNumber).Scan(&players).Error; err != nil {
		return err
	}
	for _, player := range players {
		if err := ctx.Err(); err != nil {
			return err
		}
		var ua account.UserAccount
		res := repo.GetAccountByUid(user.UserNumber, player.Uid).Scan(&ua)
		if res.Error != nil {
			log.Println("查询签到账号失败:", res.Error)
			continue
		}
		if res.RowsAffected == 0 {
			continue
		}
		var ska skland.Account
		ska.Hypergryph.Token = ua.HypergryphToken
		ska.Skland.Token = ua.SklandToken
		ska.Skland.Cred = ua.SklandCred
		award, alreadySigned, err := signPlayer(player.Uid, ska, ua.ServerName)
		if err != nil {
			if user.NotifyMode == 0 || user.NotifyMode == 1 {
				notify(user.UserNumber, fmt.Sprintf("角色 %s 签到失败!\n失败原因:%s", player.PlayerName, err))
			}
			log.Println(player.PlayerName, err)
			continue
		}
		if alreadySigned {
			if user.NotifyMode == 0 {
				notify(user.UserNumber, fmt.Sprintf("角色 %s 今天已经签到过了", player.PlayerName))
			}
			continue
		}
		if user.NotifyMode == 0 || user.NotifyMode == 2 {
			notify(user.UserNumber, fmt.Sprintf("角色 %s 签到成功!\n今日奖励：%s", player.PlayerName, award))
		}
	}
	return nil
}
