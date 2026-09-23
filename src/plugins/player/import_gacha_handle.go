package player

import (
	"arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/commandoperation"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/tidwall/sjson"
	"io"
	"log"
	"strconv"
	"strings"
)

type ImportGachaData struct {
	Data map[string]struct {
		P   string          `json:"p"`
		Pi  string          `json:"pi"`
		C   [][]interface{} `json:"c"`
		Cn  string          `json:"cn"`
		Ci  string          `json:"ci"`
		Pos int             `json:"pos"`
	} `json:"data"`
}

type PlayerOperationImportS1 struct {
	commandoperation.MultiStepOperation
}

func (o PlayerOperationImportS1) Run(uid string, userAccount account.UserAccount, chatId int64, message *tgbotapi.Message) error {
	commandoperation.AddNextStep(chatId, *o.NextStepOperation(uid, userAccount, message.CommandArguments()), "importGacha")
	sent, sendErr := config.Arknights.SendMarkdownV2(chatId, "请将[网站](https://arkgacha.kwer.top/)导出的json文件发送给机器人或使用 /cancel 指令取消操作。")
	if sendErr != nil {
		log.Printf("%v can not be send error : %v", sent, sendErr)
	}
	return nil
}
func (_ PlayerOperationImportS1) NextStepOperation(playerUID string, userAccount account.UserAccount, param string) *commandoperation.NextStepOperation {
	return &commandoperation.NextStepOperation{
		PlayerID:      playerUID,
		Account:       userAccount,
		Param:         param,
		NextOperation: new(PlayerOperationImportS2),
	}
}

type PlayerOperationImportS2 struct {
	commandoperation.OperationAbstract
}

func (o PlayerOperationImportS2) Run(uid string, userAccount account.UserAccount, chatId int64, message *tgbotapi.Message) error {
	var importGachaData ImportGachaData
	var k = *message.Document
	f, _, err := config.Arknights.DownloadFile(k.FileID)
	if err != nil {
		config.Arknights.SendText(chatId, "下载文件失败！")
		return err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	j, _ := sjson.SetRaw("{}", "data", string(data))
	err = json.Unmarshal([]byte(j), &importGachaData)
	if err != nil {
		config.Arknights.SendText(chatId, "解析抽卡记录失败！")
		return err
	}

	if err := addGacha(importGachaData, userAccount.UserNumber, uid, message.From.FullName()); err != nil {
		config.Arknights.SendText(chatId, "抽卡记录导入失败，请稍后重试。")
		return err
	}
	config.Arknights.SendText(chatId, "抽卡记录导入成功！")
	return nil
}
func (o PlayerOperationImportS2) CheckRequirementsAndPrepare(update tgbotapi.Update) bool {
	doc := update.Message.Document
	result := doc != nil && strings.HasSuffix(doc.FileName, ".json")
	return result
}
func (operation PlayerOperationImportS2) HintOnRequirementsFailed() (string, bool) {
	return "导入文件格式错误！", false
}

// 在同一事务中导入全部记录，避免十连记录仅写入一部分后失败，
// 重试时因已有相同时间的记录而跳过剩余数据。
func addGacha(data ImportGachaData, userNumber int64, uid, name string) error {
	return config.DBEngine.Transaction(func(tx *gorm.DB) error {
		for k, d := range data.Data {
			key, err := strconv.ParseFloat(k, 64)
			if err != nil {
				return fmt.Errorf("抽卡时间格式错误: %w", err)
			}
			key *= 1000
			var existing UserGacha
			res := tx.Raw("select * from user_gacha where user_number = ? and uid = ? and ts = ?", userNumber, uid, key).Scan(&existing)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected > 0 {
				continue
			}
			for _, c := range d.C {
				if len(c) < 3 {
					return fmt.Errorf("抽卡记录字段不完整")
				}
				charName, nameOK := c[0].(string)
				rarity, rarityOK := c[1].(float64)
				flag, flagOK := c[2].(float64)
				if !nameOK || !rarityOK || !flagOK {
					return fmt.Errorf("抽卡记录字段类型错误")
				}
				id, err := gonanoid.New(32)
				if err != nil {
					return err
				}
				record := UserGacha{Id: id, UserName: name, UserNumber: userNumber, Uid: uid, PoolName: d.P, PoolOrder: int(flag), CharName: charName, IsNew: flag == 1, Rarity: int64(rarity), Ts: int64(key)}
				if err := tx.Table("user_gacha").Create(&record).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
