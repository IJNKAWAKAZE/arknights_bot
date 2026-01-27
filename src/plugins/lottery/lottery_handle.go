package lottery

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// StartLotteryHandle 开启抽奖活动
func StartLotteryHandle(update tgbotapi.Update) error {
	message := update.Message
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	param := update.Message.CommandArguments()
	messageId := message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 5)

	if !bot.Arknights.IsAdminWithPermissions(chatId, userId, tgbotapi.AdminCanRestrictMembers) {
		sendMessage := tgbotapi.NewMessage(chatId, "无使用权限！")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}
	// 检查参数是否为 24 小时制时间的年月日时分秒格式
	var endTime time.Time
	if param != "" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", param, time.Local)
		if err != nil {
			sendMessage := tgbotapi.NewMessage(chatId, "参数格式错误！请输入 YYYY-MM-DD HH:MM:SS 格式")
			sendMessage.ReplyToMessageID = messageId
			msg, err := bot.Arknights.Send(sendMessage)
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
			return nil
		}
		endTime = t
	}

	// 如果没有设置报名截止时间（参数为空）
	if endTime.IsZero() {
		endTime = time.Now().Add(time.Hour * 24 * 7) // 默认 7 天后截止报名
	}
	// 检查是否存在已开启的抽奖
	var lottery utils.GroupLottery
	utils.GetGroupLottery(chatId).Scan(&lottery)
	if lottery.Id != "" {
		sendMessage := tgbotapi.NewMessage(chatId, "已有正在进行的抽奖活动，请先结束当前抽奖！")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	// 保存抽奖记录
	id, _ := gonanoid.New(32)
	groupLottery := utils.GroupLottery{
		Id:          id,
		GroupName:   message.Chat.Title,
		GroupNumber: message.Chat.ID,
		Status:      1,
		EndTime:     endTime,
	}
	res := bot.DBEngine.Table("group_lottery").Create(&groupLottery)
	log.Println(res.Error)
	sendMessage := tgbotapi.NewMessage(chatId, tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, fmt.Sprintf("🎉 *抽奖活动已开启*\n\n📅 *报名截止时间*：%s\n\n📝 *指令说明*：\n🔹 参与选号：`/join_lottery [1-100]`\n🔹 查看详情：`/lottery_detail`\n\n⚙️ *管理指令*：\n🔸 停止报名：`/stop_lottery`\n🔸 进行抽奖：`/lottery`\n🔸 结束抽奖：`/end_lottery`", endTime.Format("2006-01-02 15:04:05"))))
	sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
	sendMessage.ReplyToMessageID = messageId
	bot.Arknights.Send(sendMessage)
	return nil
}

// StopLotteryHandle 停止抽奖报名
func StopLotteryHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 5)
	if !bot.Arknights.IsAdminWithPermissions(chatId, userId, tgbotapi.AdminCanRestrictMembers) {
		sendMessage := tgbotapi.NewMessage(chatId, "无使用权限！")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}
	var lottery utils.GroupLottery
	utils.GetGroupLottery(chatId).Scan(&lottery)
	if lottery.Id == "" {
		sendMessage := tgbotapi.NewMessage(chatId, "当前群组暂无正在进行的抽奖活动！")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}
	if lottery.Status == 2 {
		sendMessage := tgbotapi.NewMessage(chatId, "抽奖活动已停止报名，请勿重复操作！")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}
	lottery.Status = 2
	bot.DBEngine.Table("group_lottery").Save(&lottery)
	sendMessage := tgbotapi.NewMessage(chatId, "抽奖活动已停止报名！")
	sendMessage.ReplyToMessageID = messageId
	msg, err := bot.Arknights.Send(sendMessage)
	if err != nil {
		return err
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
	return nil
}

// EndLotteryHandle 结束抽奖活动
func EndLotteryHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 5)
	if !bot.Arknights.IsAdminWithPermissions(chatId, userId, tgbotapi.AdminCanRestrictMembers) {
		sendMessage := tgbotapi.NewMessage(chatId, "无使用权限！")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}
	var lottery utils.GroupLottery
	utils.GetGroupLottery(chatId).Scan(&lottery)
	if lottery.Id == "" {
		sendMessage := tgbotapi.NewMessage(chatId, "当前群组暂无抽奖活动！")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}
	lottery.Status = 0
	bot.DBEngine.Table("group_lottery").Save(&lottery)
	sendMessage := tgbotapi.NewMessage(chatId, "抽奖活动已结束！")
	sendMessage.ReplyToMessageID = messageId
	msg, err := bot.Arknights.Send(sendMessage)
	if err != nil {
		return err
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
	return nil
}

