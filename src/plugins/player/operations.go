package player

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type PlayerOperation int

const (
	OP_STATE   PlayerOperation = iota // 实时数据
	OP_BOX                            // 我的干员
	OP_GACHA                          // 抽卡记录
	OP_CARD                           // 我的名片
	OP_IMPORT                         // 导入抽卡记录
	OP_EXPORT                         // 导出抽卡记录
	OP_MISSING                        // 为获取干员
	OP_BASE                           // 基建信息
	OP_SYNC                           // 同步抽卡记录
	OP_REDEEM                         // CDK兑换
)

var (
	playerOperationMap = map[string]PlayerOperation{
		"state":      OP_STATE,
		"box":        OP_BOX,
		"gacha":      OP_GACHA,
		"card":       OP_CARD,
		"import":     OP_IMPORT,
		"export":     OP_EXPORT,
		"missing":    OP_MISSING,
		"base":       OP_BASE,
		"sync_gacha": OP_SYNC,
		"redeem":     OP_REDEEM,
	}
)

func parseIntStringToOperation(str string) (PlayerOperation, bool) {
	result, err := strconv.Atoi(str)
	if err != nil || result < 0 || result >= len(playerOperationMap) {
		return -1, false
	}
	return PlayerOperation(result), true
}
func parseStringToOperation(str string) (PlayerOperation, bool) {
	c, ok := playerOperationMap[strings.ToLower(str)]
	return c, ok
}
func (operation PlayerOperation) getHintWordForPlayerSelection() string {
	switch operation {
	case OP_EXPORT:
		return "请选择要导出的角色"
	case OP_REDEEM:
		return "请选择要兑换的角色"
	default:
		return "请选择要查询的角色"
	}
}
func (operation PlayerOperation) getPerReqForPlayerSelection() func(update tgbotapi.Update) bool {
	switch operation {
	case OP_REDEEM:
		return getRedeemPerFeq
	default:
		return NO_REQUIREMENT
	}
}
