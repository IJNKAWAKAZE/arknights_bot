package system

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"log"
	"strings"
)

var tagMap = make(map[string]string)

func init() {
	tagMap["高级资深干员"] = "高资"
	tagMap["资深干员"] = "资深"
	tagMap["新手"] = "新手"
	tagMap["近卫干员"] = "近卫干员"
	tagMap["狙击干员"] = "狙击干员"
	tagMap["重装干员"] = "重装干员"
	tagMap["医疗干员"] = "医疗干员"
	tagMap["辅助干员"] = "辅助干员"
	tagMap["术师干员"] = "术师干员"
	tagMap["特种干员"] = "特种干员"
	tagMap["先锋干员"] = "先锋干员"
	tagMap["近战位"] = "近战位"
	tagMap["远程位"] = "远程位"
	tagMap["支援机械"] = "机械"
	tagMap["控场"] = "控场"
	tagMap["爆发"] = "爆发"
	tagMap["治疗"] = "治疗"
	tagMap["支援"] = "支援"
	tagMap["费用回复"] = "费用回复"
	tagMap["發用回复"] = "费用回复"
	tagMap["输出"] = "输出"
	tagMap["生存"] = "生存"
	tagMap["群攻"] = "群攻"
	tagMap["防护"] = "防护"
	tagMap["减速"] = "减速"
	tagMap["削弱"] = "削弱"
	tagMap["快速复活"] = "快速复活"
	tagMap["位移"] = "位移"
	tagMap["召唤"] = "召唤"

	tagMap["上級エリート"] = "高资"
	tagMap["上級工リート"] = "高资"
	tagMap["エリート"] = "资深"
	tagMap["工リート"] = "资深"
	tagMap["初期"] = "新手"
	tagMap["前衛タイプ"] = "近卫干员"
	tagMap["狙撃タイプ"] = "狙击干员"
	tagMap["重装タイプ"] = "重装干员"
	tagMap["医療タイプ"] = "医疗干员"
	tagMap["補助タイプ"] = "辅助干员"
	tagMap["術師タイプ"] = "术师干员"
	tagMap["特殊タイプ"] = "特种干员"
	tagMap["先鋒タイプ"] = "先锋干员"
	tagMap["近距離"] = "近战位"
	tagMap["遠距離"] = "远程位"
	tagMap["ロボット"] = "机械"
	tagMap["ロポット"] = "机械"
	tagMap["牽制"] = "控场"
	tagMap["爆発力"] = "爆发"
	tagMap["治療"] = "治疗"
	tagMap["COST回復"] = "费用回复"
	tagMap["火力"] = "输出"
	tagMap["範囲攻撃"] = "群攻"
	tagMap["防御"] = "防护"
	tagMap["減速"] = "减速"
	tagMap["弱化"] = "削弱"
	tagMap["高速再配置"] = "快速复活"
	tagMap["強制移動"] = "位移"
	tagMap["召喚"] = "召唤"
}

func RecruitHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	caption := strings.Split(update.Message.Caption, " ")
	param := caption[len(caption)-1]
	var tags []string
	photos := update.Message.Photo
	file, _ := utils.DownloadFile(photos[len(photos)-1].FileID)
	lang, engine, sep := "chs", "2", "\n"
	if param == "jp" {
		lang, engine, sep = "jpn", "1", "\r\n"
	}

	results := utils.OCR(file, lang, engine, sep)
	if results == nil {
		log.Println("图片识别失败")
		return nil
	}
	for _, result := range results {
		if tagMap[result] != "" {
			tags = append(tags, tagMap[result])
		}
	}
	if len(tags) != 5 {
		sendMessage := tgbotapi.NewMessage(chatId, "标签数量错误，请更换图片。")
		sendMessage.ReplyToMessageID = messageId
		_, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		return nil
	}
	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	bot.Arknights.Send(sendAction)

	port := viper.GetString("http.port")
	pic := utils.Screenshot(fmt.Sprintf("http://localhost:%s/recruit?tags=%s&client=%s", port, strings.Join(tags, " "), param), 0, 1.5)
	if pic == nil {
		sendMessage := tgbotapi.NewMessage(chatId, "生成图片失败，请重试。")
		sendMessage.ReplyToMessageID = messageId
		_, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		return nil
	}
	sendDocument := tgbotapi.NewDocument(chatId, tgbotapi.FileBytes{Bytes: pic, Name: "recruit.jpg"})
	sendDocument.ReplyToMessageID = messageId
	_, err := bot.Arknights.Send(sendDocument)
	if err != nil {
		return err
	}
	return nil
}
