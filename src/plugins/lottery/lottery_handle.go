package lottery

import (
	"arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils/media"
	"arknights_bot/utils/model"
	"arknights_bot/utils/repo"
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

	if !config.Arknights.IsAdminWithPermissions(chatId, userId, tgbotapi.AdminCanRestrictMembers) {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "无使用权限！")
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}
	// 检查参数是否为 24 小时制时间的年月日时分秒格式
	var endTime time.Time
	if param != "" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", param, time.Local)
		if err != nil {
			msg, err := config.Arknights.ReplyText(chatId, messageId, "参数格式错误！请输入 YYYY-MM-DD HH:MM:SS 格式")
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
			return nil
		}
		endTime = t
	}

	// 如果没有设置报名截止时间（参数为空）
	if endTime.IsZero() {
		endTime = time.Now().Add(time.Hour * 24 * 7) // 默认 7 天后截止报名
	}
	// 检查是否存在已开启的抽奖
	var lottery model.GroupLottery
	repo.GetGroupLottery(chatId).Scan(&lottery)
	if lottery.Id != "" {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "已有正在进行的抽奖活动，请先结束当前抽奖！")
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}

	// 保存抽奖记录
	id, _ := gonanoid.New(32)
	groupLottery := model.GroupLottery{
		Id:          id,
		GroupName:   message.Chat.Title,
		GroupNumber: message.Chat.ID,
		Status:      1,
		EndTime:     endTime,
	}
	res := config.DBEngine.Table("group_lottery").Create(&groupLottery)
	log.Println(res.Error)
	config.Arknights.SendMarkdownV2(chatId, tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, fmt.Sprintf("🎉 *抽奖活动已开启*\n\n📅 *报名截止时间*：%s\n\n📝 *指令说明*：\n🔹 参与选号：`/join_lottery [1-100]`\n🔹 查看详情：`/lottery_detail`\n\n⚙️ *管理指令*：\n🔸 停止报名：`/stop_lottery`\n🔸 进行抽奖：`/lottery`\n🔸 结束抽奖：`/end_lottery`", endTime.Format("2006-01-02 15:04:05"))), messageId)
	return nil
}

// StopLotteryHandle 停止抽奖报名
func StopLotteryHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 5)
	if !config.Arknights.IsAdminWithPermissions(chatId, userId, tgbotapi.AdminCanRestrictMembers) {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "无使用权限！")
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}
	var lottery model.GroupLottery
	repo.GetGroupLottery(chatId).Scan(&lottery)
	if lottery.Id == "" {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "当前群组暂无正在进行的抽奖活动！")
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}
	if lottery.Status == 2 {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "抽奖活动已停止报名，请勿重复操作！")
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}
	lottery.Status = 2
	config.DBEngine.Table("group_lottery").Save(&lottery)
	msg, err := config.Arknights.ReplyText(chatId, messageId, "抽奖活动已停止报名！")
	if err != nil {
		return err
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
	return nil
}

