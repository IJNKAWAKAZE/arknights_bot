package cron

import "arknights_bot/utils"

// UpdateDataSource 更新数据源
func UpdateDataSource() func() {
	updateDataSource := func() {
		go utils.UpdateDataSource()
	}
	return updateDataSource
}