// JoinLotteryHandle 参加抽奖活动
func JoinLotteryHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	param := update.Message.CommandArguments()
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 5)

	// 检查当前群组是否存在已开启的抽奖
	var lottery utils.GroupLottery
	utils.GetGroupLottery(chatId).Scan(&lottery)
	if lottery.Id == "" {
		sendMessage := tgbotapi.NewMessage(chatId, "当前群组暂无正在进行的抽奖活动！")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	if lottery.Status == 2 {
		sendMessage := tgbotapi.NewMessage(chatId, "抽奖活动已停止报名！")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	// 检查用户输入的数字是否合法 (1-100)
	lotteryNum, err := strconv.Atoi(param)
	if err != nil || lotteryNum < 1 || lotteryNum > 100 {
		sendMessage := tgbotapi.NewMessage(chatId, "输入的数字不合法，请输入 1-100 之间的整数！")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	// 检查用户是否已经参与过本次抽奖
	var detail utils.GroupLotteryDetail
	bot.DBEngine.Raw("select * from group_lottery_detail where lottery_id = ? and user_number = ?", lottery.Id, userId).Scan(&detail)
	if detail.Id != "" {
		sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("您已参加过本次抽奖，选择的数字是：%d", detail.LotteryNumber))
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	// 检查数字是否已被其他用户选中
	var otherDetail utils.GroupLotteryDetail
	utils.GetLotteryDetail(lottery.Id, lotteryNum).Scan(&otherDetail)
	if otherDetail.Id != "" {
		sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("数字 %d 已被其他用户选择，请尝试其他数字！", lotteryNum))
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	// 保存抽奖详情
	id, _ := gonanoid.New(32)
	groupLotteryDetail := utils.GroupLotteryDetail{
		Id:            id,
		LotteryId:     lottery.Id,
		UserName:      update.Message.From.FullName(),
		UserNumber:    userId,
		LotteryNumber: int64(lotteryNum),
		Status:        0,
	}
	bot.DBEngine.Table("group_lottery_detail").Create(&groupLotteryDetail)

	sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("参与成功！您选择的数字是：%d", lotteryNum))
	sendMessage.ReplyToMessageID = messageId
	msg, err := bot.Arknights.Send(sendMessage)
	if err != nil {
		return err
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
	return nil
}

// LotteryDetailHandle 查看参加详情
func LotteryDetailHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 5)

	// 检查是否存在已开启的抽奖
	var lottery utils.GroupLottery
	utils.GetGroupLottery(chatId).Scan(&lottery)
	if lottery.Id == "" {
		sendMessage := tgbotapi.NewMessage(chatId, "当前群组暂无正在进行的抽奖活动！")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	bot.Arknights.Send(sendAction)

	port := viper.GetString("http.port")
	url := fmt.Sprintf("http://localhost:%s/lotteryDetail?lotteryId=%s", port, lottery.Id)
	pic, err := utils.Screenshot(url, 0, 1.5)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, "生成详情图片失败，请稍后再试！")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	sendPhoto.ReplyToMessageID = messageId
	_, err = bot.Arknights.Send(sendPhoto)
	if err != nil {
		return err
	}

	return nil
}

// LotteryHandle 进行抽奖
func LotteryHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 5)

	if !bot.Arknights.IsAdminWithPermissions(chatId, userId, tgbotapi.AdminCanRestrictMembers) {
		sendMessage := tgbotapi.NewMessage(chatId, "无使用权限！")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	// 检查当前群组是否存在已开启的抽奖
	var lottery utils.GroupLottery
	utils.GetGroupLottery(chatId).Scan(&lottery)
	if lottery.Id == "" {
		sendMessage := tgbotapi.NewMessage(chatId, "当前群组暂无正在进行的抽奖活动！")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	// 开始随机抽奖
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var winningWinner *utils.GroupLotteryDetail
	var results []string

	for i := 1; i <= 5; i++ {
		luckyNum := r.Intn(100) + 1 // 1-100
		var winner utils.GroupLotteryDetail
		utils.GetLotteryDetail(lottery.Id, luckyNum).Scan(&winner)

		if winner.Id != "" && winner.Status == 0 && winningWinner == nil {
			// 找到第一个未中奖用户，设为本轮唯一中奖者
			winner.Status = 1
			bot.DBEngine.Table("group_lottery_detail").Save(&winner)
			winningWinner = &winner
			results = append(results, fmt.Sprintf("第 %d 个号码：%d — 中奖", i, luckyNum))
		} else if winner.Id != "" {
			// 有人选但已有中奖者产生或该用户已在中奖名单中
			statusText := "已选出中奖者"
			if winner.Status == 1 {
				statusText = "已中奖"
			}
			results = append(results, fmt.Sprintf("第 %d 个号码：%d — %s", i, luckyNum, statusText))
		} else {
			// 无人选择
			results = append(results, fmt.Sprintf("第 %d 个号码：%d — 无人选择", i, luckyNum))
		}
	}

	// 发送抽奖报告
	msgText := "*💎 抽奖选号结果：*\n\n" + strings.Join(results, "\n")
	if winningWinner != nil {
		msgText += fmt.Sprintf("\n\n恭喜 [%s](tg://user?id=%d) 成为本次幸运儿！🎉", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, winningWinner.UserName), winningWinner.UserNumber)
	} else {
		msgText += "\n\n很遗憾，本次无人中奖 🌚"
	}

	sendMessage := tgbotapi.NewMessage(chatId, msgText)
	sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
	sendMessage.ReplyToMessageID = messageId
	bot.Arknights.Send(sendMessage)
	return nil
}

// CheckStopLottery 检查抽奖是否停止报名
func CheckStopLottery() {
	var lotteryList []utils.GroupLottery
	utils.GetAllGroupLottery().Scan(&lotteryList)
	for _, lottery := range lotteryList {
		if lottery.EndTime.Before(time.Now()) {
			lottery.Status = 2
			bot.DBEngine.Table("group_lottery").Save(&lottery)
			log.Println("抽奖报名截止时间到达，报名已结束")
		}
	}
}
