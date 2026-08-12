package player

import (
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils/repo"
)

// GetPlayerData 查询账号并获取玩家数据，成功后同步玩家名称
func GetPlayerData(userId int64, sklandId, uid string) (*skland.PlayerData, account.UserAccount, skland.Account, error) {
	var userAccount account.UserAccount
	var skAccount skland.Account
	repo.GetAccountByUserIdAndSklandId(userId, sklandId).Scan(&userAccount)
	skAccount.Hypergryph.Token = userAccount.HypergryphToken
	skAccount.Skland.Token = userAccount.SklandToken
	skAccount.Skland.Cred = userAccount.SklandCred
	playerData, skAccount, err := skland.GetPlayerInfo(uid, skAccount, userAccount.ServerName)
	if err != nil {
		return nil, userAccount, skAccount, err
	}
	repo.UpdatePlayerName(uid, playerData.Status.Name)
	return playerData, userAccount, skAccount, nil
}
