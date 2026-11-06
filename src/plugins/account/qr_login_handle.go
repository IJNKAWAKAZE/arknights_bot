package account

import (
	"arknights_bot/config"
	"arknights_bot/plugins/skland"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	qrcode "github.com/skip2/go-qrcode"
)

const (
	qrServerName   = "国服"              // 扫码登录只做国服
	qrPollInterval = 2 * time.Second   // 轮询间隔
	qrPollTimeout  = 120 * time.Second // 二维码有效期（服务端约 2 分钟）

	qrModeBind  = "bind"  // 扫码绑定：落库 + 选角色
	qrModeReset = "reset" // 扫码重置：只刷新已有账号的凭据
)

// QRLoginHandle 森空岛扫码登录（国服）
//
// callback data: qrLogin            -> 绑定
// callback data: qrLogin,reset      -> 重设 token
//
// 流程：gen_scan/login 拿二维码 -> 发图 -> 轮询 scan_status -> token_by_scan_code 换 token -> 绑定/重置
func QRLoginHandle(callBack tgbotapi.Update) error {
	callbackQuery := callBack.CallbackQuery
	callbackQuery.Answer(false, "")

	mode := qrModeBind
	if d := strings.Split(callBack.CallbackData(), ","); len(d) >= 2 && d[1] == qrModeReset {
		mode = qrModeReset
	}

	chatId := callbackQuery.Message.Chat.ID
	userId := callbackQuery.From.ID
	userName := callbackQuery.From.FullName()

	setServerName(chatId, qrServerName)
	callbackQuery.Message.Delete()

	login, err := skland.GenScanLogin()
	if err != nil {
		config.Arknights.SendText(chatId, err.Error())
		return err
	}

	png, err := qrcode.Encode(login.ScanUrl, qrcode.Medium, 512)
	if err != nil {
		config.Arknights.SendText(chatId, "二维码生成失败，请重试")
		return err
	}

	photo := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: png, Name: "qr_login.png"})
	if mode == qrModeReset {
		photo.Caption = "📷 请用「森空岛」App 扫码，刷新当前账号的登录凭证\n" +
			"扫码后还需在手机上点击确认，二维码 2 分钟内有效\n" +
			"国际服请返回上一页，使用 Token 方式重设"
	} else {
		photo.Caption = "📷 请用「森空岛」App 扫码登录\n" +
			"扫码后还需在手机上点击确认，二维码 2 分钟内有效\n" +
			"国际服请返回上一页，使用 Token 方式绑定"
	}
	config.Arknights.Send(photo)

	go pollScanLogin(chatId, userId, userName, mode, login.ScanId)
	return nil
}

// pollScanLogin 轮询扫码状态，确认后换取 token 并完成绑定/重置
func pollScanLogin(chatId, userId int64, userName, mode, scanId string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("扫码登录轮询异常:", r)
		}
	}()

	deadline := time.Now().Add(qrPollTimeout)
	notified := false
	for time.Now().Before(deadline) {
		time.Sleep(qrPollInterval)

		status, scanCode, err := skland.GetScanStatus(scanId)
		if err != nil {
			log.Println(err)
			continue
		}
		switch status {
		case skland.ScanStatusScanned:
			if !notified {
				config.Arknights.SendText(chatId, "已扫码，请在手机上点击「确认登录」")
				notified = true
			}
		case skland.ScanStatusExpired:
			config.Arknights.SendText(chatId, "二维码已过期，请重新操作")
			return
		case skland.ScanStatusConfirmed:
			token, err := skland.TokenByScanCode(scanCode)
			if err != nil {
				config.Arknights.SendText(chatId, "登录凭证获取失败："+err.Error())
				return
			}
			account, err := skland.Login(token, qrServerName)
			if err != nil {
				config.Arknights.SendText(chatId, "登录失败！请重新扫码。")
				return
			}
			if mode == qrModeReset {
				if err := finishReset(chatId, userId, token, account); err != nil {
					log.Println("扫码重设失败:", err)
				}
			} else if err := finishBind(chatId, userId, userName, token, account); err != nil {
				log.Println("扫码绑定失败:", err)
			}
			return
		}
	}
	config.Arknights.SendText(chatId, "二维码已过期，请重新操作")
}
