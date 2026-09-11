package account

import (
	"arknights_bot/config"
	"arknights_bot/utils/repo"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"strings"
	"sync"
)

// 指令并发执行后，这两个 map 会被多个 goroutine 访问，需要加锁保护。
var (
	accountMapMu  sync.Mutex
	sklandIdMap   = make(map[int64]string)
	serverNameMap = make(map[int64]string)
)

func setServerName(chatId int64, name string) {
	accountMapMu.Lock()
	serverNameMap[chatId] = name
	accountMapMu.Unlock()
}
func getServerName(chatId int64) string {
	accountMapMu.Lock()
	defer accountMapMu.Unlock()
	return serverNameMap[chatId]
}
func setSklandId(chatId int64, id string) {
	accountMapMu.Lock()
	sklandIdMap[chatId] = id
	accountMapMu.Unlock()
}
func getSklandId(chatId int64) string {
	accountMapMu.Lock()
	defer accountMapMu.Unlock()
	return sklandIdMap[chatId]
}
func delSklandId(chatId int64) {
	accountMapMu.Lock()
	delete(sklandIdMap, chatId)
	accountMapMu.Unlock()
}

// ChooseServer 选择服务器
func ChooseServer(callBack tgbotapi.Update) error {
	callbackQuery := callBack.CallbackQuery
	callbackQuery.Answer(false, "")
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 3 {
		return nil
	}

	chatId := callbackQuery.Message.Chat.ID
	userId := callbackQuery.From.ID
	setServerName(chatId, d[1])
	operType := d[2]

	sendMessage := tgbotapi.NewMessage(chatId,
		"🔑 *如何获取 Token*\n\n"+
			"🌏 *国服*\n"+
			"1️⃣ 前往森空岛登录\n"+
			"2️⃣ 打开下方「获取国服 Token」按钮\n"+
			"3️⃣ 复制 `content` 中的 token\n\n"+
			"🌍 *国际服*\n"+
			"1️⃣ 前往森空港登录\n"+
			"2️⃣ 打开下方「获取国际服 Token」按钮\n"+
			"3️⃣ 复制 `content` 中的 token\n\n"+
			"请直接发送 token，或使用 `/cancel` 取消",
	)
	sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
	sendMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🌏 森空岛登录", "https://www.skland.com"),
			tgbotapi.NewInlineKeyboardButtonURL("🔑 获取国服 Token", "https://web-api.skland.com/account/info/hg"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🌍 森空港登录", "https://www.skport.com"),
			tgbotapi.NewInlineKeyboardButtonURL("🔑 获取国际服 Token", "https://web-api.skport.com/cookie_store/account_token"),
		),
	)
	config.Arknights.SetWaitMessage(userId, operType)
	config.Arknights.Send(sendMessage)
	callbackQuery.Message.Delete()
	return nil
}

// ChoosePlayer 选择绑定角色
func ChoosePlayer(callBack tgbotapi.Update) error {
	callbackQuery := callBack.CallbackQuery
	callbackQuery.Answer(false, "")
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 4 {
		return nil
	}

	userId := callbackQuery.From.ID
	chatId := callbackQuery.Message.Chat.ID

	uid := d[1]
	serverName := d[2]
	playerName := d[3]
	sklandId := getSklandId(chatId)

	var userAccount UserAccount
	var userPlayer UserPlayer
	repo.GetAccountByUserIdAndSklandId(userId, sklandId).Scan(&userAccount)
	res := repo.GetPlayerByUserId(userId, uid).Scan(&userPlayer)
	if res.RowsAffected == 0 {
		id, _ := gonanoid.New(32)
		userPlayer = UserPlayer{
			Id:         id,
			AccountId:  userAccount.Id,
			UserName:   userAccount.UserName,
			UserNumber: userAccount.UserNumber,
			Uid:        uid,
			ServerName: serverName,
			PlayerName: playerName,
		}
		config.DBEngine.Table("user_player").Create(&userPlayer)
	} else {
		userPlayer.PlayerName = playerName
		userPlayer.ServerName = serverName
		config.DBEngine.Table("user_player").Save(&userPlayer)
		config.Arknights.SendText(chatId, "此角色已绑定，更新角色信息。")
		return nil
	}
	config.Arknights.SendText(chatId, "角色绑定成功！")
	delSklandId(chatId)
	return nil
}

// UnbindPlayer 解绑角色
func UnbindPlayer(callBack tgbotapi.Update) error {
	callbackQuery := callBack.CallbackQuery
	callbackQuery.Answer(false, "")
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 2 {
		return nil
	}

	userId := callbackQuery.From.ID
	chatId := callbackQuery.Message.Chat.ID

	uid := d[1]
	config.DBEngine.Exec("delete from user_player where user_number = ? and uid = ?", userId, uid)
	config.Arknights.SendText(chatId, "角色解绑成功！")
	callbackQuery.Message.Delete()
	return nil
}
