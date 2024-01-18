package skland

type PlayerData struct {
	CurrentTs  int `json:"currentTs"`
	ShowConfig struct {
		CharSwitch      bool `json:"charSwitch"`
		SkinSwitch      bool `json:"skinSwitch"`
		StandingsSwitch bool `json:"standingsSwitch"`
	} `json:"showConfig"`
	Status struct {
		UID    string `json:"uid"`
		Name   string `json:"name"`
		Level  int    `json:"level"`
		Avatar struct {
			Type string `json:"type"`
			Id   string `json:"id"`
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
	} `json:"status"`
	AssistChars []struct {
		Name            string `json:"name"`
		CharID          string `json:"charId"`
		SkinID          string `json:"skinId"`
		Level           int    `json:"level"`
		EvolvePhase     int    `json:"evolvePhase"`
		PotentialRank   int    `json:"potentialRank"`
		SkillID         string `json:"skillId"`
		MainSkillLvl    int    `json:"mainSkillLvl"`
		SpecializeLevel int    `json:"specializeLevel"`
		Equip           struct {
			ID           string `json:"id"`
			Level        int    `json:"level"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"equip"`
	} `json:"assistChars"`
	Chars []struct {
		CharID        string `json:"charId"`
		SkinID        string `json:"skinId"`
		Level         int    `json:"level"`
		EvolvePhase   int    `json:"evolvePhase"`
		PotentialRank int    `json:"potentialRank"`
		MainSkillLvl  int    `json:"mainSkillLvl"`
		Skills        []struct {
			ID              string `json:"id"`
			SpecializeLevel int    `json:"specializeLevel"`
		} `json:"skills"`
		Equip []struct {
			ID    string `json:"id"`
			Level int    `json:"level"`
		} `json:"equip"`
		FavorPercent   int    `json:"favorPercent"`
		DefaultSkillID string `json:"defaultSkillId"`
		GainTime       int    `json:"gainTime"`
		DefaultEquipID string `json:"defaultEquipId"`
	} `json:"chars"`
	Skins []struct {
		ID string `json:"id"`
		Ts int    `json:"ts"`
	} `json:"skins"`
	Building struct {
		TiredChars []struct {
			CharID        string `json:"charId"`
			Ap            int    `json:"ap"`
			LastApAddTime int    `json:"lastApAddTime"`
			RoomSlotID    string `json:"roomSlotId"`
			Index         int    `json:"index"`
			Bubble        struct {
				Normal struct {
					Add int `json:"add"`
					Ts  int `json:"ts"`
				} `json:"normal"`
				Assist struct {
					Add int `json:"add"`
					Ts  int `json:"ts"`
				} `json:"assist"`
			} `json:"bubble"`
			WorkTime int `json:"workTime"`
		} `json:"tiredChars"`
		TiredCharsCount int `json:"tiredCharsCount"`
		Powers          []struct {
			SlotID string `json:"slotId"`
			Level  int    `json:"level"`
			Chars  []struct {
				CharID        string `json:"charId"`
				Ap            int    `json:"ap"`
				LastApAddTime int    `json:"lastApAddTime"`
				Index         int    `json:"index"`
				Bubble        struct {
					Normal struct {
						Add int `json:"add"`
						Ts  int `json:"ts"`
					} `json:"normal"`
					Assist struct {
						Add int `json:"add"`
						Ts  int `json:"ts"`
					} `json:"assist"`
				} `json:"bubble"`
				WorkTime int `json:"workTime"`
			} `json:"chars"`
		} `json:"powers"`
		Manufactures []struct {
			SlotID string `json:"slotId"`
			Level  int    `json:"level"`
			Chars  []struct {
				CharID        string `json:"charId"`
				Ap            int    `json:"ap"`
				LastApAddTime int    `json:"lastApAddTime"`
				Index         int    `json:"index"`
				Bubble        struct {
					Normal struct {
						Add int `json:"add"`
						Ts  int `json:"ts"`
					} `json:"normal"`
					Assist struct {
						Add int `json:"add"`
						Ts  int `json:"ts"`
					} `json:"assist"`
				} `json:"bubble"`
				WorkTime int `json:"workTime"`
			} `json:"chars"`
			CompleteWorkTime int    `json:"completeWorkTime"`
			LastUpdateTime   int    `json:"lastUpdateTime"`
			FormulaID        string `json:"formulaId"`
			Capacity         int    `json:"capacity"`
			Weight           int    `json:"weight"`
			Complete         int    `json:"complete"`
			Remain           int    `json:"remain"`
			Speed            int    `json:"speed"`
		} `json:"manufactures"`
		ManufacturesCurrent int `json:"manufacturesCurrent"`
		ManufacturesTotal   int `json:"manufacturesTotal"`
		Tradings            []struct {
			SlotID string `json:"slotId"`
			Level  int    `json:"level"`
			Chars  []struct {
				Index         int    `json:"index"`
				CharID        string `json:"charId"`
				Ap            int    `json:"ap"`
				LastApAddTime int    `json:"lastApAddTime"`
				Bubble        struct {
					Normal struct {
						Add int `json:"add"`
						Ts  int `json:"ts"`
					} `json:"normal"`
					Assist struct {
						Add int `json:"add"`
						Ts  int `json:"ts"`
					} `json:"assist"`
				} `json:"bubble"`
				WorkTime int `json:"workTime"`
			} `json:"chars"`
			CompleteWorkTime int    `json:"completeWorkTime"`
			LastUpdateTime   int    `json:"lastUpdateTime"`
			Strategy         string `json:"strategy"`
			Stock            []struct {
				Delivery []struct {
					Count int    `json:"count"`
					ID    string `json:"id"`
					Type  string `json:"type"`
				} `json:"delivery"`
				Gain struct {
					Count int    `json:"count"`
					ID    string `json:"id"`
					Type  string `json:"type"`
				} `json:"gain"`
				InstID     int    `json:"instId"`
				IsViolated bool   `json:"isViolated"`
				Type       string `json:"type"`
			} `json:"stock"`
			StockLimit int `json:"stockLimit"`
		} `json:"tradings"`
		TradingsCurrent int `json:"tradingsCurrent"`
		TradingsTotal   int `json:"tradingsTotal"`
		Dormitories     []struct {
			SlotID string `json:"slotId"`
			Level  int    `json:"level"`
			Chars  []struct {
				Index         int    `json:"index"`
				CharID        string `json:"charId"`
				Ap            int    `json:"ap"`
				LastApAddTime int    `json:"lastApAddTime"`
				Bubble        struct {
					Normal struct {
						Add int `json:"add"`
						Ts  int `json:"ts"`
					} `json:"normal"`
					Assist struct {
						Add int `json:"add"`
						Ts  int `json:"ts"`
					} `json:"assist"`
				} `json:"bubble"`
				WorkTime int `json:"workTime"`
			} `json:"chars"`
			Comfort int `json:"comfort"`
		} `json:"dormitories"`
		Meeting struct {
			SlotID string `json:"slotId"`
			Level  int    `json:"level"`
			Chars  []struct {
				Index         int    `json:"index"`
				CharID        string `json:"charId"`
				Ap            int    `json:"ap"`
				LastApAddTime int    `json:"lastApAddTime"`
				Bubble        struct {
					Normal struct {
						Add int `json:"add"`
						Ts  int `json:"ts"`
					} `json:"normal"`
					Assist struct {
						Add int `json:"add"`
						Ts  int `json:"ts"`
					} `json:"assist"`
				} `json:"bubble"`
				WorkTime int `json:"workTime"`
			} `json:"chars"`
			Clue struct {
				Own               int      `json:"own"`
				Received          int      `json:"received"`
				DailyReward       bool     `json:"dailyReward"`
				NeedReceive       int      `json:"needReceive"`
				Board             []string `json:"board"`
				Sharing           bool     `json:"sharing"`
				ShareCompleteTime int      `json:"shareCompleteTime"`
			} `json:"clue"`
			LastUpdateTime   int `json:"lastUpdateTime"`
			CompleteWorkTime int `json:"completeWorkTime"`
		} `json:"meeting"`
		Hire struct {
			SlotID string `json:"slotId"`
			Level  int    `json:"level"`
			Chars  []struct {
				Index         int    `json:"index"`
				CharID        string `json:"charId"`
				Ap            int    `json:"ap"`
				LastApAddTime int    `json:"lastApAddTime"`
				Bubble        struct {
					Normal struct {
						Add int `json:"add"`
						Ts  int `json:"ts"`
					} `json:"normal"`
					Assist struct {
						Add int `json:"add"`
						Ts  int `json:"ts"`
					} `json:"assist"`
				} `json:"bubble"`
				WorkTime int `json:"workTime"`
			} `json:"chars"`
			State            int `json:"state"`
			RefreshCount     int `json:"refreshCount"`
			CompleteWorkTime int `json:"completeWorkTime"`
			SlotState        int `json:"slotState"`
		} `json:"hire"`
		Training struct {
			SlotID  string `json:"slotId"`
			Level   int    `json:"level"`
			Trainee struct {
				Ap            int    `json:"ap"`
				CharID        string `json:"charId"`
				LastApAddTime int    `json:"lastApAddTime"`
				TargetSkill   int    `json:"targetSkill"`
			} `json:"trainee"`
			Trainer struct {
				CharID        string `json:"charId"`
				Ap            int    `json:"ap"`
				LastApAddTime int    `json:"lastApAddTime"`
			} `json:"trainer"`
			RemainPoint    int     `json:"remainPoint"`
			Speed          float64 `json:"speed"`
			LastUpdateTime int     `json:"lastUpdateTime"`
			RemainSecs     int     `json:"remainSecs"`
			SlotState      int     `json:"slotState"`
		} `json:"training"`
		Labor struct {
			MaxValue       int `json:"maxValue"`
			Value          int `json:"value"`
			LastUpdateTime int `json:"lastUpdateTime"`
			RemainSecs     int `json:"remainSecs"`
		} `json:"labor"`
		Furniture struct {
			Total int `json:"total"`
		} `json:"furniture"`
		Elevators []struct {
			SlotID    string `json:"slotId"`
			SlotState int    `json:"slotState"`
			Level     int    `json:"level"`
		} `json:"elevators"`
		Corridors []struct {
			SlotID    string `json:"slotId"`
			SlotState int    `json:"slotState"`
			Level     int    `json:"level"`
		} `json:"corridors"`
		Control struct {
			SlotID    string `json:"slotId"`
			SlotState int    `json:"slotState"`
			Level     int    `json:"level"`
			Chars     []struct {
				CharID        string `json:"charId"`
				Ap            int    `json:"ap"`
				LastApAddTime int    `json:"lastApAddTime"`
				Index         int    `json:"index"`
				Bubble        struct {
					Normal struct {
						Add int `json:"add"`
						Ts  int `json:"ts"`
					} `json:"normal"`
					Assist struct {
						Add int `json:"add"`
						Ts  int `json:"ts"`
					} `json:"assist"`
				} `json:"bubble"`
				WorkTime int `json:"workTime"`
			} `json:"chars"`
		} `json:"control"`
	} `json:"building"`
	Recruit []struct {
		StartTs  int `json:"startTs"`
		FinishTs int `json:"finishTs"`
		State    int `json:"state"`
	} `json:"recruit"`
	RecruitFinished int `json:"recruitFinished"`
	RecruitTotal    int `json:"recruitTotal"`
	Campaign        struct {
		Records []struct {
			CampaignID string `json:"campaignId"`
			MaxKills   int    `json:"maxKills"`
		} `json:"records"`
		Reward struct {
			Current int `json:"current"`
			Total   int `json:"total"`
		} `json:"reward"`
	} `json:"campaign"`
	Tower struct {
		Records []struct {
			TowerID string `json:"towerId"`
			Best    int    `json:"best"`
		} `json:"records"`
		Reward struct {
			HigherItem struct {
				Current int `json:"current"`
				Total   int `json:"total"`
			} `json:"higherItem"`
			LowerItem struct {
				Current int `json:"current"`
				Total   int `json:"total"`
			} `json:"lowerItem"`
			TermTs int `json:"termTs"`
		} `json:"reward"`
	} `json:"tower"`
	Rogue struct {
		Records []struct {
			RogueID  string `json:"rogueId"`
			RelicCnt int    `json:"relicCnt"`
			Bank     struct {
				Current int `json:"current"`
				Record  int `json:"record"`
			} `json:"bank"`
		} `json:"records"`
	} `json:"rogue"`
	Routine struct {
		Daily struct {
			Current int `json:"current"`
			Total   int `json:"total"`
		} `json:"daily"`
		Weekly struct {
			Current int `json:"current"`
			Total   int `json:"total"`
		} `json:"weekly"`
	} `json:"routine"`
	Activity []struct {
		ActID        string `json:"actId"`
		ActReplicaID string `json:"actReplicaId"`
		Zones        []struct {
			ZoneID        string `json:"zoneId"`
			ZoneReplicaID string `json:"zoneReplicaId"`
			ClearedStage  int    `json:"clearedStage"`
			TotalStage    int    `json:"totalStage"`
		} `json:"zones"`
	} `json:"activity"`
	CharInfoMap map[string]struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		NationID        string `json:"nationId"`
		GroupID         string `json:"groupId"`
		DisplayNumber   string `json:"displayNumber"`
		Rarity          int    `json:"rarity"`
		Profession      string `json:"profession"`
		SubProfessionID string `json:"subProfessionId"`
	} `json:"charInfoMap"`
	SkinInfoMap map[string]struct {
		ID           string `json:"id"`
		BrandID      string `json:"brandId"`
		SortID       int    `json:"sortId"`
		DisplayTagID string `json:"displayTagId"`
	} `json:"skinInfoMap"`
	StageInfoMap map[string]struct {
		ID   string `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"stageInfoMap"`
	ActivityInfoMap map[string]struct {
		ID            string `json:"id"`
		Name          string `json:"name"`
		StartTime     int    `json:"startTime"`
		EndTime       int    `json:"endTime"`
		RewardEndTime int    `json:"rewardEndTime"`
		IsReplicate   bool   `json:"isReplicate"`
		Type          string `json:"type"`
	} `json:"activityInfoMap"`
	TowerInfoMap struct {
	} `json:"towerInfoMap"`
	RogueInfoMap map[string]struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Sort int    `json:"sort"`
	} `json:"rogueInfoMap"`
	CampaignInfoMap map[string]struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		CampaignZoneID string `json:"campaignZoneId"`
	} `json:"campaignInfoMap"`
	CampaignZoneInfoMap map[string]struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"campaignZoneInfoMap"`
	EquipmentInfoMap map[string]struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		TypeIcon     string `json:"typeIcon"`
		ShiningColor string `json:"shiningColor"`
	} `json:"equipmentInfoMap"`
	ManufactureFormulaInfoMap map[string]struct {
		ID        string `json:"id"`
		ItemID    string `json:"itemId"`
		Count     int    `json:"count"`
		Weight    int    `json:"weight"`
		CostPoint int    `json:"costPoint"`
	} `json:"manufactureFormulaInfoMap"`
	CharAssets         []interface{} `json:"charAssets"`
	SkinAssets         []interface{} `json:"skinAssets"`
	ActivityBannerList struct {
		List []struct {
			ActivityID string `json:"activityId"`
			ImgURL     string `json:"imgUrl"`
			URL        string `json:"url"`
			StartTs    int    `json:"startTs"`
			EndTs      int    `json:"endTs"`
			OfflineTs  int    `json:"offlineTs"`
		} `json:"list"`
	} `json:"activityBannerList"`
	BossRush []struct {
		ID     string `json:"id"`
		Record struct {
			Played     bool   `json:"played"`
			StageID    string `json:"stageId"`
			Difficulty string `json:"difficulty"`
		} `json:"record"`
	} `json:"bossRush"`
}
