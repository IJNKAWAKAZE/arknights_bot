package player

import (
	"strconv"
	"strings"
)

type PlayerOperation int

const (
	OP_STATE  PlayerOperation = iota // 实时数据
	OP_BOX                           // 我的干员
	OP_GACHA                         // 抽卡记录
	OP_CARD                          // 我的名片
	OP_IMPORT                        // 导入抽卡记录
	OP_EXPORT                        // 导出抽卡记录
)

var (
	playerOperationMap = map[string]PlayerOperation{
		"state":  OP_STATE,
		"box":    OP_BOX,
		"gacha":  OP_GACHA,
		"card":   OP_CARD,
		"import": OP_IMPORT,
		"export": OP_EXPORT,
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
	default:
		return "请选择要查询的角色"
	}
}
