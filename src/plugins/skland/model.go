package skland

// Account 账号信息
type Account struct {
	UserId     string            `json:"userId"`
	Hypergryph AccountHypergryph `json:"hypergryph"`
	Skland     AccountSkland     `json:"skland"`
}

type AccountHypergryph struct {
	Token string `json:"token"`
	Code  string `json:"code"`
}

type AccountSkland struct {
	Cred  string `json:"cred"`
	Token string `json:"token"`
}

type GrantAppData struct {
	Uid  string `json:"uid"`
	Code string `json:"code"`
}

type GenCredByCodeData struct {
	UserId string `json:"userId"`
	Cred   string `json:"cred"`
	Token  string `json:"token"`
}

type AuthRefreshData struct {
	Token string `json:"token"`
}

type ListPlayerData struct {
	List []*PlayersByApp `json:"list"`
}

type PlayersByApp struct {
	AppCode     string    `json:"appCode"`
	AppName     string    `json:"appName"`
	DefaultUid  string    `json:"defaultUid"`
	BindingList []*Player `json:"bindingList"`
}

type Player struct {
	Uid             string `json:"uid"`
	ChannelName     string `json:"channelName"`
	ChannelMasterId string `json:"channelMasterId"`
	NickName        string `json:"nickName"`
	IsOfficial      bool   `json:"isOfficial"`
	IsDefault       bool   `json:"isDefault"`
	IsDelete        bool   `json:"isDelete"`
}

type User struct {
	HgId string `json:"hgId"`
}

// PlayerStatistic 玩家概况
type PlayerStatistic struct {
	CurrentTs  string `json:"currentTs"`
	PlayerName string `json:"playerName"`
	Avatar     string `json:"avatar"`
	Ap         struct {
		Current   int    `json:"current"`
		Max       int    `json:"max"`
		RecoverTs string `json:"recoverTs"`
	} `json:"ap"`
	CheckedIn  bool `json:"checkedIn"`
	TowerLower struct {
		Current   int    `json:"current"`
		Max       int    `json:"max"`
		RecoverTs string `json:"recoverTs"`
	} `json:"towerLower"`
	TowerHigher struct {
		Current   int    `json:"current"`
		Max       int    `json:"max"`
		RecoverTs string `json:"recoverTs"`
	} `json:"towerHigher"`
	Reward struct {
		Current   int    `json:"current"`
		Max       int    `json:"max"`
		RecoverTs string `json:"recoverTs"`
	} `json:"reward"`
	Recruitment struct {
		Current int `json:"current"`
		Max     int `json:"max"`
	} `json:"recruitment"`
	Trading struct {
		Current int `json:"current"`
		Max     int `json:"max"`
	} `json:"trading"`
	Manufacture struct {
		Current int `json:"current"`
		Max     int `json:"max"`
	} `json:"manufacture"`
	TiredChars int `json:"tiredChars"`
	Training   struct {
		CharIcon    string `json:"charIcon"`
		LeftSeconds string `json:"leftSeconds"`
	} `json:"training"`
	BgURL string `json:"bgURL"`
}

// PlayerCards 玩家名片
type PlayerCards struct {
	List []struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Icon    string `json:"icon"`
		BgURL   string `json:"bgUrl"`
		Privacy struct {
			CardOn         bool `json:"cardOn"`
			DetailOn       bool `json:"detailOn"`
			GameRelationOn bool `json:"gameRelationOn"`
		} `json:"privacy"`
		Link      string `json:"link"`
		Arknights struct {
			UID    string `json:"uid"`
			Name   string `json:"name"`
			Level  int    `json:"level"`
			Avatar struct {
				Type string `json:"type"`
				ID   string `json:"id"`
				URL  string `json:"url"`
			} `json:"avatar"`
			RegisterTs        int    `json:"registerTs"`
			MainStageProgress string `json:"mainStageProgress"`
			Secretary         struct {
				CharID string `json:"charId"`
				SkinID string `json:"skinId"`
			} `json:"secretary"`
			Resume          string `json:"resume"`
			SubscriptionEnd int    `json:"subscriptionEnd"`
			Ap              struct {
				Current              int `json:"current"`
				Max                  int `json:"max"`
				LastApAddTime        int `json:"lastApAddTime"`
				CompleteRecoveryTime int `json:"completeRecoveryTime"`
			} `json:"ap"`
			StoreTs      int `json:"storeTs"`
			LastOnlineTs int `json:"lastOnlineTs"`
			CharCnt      int `json:"charCnt"`
			FurnitureCnt int `json:"furnitureCnt"`
			SkinCnt      int `json:"skinCnt"`
			Exp          struct {
				Current int `json:"current"`
				Max     int `json:"max"`
			} `json:"exp"`
		} `json:"arknights"`
		IconBorderColor string `json:"iconBorderColor"`
		GameChar        string `json:"gameChar"`
		Decoration      struct {
			ID           int    `json:"id"`
			URL          string `json:"url"`
			Kind         int    `json:"kind"`
			ResourceKind int    `json:"resourceKind"`
			TopColor     string `json:"topColor"`
			TextColor    string `json:"textColor"`
		} `json:"decoration"`
	} `json:"list"`
}

// PlayerCultivate 玩家养成数据
type PlayerCultivate struct {
	Characters []struct {
		ID             string `json:"id"`
		Level          int    `json:"level"`
		EvolvePhase    int    `json:"evolvePhase"`
		MainSkillLevel int    `json:"mainSkillLevel"`
		Skills         []struct {
			ID    string `json:"id"`
			Level int    `json:"level"`
		} `json:"skills"`
		Equips []struct {
			ID    string `json:"id"`
			Level int    `json:"level"`
		} `json:"equips"`
		PotentialRank int `json:"potentialRank"`
	} `json:"characters"`
	Items []struct {
		ID    string `json:"id"`
		Count string `json:"count"`
	} `json:"items"`
}

// PoolInfo 卡池信息
type PoolInfo struct {
	PoolName string `json:"poolName"`
	PoolId   string `json:"poolId"`
}

// Char 抽卡记录
type Char struct {
	PoolName  string `json:"poolName"`
	PoolOrder int    `json:"poolOrder"`
	Name      string `json:"name"`
	IsNew     bool   `json:"isNew"`
	Rarity    int64  `json:"rarity"`
	Ts        int64  `json:"ts"`
}

// SignGameData 签到奖励
type SignGameData struct {
	Ts     string         `json:"ts"`
	Awards SignGameAwards `json:"awards"`
}

type SignGameAward struct {
	Type     string       `json:"type"`
	Count    int          `json:"count"`
	Resource *SignGameRes `json:"resource"`
}

type SignGameRes struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Rarity int    `json:"rarity"`
}

type SignGameAwards []*SignGameAward