// EndLotteryHandle 结束抽奖活动
func EndLotteryHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 5)
	if !config.Arknights.IsAdminWithPermissions(chatId, userId, tgbotapi.AdminCanRestrictMembers) {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "无使用权限！")
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}
	var lottery model.GroupLottery
	repo.GetGroupLottery(chatId).Scan(&lottery)
	if lottery.Id == "" {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "当前群组暂无抽奖活动！")
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}
	lottery.Status = 0
	config.DBEngine.Table("group_lottery").Save(&lottery)
	msg, err := config.Arknights.ReplyText(chatId, messageId, "抽奖活动已结束！")
	if err != nil {
		return err
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
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
	var lottery model.GroupLottery
	repo.GetGroupLottery(chatId).Scan(&lottery)
	if lottery.Id == "" {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "当前群组暂无正在进行的抽奖活动！")
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}

	if lottery.Status == 2 {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "抽奖活动已停止报名！")
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}

	// 检查用户输入的数字是否合法 (1-100)
	lotteryNum, err := strconv.Atoi(param)
	if err != nil || lotteryNum < 1 || lotteryNum > 100 {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "输入的数字不合法，请输入 1-100 之间的整数！")
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}

	// 检查用户是否已经参与过本次抽奖
	var detail model.GroupLotteryDetail
	config.DBEngine.Raw("select * from group_lottery_detail where lottery_id = ? and user_number = ?", lottery.Id, userId).Scan(&detail)
	if detail.Id != "" {
		msg, err := config.Arknights.ReplyText(chatId, messageId, fmt.Sprintf("您已参加过本次抽奖，选择的数字是：%d", detail.LotteryNumber))
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}

	// 检查数字是否已被其他用户选中
	var otherDetail model.GroupLotteryDetail
	repo.GetLotteryDetail(lottery.Id, lotteryNum).Scan(&otherDetail)
	if otherDetail.Id != "" {
		msg, err := config.Arknights.ReplyText(chatId, messageId, fmt.Sprintf("数字 %d 已被其他用户选择，请尝试其他数字！", lotteryNum))
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}

	// 保存抽奖详情
	id, _ := gonanoid.New(32)
	groupLotteryDetail := model.GroupLotteryDetail{
		Id:            id,
		LotteryId:     lottery.Id,
		UserName:      update.Message.From.FullName(),
		UserNumber:    userId,
		LotteryNumber: int64(lotteryNum),
		Status:        0,
	}
	config.DBEngine.Table("group_lottery_detail").Create(&groupLotteryDetail)

	msg, err := config.Arknights.ReplyText(chatId, messageId, fmt.Sprintf("参与成功！您选择的数字是：%d", lotteryNum))
	if err != nil {
		return err
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
	return nil
}

// LotteryDetailHandle 查看参加详情
func LotteryDetailHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 5)

	// 检查是否存在已开启的抽奖
	var lottery model.GroupLottery
	repo.GetGroupLottery(chatId).Scan(&lottery)
	if lottery.Id == "" {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "当前群组暂无正在进行的抽奖活动！")
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}

	_, _ = config.Arknights.SendChatAction(chatId, "upload_photo")

	port := viper.GetString("http.port")
	url := fmt.Sprintf("http://localhost:%s/lotteryDetail?lotteryId=%s", port, lottery.Id)
	pic, err := media.Screenshot(url, 0, 1.5)
	if err != nil {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "生成详情图片失败，请稍后再试！")
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}

	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	sendPhoto.ReplyToMessageID = messageId
	msg, err := config.Arknights.Send(sendPhoto)
	if err != nil {
		return err
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, 60)
	return nil
}

// LotteryHandle 进行抽奖
func LotteryHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 5)

	if !config.Arknights.IsAdminWithPermissions(chatId, userId, tgbotapi.AdminCanRestrictMembers) {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "无使用权限！")
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}

	// 检查当前群组是否存在已开启的抽奖
	var lottery model.GroupLottery
	repo.GetGroupLottery(chatId).Scan(&lottery)
	if lottery.Id == "" {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "当前群组暂无正在进行的抽奖活动！")
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}

	// 开始随机抽奖
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var winningWinner *model.GroupLotteryDetail
	var results []string

	for i := 1; i <= 5; i++ {
		luckyNum := r.Intn(100) + 1 // 1-100
		var winner model.GroupLotteryDetail
		repo.GetLotteryDetail(lottery.Id, luckyNum).Scan(&winner)

		if winner.Id != "" && winner.Status == 0 && winningWinner == nil {
			// 找到第一个未中奖用户，设为本轮唯一中奖者
			winner.Status = 1
			config.DBEngine.Table("group_lottery_detail").Save(&winner)
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

	config.Arknights.SendMarkdownV2(chatId, msgText, messageId)
	return nil
}

// CheckStopLottery 检查抽奖是否停止报名
func CheckStopLottery() {
	var lotteryList []model.GroupLottery
	repo.GetAllGroupLottery().Scan(&lotteryList)
	for _, lottery := range lotteryList {
		if lottery.EndTime.Before(time.Now()) {
			lottery.Status = 2
			config.DBEngine.Table("group_lottery").Save(&lottery)
			log.Println("抽奖报名截止时间到达，报名已结束")
		}
	}
}
