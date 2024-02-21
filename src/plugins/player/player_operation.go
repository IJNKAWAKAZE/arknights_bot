package player

import (
	"arknights_bot/plugins/commandoperation"
)

var (
	playerOperationMap = map[string]commandoperation.OperationI{
		"state":        PlayerOperationState{},
		"box":          PlayerOperationBox{},
		"gacha":        PlayerOperationGacha{},
		"card":         PlayerOperationCard{},
		"import_gacha": PlayerOperationImportS1{},
		"export_gacha": PlayerOperationExport{},
		"missing":      PlayerOperationMissing{},
		"base":         PlayerOperationBase{},
		"sync_gacha":   PlayerOperationSyncGacha{},
		"redeem":       PlayerOperationRedeem{},
	}
)

func initFactory() {
	for k, f := range playerOperationMap {
		commandoperation.OperationTypeMaps[k] = f
	}
}
func playerOperationFactory(operation string) *commandoperation.OperationI {
	result, ok := playerOperationMap[operation]
	if !ok {
		return nil
	} else {
		return &result
	}
}
