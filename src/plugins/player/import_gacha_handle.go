package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/commandoperation"
	"arknights_bot/utils"
	"encoding/json"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/tidwall/sjson"
	"io"
	"strconv"
	"strings"
)

type ImportGachaData struct {
	Data map[string]struct {
		C [][]interface{} `json:"c"`
		P string          `json:"p"`
	} `json:"data"`
}

type PlayerOperationImportS1 struct {
	commandoperation.MultiStepOperation
}

func (o PlayerOperationImportS1) Run(uid string, userAccount account.UserAccount, chatId int64, message *tgbotapi.Message) error {
	utils.SendMessage(chatId, "请将[网站](https://arkgacha.kwer.top/)导出的json文件发送给机器人或使用 /cancel 指令取消操作。", true, nil)
	commandoperation.AddNextStep(chatId, *o.NextStepOperation(uid, userAccount, message.CommandArguments()), "importGacha")
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
	f, _ := utils.DownloadFile(k.FileID)
	data, _ := io.ReadAll(f)
	j, _ := sjson.SetRaw("{}", "data", string(data))
	err := json.Unmarshal([]byte(j), &importGachaData)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, "解析抽卡记录失败！")
		bot.Arknights.Send(sendMessage)
		return err
	}

	go addGacha(importGachaData, userAccount.UserNumber, uid, message.From.FullName())
	sendMessage := tgbotapi.NewMessage(chatId, "抽卡记录导入成功！")
	bot.Arknights.Send(sendMessage)
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

func addGacha(importGachaData ImportGachaData, userNumber int64, uid string, name string) {
	// 遍历抽卡记录
	for k, d := range importGachaData.Data {
		key, _ := strconv.ParseInt(k, 10, 64)
		var gacha UserGacha
		res := bot.DBEngine.Raw("select * from user_gacha where user_number = ? and uid = ? and ts = ?", userNumber, uid, key).Scan(&gacha)
		if res.RowsAffected == 0 {
			// 同步抽卡数据
			for i, c := range d.C {
				id, _ := gonanoid.New(32)
				n := strconv.Itoa(int(c[2].(float64)))
				isNew, _ := strconv.ParseBool(n)
				userGacha := UserGacha{
					Id:         id,
					UserName:   name,
					UserNumber: userNumber,
					Uid:        uid,
					PoolName:   d.P,
					PoolOrder:  i + 1,
					CharName:   c[0].(string),
					IsNew:      isNew,
					Rarity:     int64(c[1].(float64)),
					Ts:         key,
				}
				bot.DBEngine.Table("user_gacha").Create(&userGacha)
			}
		}
	}
}
