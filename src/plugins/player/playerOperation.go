package player

import (
	"arknights_bot/plugins/commandOperation"
)

var (
	playerOperationMap = map[string]commandOperation.OperationI{
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

func initFactoey() {
	for k, f := range playerOperationMap {
		commandOperation.OperationTypeMaps[k] = f
	}
}
func playerOperationFactory(operation string) *commandOperation.OperationI {

	result, ok := playerOperationMap[operation]
	if !ok {
		return nil
	} else {
		return &result
	}
}
