package skland

type PlayerData struct {
	CurrentTs  int `json:"currentTs"`
	ShowConfig struct {
		CharSwitch      bool `json:"charSwitch"`
		SkinSwitch      bool `json:"skinSwitch"`
		StandingsSwitch bool `json:"standingsSwitch"`
	} `json:"showConfig"`
	Status struct {
		UID               string `json:"uid"`
		Name              string `json:"name"`
		Level             int    `json:"level"`
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
		CharID          string `json:"charId"`
		SkinID          string `json:"skinId"`
		Level           int    `json:"level"`
		EvolvePhase     int    `json:"evolvePhase"`
		PotentialRank   int    `json:"potentialRank"`
		SkillID         string `json:"skillId"`
		MainSkillLvl    int    `json:"mainSkillLvl"`
		SpecializeLevel int    `json:"specializeLevel"`
		Equip           struct {
			ID    string `json:"id"`
			Level int    `json:"level"`
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
	CharInfoMap struct {
		Char002Amiya struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_002_amiya"`
		Char00912Fce struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_009_12fce"`
		Char010Chen struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_010_chen"`
		Char017Huang struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_017_huang"`
		Char1013Chen2 struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_1013_chen2"`
		Char1023Ghost2 struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_1023_ghost2"`
		Char1024Hbisc2 struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_1024_hbisc2"`
		Char1028Texas2 struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_1028_texas2"`
		Char102Texas struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_102_texas"`
		Char106Franka struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_106_franka"`
		Char107Liskam struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_107_liskam"`
		Char108Silent struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_108_silent"`
		Char109Fmout struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_109_fmout"`
		Char110Deepcl struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_110_deepcl"`
		Char115Headbr struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_115_headbr"`
		Char117Myrrh struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_117_myrrh"`
		Char118Yuki struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_118_yuki"`
		Char120Hibisc struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_120_hibisc"`
		Char121Lava struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_121_lava"`
		Char122Beagle struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_122_beagle"`
		Char123Fang struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_123_fang"`
		Char124Kroos struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_124_kroos"`
		Char126Shotst struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_126_shotst"`
		Char127Estell struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_127_estell"`
		Char129Bluep struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_129_bluep"`
		Char130Doberm struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_130_doberm"`
		Char133Mm struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_133_mm"`
		Char134Ifrit struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_134_ifrit"`
		Char136Hsguma struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_136_hsguma"`
		Char137Brownb struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_137_brownb"`
		Char140Whitew struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_140_whitew"`
		Char141Nights struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_141_nights"`
		Char143Ghost struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_143_ghost"`
		Char144Red struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_144_red"`
		Char145Prove struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_145_prove"`
		Char147Shining struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_147_shining"`
		Char149Scave struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_149_scave"`
		Char150Snakek struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_150_snakek"`
		Char151Myrtle struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_151_myrtle"`
		Char155Tiger struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_155_tiger"`
		Char166Skfire struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_166_skfire"`
		Char171Bldsk struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_171_bldsk"`
		Char173Slchan struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_173_slchan"`
		Char174Slbell struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_174_slbell"`
		Char179Cgbird struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_179_cgbird"`
		Char180Amgoat struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_180_amgoat"`
		Char181Flower struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_181_flower"`
		Char183Skgoat struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_183_skgoat"`
		Char185Frncat struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_185_frncat"`
		Char187Ccheal struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_187_ccheal"`
		Char188Helage struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_188_helage"`
		Char190Clour struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_190_clour"`
		Char192Falco struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_192_falco"`
		Char193Frostl struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_193_frostl"`
		Char195Glassb struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_195_glassb"`
		Char196Sunbr struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_196_sunbr"`
		Char197Poca struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_197_poca"`
		Char198Blackd struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_198_blackd"`
		Char199Yak struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_199_yak"`
		Char201Moeshd struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_201_moeshd"`
		Char2024Chyue struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_2024_chyue"`
		Char204Platnm struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_204_platnm"`
		Char208Melan struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_208_melan"`
		Char209Ardign struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_209_ardign"`
		Char210Stward struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_210_stward"`
		Char211Adnach struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_211_adnach"`
		Char212Ansel struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_212_ansel"`
		Char213Mostma struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_213_mostma"`
		Char215Mantic struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_215_mantic"`
		Char218Cuttle struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_218_cuttle"`
		Char219Meteo struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_219_meteo"`
		Char220Grani struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_220_grani"`
		Char226Hmau struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_226_hmau"`
		Char230Savage struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_230_savage"`
		Char235Jesica struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_235_jesica"`
		Char236Rope struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_236_rope"`
		Char237Gravel struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_237_gravel"`
		Char240Wyvern struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_240_wyvern"`
		Char241Panda struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_241_panda"`
		Char242Otter struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_242_otter"`
		Char243Waaifu struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_243_waaifu"`
		Char250Phatom struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_250_phatom"`
		Char253Greyy struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_253_greyy"`
		Char254Vodfox struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_254_vodfox"`
		Char258Podego struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_258_podego"`
		Char260Durnar struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_260_durnar"`
		Char263Skadi struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_263_skadi"`
		Char271Spikes struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_271_spikes"`
		Char272Strong struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_272_strong"`
		Char274Astesi struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_274_astesi"`
		Char275Breeze struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_275_breeze"`
		Char277Sqrrel struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_277_sqrrel"`
		Char278Orchid struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_278_orchid"`
		Char279Excu struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_279_excu"`
		Char281Popka struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_281_popka"`
		Char282Catap struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_282_catap"`
		Char283Midn struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_283_midn"`
		Char284Spot struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_284_spot"`
		Char285Medic2 struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_285_medic2"`
		Char286Cast3 struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_286_cast3"`
		Char289Gyuki struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_289_gyuki"`
		Char290Vigna struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_290_vigna"`
		Char291Aglina struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_291_aglina"`
		Char293Thorns struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_293_thorns"`
		Char294Ayer struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_294_ayer"`
		Char298Susuro struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_298_susuro"`
		Char301Cutter struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_301_cutter"`
		Char302Glaze struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_302_glaze"`
		Char308Swire struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_308_swire"`
		Char325Bison struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_325_bison"`
		Char326Glacus struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_326_glacus"`
		Char328Cammou struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_328_cammou"`
		Char332Archet struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_332_archet"`
		Char333Sidero struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_333_sidero"`
		Char337Utage struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_337_utage"`
		Char338Iris struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_338_iris"`
		Char340Shwaz struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_340_shwaz"`
		Char344Beewax struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_344_beewax"`
		Char345Folnic struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_345_folnic"`
		Char346Aosta struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_346_aosta"`
		Char347Jaksel struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_347_jaksel"`
		Char348Ceylon struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_348_ceylon"`
		Char349Chiave struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_349_chiave"`
		Char355Ethan struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_355_ethan"`
		Char356Broca struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_356_broca"`
		Char366Acdrop struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_366_acdrop"`
		Char367Swllow struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_367_swllow"`
		Char369Bena struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_369_bena"`
		Char373Lionhd struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_373_lionhd"`
		Char376Therex struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_376_therex"`
		Char381Bubble struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_381_bubble"`
		Char383Snsant struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_383_snsant"`
		Char385Finlpp struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_385_finlpp"`
		Char388Mint struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_388_mint"`
		Char4000Jnight struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_4000_jnight"`
		Char4004Pudd struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_4004_pudd"`
		Char4009Irene struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_4009_irene"`
		Char4014Lunacu struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_4014_lunacu"`
		Char4019Ncdeer struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_4019_ncdeer"`
		Char401Elysm struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_401_elysm"`
		Char4032Provs struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_4032_provs"`
		Char4039Horn struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_4039_horn"`
		Char4041Chnut struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_4041_chnut"`
		Char4042Lumen struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_4042_lumen"`
		Char4054Malist struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_4054_malist"`
		Char4062Totter struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_4062_totter"`
		Char4063Quartz struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_4063_quartz"`
		Char4065Judge struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_4065_judge"`
		Char4067Lolxh struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_4067_lolxh"`
		Char411Tomimi struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_411_tomimi"`
		Char421Crow struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_421_crow"`
		Char422Aurora struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_422_aurora"`
		Char427Vigil struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_427_vigil"`
		Char437Mizuki struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_437_mizuki"`
		Char440Pinecn struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_440_pinecn"`
		Char452Bstalk struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_452_bstalk"`
		Char455Nothin struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_455_nothin"`
		Char457Blitz struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_457_blitz"`
		Char458Rfrost struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_458_rfrost"`
		Char459Tachak struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_459_tachak"`
		Char469Indigo struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_469_indigo"`
		Char475Akafyu struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_475_akafyu"`
		Char476Blkngt struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_476_blkngt"`
		Char484Robrta struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_484_robrta"`
		Char489Serum struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_489_serum"`
		Char491Humus struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_491_humus"`
		Char493Firwhl struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_493_firwhl"`
		Char496Wildmn struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_496_wildmn"`
		Char497Ctable struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_497_ctable"`
		Char500Noirc struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_500_noirc"`
		Char501Durin struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_501_durin"`
		Char502Nblade struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_502_nblade"`
		Char503Rang struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			NationID        string `json:"nationId"`
			GroupID         string `json:"groupId"`
			DisplayNumber   string `json:"displayNumber"`
			Rarity          int    `json:"rarity"`
			Profession      string `json:"profession"`
			SubProfessionID string `json:"subProfessionId"`
		} `json:"char_503_rang"`
	} `json:"charInfoMap"`
	SkinInfoMap struct {
		Char002Amiya1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_002_amiya#1"`
		Char00912Fce1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_009_12fce#1"`
		Char010Chen1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_010_chen#1"`
		Char017Huang1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_017_huang#1"`
		Char1013Chen21 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_1013_chen2#1"`
		Char1023Ghost21 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_1023_ghost2#1"`
		Char1024Hbisc21 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_1024_hbisc2#1"`
		Char1028Texas22 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_1028_texas2#2"`
		Char102Texas1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_102_texas#1"`
		Char106Franka1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_106_franka#1"`
		Char107Liskam1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_107_liskam#1"`
		Char108Silent1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_108_silent#1"`
		Char109Fmout1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_109_fmout#1"`
		Char110Deepcl2 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_110_deepcl#2"`
		Char115Headbr1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_115_headbr#1"`
		Char117Myrrh1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_117_myrrh#1"`
		Char118Yuki1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_118_yuki#1"`
		Char120Hibisc1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_120_hibisc#1"`
		Char121Lava1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_121_lava#1"`
		Char122Beagle1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_122_beagle#1"`
		Char123Fang1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_123_fang#1"`
		Char123FangWinter1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_123_fang@winter#1"`
		Char124Kroos1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_124_kroos#1"`
		Char126Shotst1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_126_shotst#1"`
		Char127Estell1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_127_estell#1"`
		Char128PlosisEpoque3 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_128_plosis@epoque#3"`
		Char129Bluep1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_129_bluep#1"`
		Char130Doberm1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_130_doberm#1"`
		Char133Mm1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_133_mm#1"`
		Char134Ifrit1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_134_ifrit#1"`
		Char136Hsguma1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_136_hsguma#1"`
		Char137Brownb1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_137_brownb#1"`
		Char140Whitew1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_140_whitew#1"`
		Char141Nights1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_141_nights#1"`
		Char143Ghost1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_143_ghost#1"`
		Char144Red1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_144_red#1"`
		Char145ProveWild5 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_145_prove@wild#5"`
		Char147Shining1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_147_shining#1"`
		Char149Scave1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_149_scave#1"`
		Char150Snakek1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_150_snakek#1"`
		Char151Myrtle1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_151_myrtle#1"`
		Char155Tiger1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_155_tiger#1"`
		Char166Skfire1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_166_skfire#1"`
		Char171Bldsk1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_171_bldsk#1"`
		Char173Slchan1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_173_slchan#1"`
		Char174Slbell1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_174_slbell#1"`
		Char179Cgbird1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_179_cgbird#1"`
		Char180Amgoat1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_180_amgoat#1"`
		Char181Flower1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_181_flower#1"`
		Char183Skgoat1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_183_skgoat#1"`
		Char185Frncat1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_185_frncat#1"`
		Char187Ccheal1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_187_ccheal#1"`
		Char188Helage1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_188_helage#1"`
		Char190Clour1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_190_clour#1"`
		Char192Falco1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_192_falco#1"`
		Char193Frostl1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_193_frostl#1"`
		Char193FrostlBoc4 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_193_frostl@boc#4"`
		Char195Glassb1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_195_glassb#1"`
		Char196Sunbr1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_196_sunbr#1"`
		Char197Poca1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_197_poca#1"`
		Char198Blackd1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_198_blackd#1"`
		Char199Yak1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_199_yak#1"`
		Char201Moeshd1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_201_moeshd#1"`
		Char2024Chyue1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_2024_chyue#1"`
		Char204Platnm1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_204_platnm#1"`
		Char208Melan1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_208_melan#1"`
		Char208MelanEpoque1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_208_melan@epoque#1"`
		Char209Ardign1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_209_ardign#1"`
		Char209ArdignSnow1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_209_ardign@snow#1"`
		Char210Stward1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_210_stward#1"`
		Char210StwardSale6 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_210_stward@sale#6"`
		Char211Adnach1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_211_adnach#1"`
		Char212Ansel1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_212_ansel#1"`
		Char212AnselSummer1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_212_ansel@summer#1"`
		Char213Mostma1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_213_mostma#1"`
		Char215Mantic1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_215_mantic#1"`
		Char218Cuttle1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_218_cuttle#1"`
		Char219Meteo1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_219_meteo#1"`
		Char220Grani1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_220_grani#1"`
		Char226Hmau1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_226_hmau#1"`
		Char230Savage1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_230_savage#1"`
		Char235Jesica1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_235_jesica#1"`
		Char236Rope1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_236_rope#1"`
		Char237Gravel1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_237_gravel#1"`
		Char240Wyvern1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_240_wyvern#1"`
		Char241Panda1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_241_panda#1"`
		Char242Otter1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_242_otter#1"`
		Char243Waaifu1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_243_waaifu#1"`
		Char250Phatom1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_250_phatom#1"`
		Char253Greyy1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_253_greyy#1"`
		Char254Vodfox1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_254_vodfox#1"`
		Char258Podego1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_258_podego#1"`
		Char260Durnar1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_260_durnar#1"`
		Char263Skadi1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_263_skadi#1"`
		Char271Spikes1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_271_spikes#1"`
		Char272Strong1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_272_strong#1"`
		Char274Astesi1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_274_astesi#1"`
		Char275Breeze1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_275_breeze#1"`
		Char277Sqrrel1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_277_sqrrel#1"`
		Char278Orchid1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_278_orchid#1"`
		Char279Excu1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_279_excu#1"`
		Char281Popka1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_281_popka#1"`
		Char282Catap1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_282_catap#1"`
		Char283Midn1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_283_midn#1"`
		Char284Spot1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_284_spot#1"`
		Char285Medic21 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_285_medic2#1"`
		Char286Cast31 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_286_cast3#1"`
		Char289Gyuki1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_289_gyuki#1"`
		Char290Vigna1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_290_vigna#1"`
		Char291Aglina1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_291_aglina#1"`
		Char293Thorns1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_293_thorns#1"`
		Char294Ayer1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_294_ayer#1"`
		Char298Susuro1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_298_susuro#1"`
		Char301Cutter1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_301_cutter#1"`
		Char302Glaze1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_302_glaze#1"`
		Char308Swire1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_308_swire#1"`
		Char325Bison1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_325_bison#1"`
		Char326Glacus1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_326_glacus#1"`
		Char328Cammou1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_328_cammou#1"`
		Char332Archet1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_332_archet#1"`
		Char333Sidero1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_333_sidero#1"`
		Char337Utage1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_337_utage#1"`
		Char337UtageSummer4 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_337_utage@summer#4"`
		Char338Iris1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_338_iris#1"`
		Char340ShwazStriker1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_340_shwaz@striker#1"`
		Char344Beewax1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_344_beewax#1"`
		Char345Folnic1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_345_folnic#1"`
		Char346Aosta1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_346_aosta#1"`
		Char347Jaksel1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_347_jaksel#1"`
		Char347JakselWhirlwind2 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_347_jaksel@whirlwind#2"`
		Char348Ceylon1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_348_ceylon#1"`
		Char349Chiave1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_349_chiave#1"`
		Char355Ethan1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_355_ethan#1"`
		Char356Broca1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_356_broca#1"`
		Char366Acdrop1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_366_acdrop#1"`
		Char367SwllowBoc1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_367_swllow@boc#1"`
		Char369Bena1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_369_bena#1"`
		Char373Lionhd1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_373_lionhd#1"`
		Char376Therex1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_376_therex#1"`
		Char381Bubble1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_381_bubble#1"`
		Char383Snsant1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_383_snsant#1"`
		Char385Finlpp1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_385_finlpp#1"`
		Char388Mint1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_388_mint#1"`
		Char4000Jnight1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_4000_jnight#1"`
		Char4004Pudd1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_4004_pudd#1"`
		Char4009Irene1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_4009_irene#1"`
		Char4014Lunacu1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_4014_lunacu#1"`
		Char4019Ncdeer1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_4019_ncdeer#1"`
		Char401Elysm1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_401_elysm#1"`
		Char4032Provs1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_4032_provs#1"`
		Char4039Horn1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_4039_horn#1"`
		Char4041Chnut1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_4041_chnut#1"`
		Char4042Lumen1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_4042_lumen#1"`
		Char4054Malist1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_4054_malist#1"`
		Char4062Totter1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_4062_totter#1"`
		Char4063Quartz1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_4063_quartz#1"`
		Char4065Judge1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_4065_judge#1"`
		Char4067Lolxh1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_4067_lolxh#1"`
		Char411Tomimi1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_411_tomimi#1"`
		Char421Crow1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_421_crow#1"`
		Char422Aurora1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_422_aurora#1"`
		Char427Vigil1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_427_vigil#1"`
		Char437Mizuki1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_437_mizuki#1"`
		Char440Pinecn1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_440_pinecn#1"`
		Char452Bstalk1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_452_bstalk#1"`
		Char455Nothin1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_455_nothin#1"`
		Char457Blitz1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_457_blitz#1"`
		Char458Rfrost1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_458_rfrost#1"`
		Char459Tachak1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_459_tachak#1"`
		Char469Indigo1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_469_indigo#1"`
		Char475Akafyu1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_475_akafyu#1"`
		Char476Blkngt1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_476_blkngt#1"`
		Char478KiraraGame2 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_478_kirara@game#2"`
		Char484Robrta1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_484_robrta#1"`
		Char489Serum1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_489_serum#1"`
		Char491Humus1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_491_humus#1"`
		Char492QuercuEpoque17 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_492_quercu@epoque#17"`
		Char493Firwhl1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_493_firwhl#1"`
		Char496Wildmn1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_496_wildmn#1"`
		Char497Ctable1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_497_ctable#1"`
		Char500Noirc1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_500_noirc#1"`
		Char501Durin1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_501_durin#1"`
		Char502Nblade1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_502_nblade#1"`
		Char503Rang1 struct {
			ID           string `json:"id"`
			BrandID      string `json:"brandId"`
			SortID       int    `json:"sortId"`
			DisplayTagID string `json:"displayTagId"`
		} `json:"char_503_rang#1"`
	} `json:"skinInfoMap"`
	StageInfoMap struct {
		Main0911 struct {
			ID   string `json:"id"`
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"main_09-11"`
	} `json:"stageInfoMap"`
	ActivityInfoMap struct {
		OneStact struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"1stact"`
		Act10D5 struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act10d5"`
		Act10Mini struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act10mini"`
		Act11D0 struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act11d0"`
		Act11Mini struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act11mini"`
		Act12D0 struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act12d0"`
		Act12Mini struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act12mini"`
		Act12Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act12side"`
		Act13D0 struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act13d0"`
		Act13D5 struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act13d5"`
		Act13Mini struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act13mini"`
		Act13Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act13side"`
		Act14Mini struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act14mini"`
		Act14Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act14side"`
		Act15D0 struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act15d0"`
		Act15D5 struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act15d5"`
		Act15Mini struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act15mini"`
		Act15Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act15side"`
		Act16D5 struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act16d5"`
		Act16Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act16side"`
		Act17D0 struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act17d0"`
		Act17Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act17side"`
		Act18D0 struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act18d0"`
		Act18D3 struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act18d3"`
		Act18Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act18side"`
		Act19Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act19side"`
		Act20Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act20side"`
		Act21Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act21side"`
		Act22Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act22side"`
		Act23Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act23side"`
		Act24Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act24side"`
		Act25Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act25side"`
		Act26Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act26side"`
		Act27Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act27side"`
		Act28Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act28side"`
		Act29Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act29side"`
		Act30Side struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act30side"`
		Act3D0 struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act3d0"`
		Act4D0 struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act4d0"`
		Act5D0 struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act5d0"`
		Act6D5 struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act6d5"`
		Act7D5 struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act7d5"`
		Act7Mini struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act7mini"`
		Act8Mini struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act8mini"`
		Act9D0 struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act9d0"`
		Act9Mini struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			StartTime     int    `json:"startTime"`
			EndTime       int    `json:"endTime"`
			RewardEndTime int    `json:"rewardEndTime"`
			IsReplicate   bool   `json:"isReplicate"`
			Type          string `json:"type"`
		} `json:"act9mini"`
	} `json:"activityInfoMap"`
	TowerInfoMap struct {
	} `json:"towerInfoMap"`
	RogueInfoMap struct {
		Rogue1 struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Sort int    `json:"sort"`
		} `json:"rogue_1"`
		Rogue2 struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Sort int    `json:"sort"`
		} `json:"rogue_2"`
		Rogue3 struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Sort int    `json:"sort"`
		} `json:"rogue_3"`
	} `json:"rogueInfoMap"`
	CampaignInfoMap struct {
		Camp01 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_01"`
		Camp02 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_02"`
		Camp03 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_03"`
		CampR01 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_r_01"`
		CampR02 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_r_02"`
		CampR03 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_r_03"`
		CampR04 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_r_04"`
		CampR05 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_r_05"`
		CampR06 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_r_06"`
		CampR07 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_r_07"`
		CampR08 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_r_08"`
		CampR09 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_r_09"`
		CampR10 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_r_10"`
		CampR11 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_r_11"`
		CampR12 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_r_12"`
		CampR13 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_r_13"`
		CampR14 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_r_14"`
		CampR15 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_r_15"`
		CampR16 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_r_16"`
		CampR17 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_r_17"`
		CampR18 struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			CampaignZoneID string `json:"campaignZoneId"`
		} `json:"camp_r_18"`
	} `json:"campaignInfoMap"`
	CampaignZoneInfoMap struct {
		CampZone1 struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"camp_zone_1"`
		CampZone10 struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"camp_zone_10"`
		CampZone11 struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"camp_zone_11"`
		CampZone12 struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"camp_zone_12"`
		CampZone2 struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"camp_zone_2"`
		CampZone3 struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"camp_zone_3"`
		CampZone4 struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"camp_zone_4"`
		CampZone5 struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"camp_zone_5"`
		CampZone6 struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"camp_zone_6"`
		CampZone7 struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"camp_zone_7"`
		CampZone8 struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"camp_zone_8"`
		CampZone9 struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"camp_zone_9"`
	} `json:"campaignZoneInfoMap"`
	EquipmentInfoMap struct {
		Uniequip001Acdrop struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_acdrop"`
		Uniequip001Aglina struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_aglina"`
		Uniequip001Akafyu struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_akafyu"`
		Uniequip001Amgoat struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_amgoat"`
		Uniequip001Amiya struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_amiya"`
		Uniequip001Archet struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_archet"`
		Uniequip001Aurora struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_aurora"`
		Uniequip001Beewax struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_beewax"`
		Uniequip001Bena struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_bena"`
		Uniequip001Bison struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_bison"`
		Uniequip001Blackd struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_blackd"`
		Uniequip001Bldsk struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_bldsk"`
		Uniequip001Bluep struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_bluep"`
		Uniequip001Breeze struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_breeze"`
		Uniequip001Broca struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_broca"`
		Uniequip001Bubble struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_bubble"`
		Uniequip001Ccheal struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_ccheal"`
		Uniequip001Ceylon struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_ceylon"`
		Uniequip001Cgbird struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_cgbird"`
		Uniequip001Chen struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_chen"`
		Uniequip001Chiave struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_chiave"`
		Uniequip001Clour struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_clour"`
		Uniequip001Crow struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_crow"`
		Uniequip001Cutter struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_cutter"`
		Uniequip001Cuttle struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_cuttle"`
		Uniequip001Deepcl struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_deepcl"`
		Uniequip001Doberm struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_doberm"`
		Uniequip001Elysm struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_elysm"`
		Uniequip001Estell struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_estell"`
		Uniequip001Ethan struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_ethan"`
		Uniequip001Finlpp struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_finlpp"`
		Uniequip001Flower struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_flower"`
		Uniequip001Fmout struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_fmout"`
		Uniequip001Folnic struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_folnic"`
		Uniequip001Franka struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_franka"`
		Uniequip001Ghost struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_ghost"`
		Uniequip001Ghost2 struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_ghost2"`
		Uniequip001Glacus struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_glacus"`
		Uniequip001Glassb struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_glassb"`
		Uniequip001Glaze struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_glaze"`
		Uniequip001Grani struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_grani"`
		Uniequip001Gravel struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_gravel"`
		Uniequip001Greyy struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_greyy"`
		Uniequip001Gyuki struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_gyuki"`
		Uniequip001Headbr struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_headbr"`
		Uniequip001Helage struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_helage"`
		Uniequip001Hmau struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_hmau"`
		Uniequip001Hsguma struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_hsguma"`
		Uniequip001Huang struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_huang"`
		Uniequip001Humus struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_humus"`
		Uniequip001Ifrit struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_ifrit"`
		Uniequip001Indigo struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_indigo"`
		Uniequip001Irene struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_irene"`
		Uniequip001Iris struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_iris"`
		Uniequip001Jesica struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_jesica"`
		Uniequip001Lionhd struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_lionhd"`
		Uniequip001Lumen struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_lumen"`
		Uniequip001Lunacu struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_lunacu"`
		Uniequip001Mantic struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_mantic"`
		Uniequip001Meteo struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_meteo"`
		Uniequip001Mint struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_mint"`
		Uniequip001Mizuki struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_mizuki"`
		Uniequip001Mm struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_mm"`
		Uniequip001Moeshd struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_moeshd"`
		Uniequip001Mostma struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_mostma"`
		Uniequip001Myrrh struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_myrrh"`
		Uniequip001Myrtle struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_myrtle"`
		Uniequip001Ncdeer struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_ncdeer"`
		Uniequip001Nights struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_nights"`
		Uniequip001Nothin struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_nothin"`
		Uniequip001Otter struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_otter"`
		Uniequip001Panda struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_panda"`
		Uniequip001Phatom struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_phatom"`
		Uniequip001Platnm struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_platnm"`
		Uniequip001Poca struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_poca"`
		Uniequip001Podego struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_podego"`
		Uniequip001Prove struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_prove"`
		Uniequip001Provs struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_provs"`
		Uniequip001Pudd struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_pudd"`
		Uniequip001Red struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_red"`
		Uniequip001Rfrost struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_rfrost"`
		Uniequip001Robrta struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_robrta"`
		Uniequip001Rope struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_rope"`
		Uniequip001Savage struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_savage"`
		Uniequip001Scave struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_scave"`
		Uniequip001Serum struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_serum"`
		Uniequip001Shining struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_shining"`
		Uniequip001Shotst struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_shotst"`
		Uniequip001Shwaz struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_shwaz"`
		Uniequip001Silent struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_silent"`
		Uniequip001Skadi struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_skadi"`
		Uniequip001Skfire struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_skfire"`
		Uniequip001Skgoat struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_skgoat"`
		Uniequip001Slbell struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_slbell"`
		Uniequip001Slchan struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_slchan"`
		Uniequip001Snakek struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_snakek"`
		Uniequip001Snsant struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_snsant"`
		Uniequip001Sqrrel struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_sqrrel"`
		Uniequip001Strong struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_strong"`
		Uniequip001Sunbr struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_sunbr"`
		Uniequip001Susuro struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_susuro"`
		Uniequip001Swire struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_swire"`
		Uniequip001Swllow struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_swllow"`
		Uniequip001Tachak struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_tachak"`
		Uniequip001Texas struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_texas"`
		Uniequip001Texas2 struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_texas2"`
		Uniequip001Tomimi struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_tomimi"`
		Uniequip001Totter struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_totter"`
		Uniequip001Utage struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_utage"`
		Uniequip001Vigna struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_vigna"`
		Uniequip001Vodfox struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_vodfox"`
		Uniequip001Waaifu struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_waaifu"`
		Uniequip001Wildmn struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_wildmn"`
		Uniequip001Yak struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_yak"`
		Uniequip001Yuki struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_001_yuki"`
		Uniequip002Acdrop struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_acdrop"`
		Uniequip002Aglina struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_aglina"`
		Uniequip002Akafyu struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_akafyu"`
		Uniequip002Amgoat struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_amgoat"`
		Uniequip002Amiya struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_amiya"`
		Uniequip002Archet struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_archet"`
		Uniequip002Aurora struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_aurora"`
		Uniequip002Beewax struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_beewax"`
		Uniequip002Bena struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_bena"`
		Uniequip002Bison struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_bison"`
		Uniequip002Blackd struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_blackd"`
		Uniequip002Bldsk struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_bldsk"`
		Uniequip002Bluep struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_bluep"`
		Uniequip002Breeze struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_breeze"`
		Uniequip002Broca struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_broca"`
		Uniequip002Bubble struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_bubble"`
		Uniequip002Ccheal struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_ccheal"`
		Uniequip002Ceylon struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_ceylon"`
		Uniequip002Cgbird struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_cgbird"`
		Uniequip002Chen struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_chen"`
		Uniequip002Chiave struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_chiave"`
		Uniequip002Clour struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_clour"`
		Uniequip002Crow struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_crow"`
		Uniequip002Cutter struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_cutter"`
		Uniequip002Cuttle struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_cuttle"`
		Uniequip002Deepcl struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_deepcl"`
		Uniequip002Doberm struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_doberm"`
		Uniequip002Elysm struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_elysm"`
		Uniequip002Estell struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_estell"`
		Uniequip002Ethan struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_ethan"`
		Uniequip002Finlpp struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_finlpp"`
		Uniequip002Flower struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_flower"`
		Uniequip002Fmout struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_fmout"`
		Uniequip002Folnic struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_folnic"`
		Uniequip002Franka struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_franka"`
		Uniequip002Ghost struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_ghost"`
		Uniequip002Ghost2 struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_ghost2"`
		Uniequip002Glacus struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_glacus"`
		Uniequip002Glassb struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_glassb"`
		Uniequip002Glaze struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_glaze"`
		Uniequip002Grani struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_grani"`
		Uniequip002Gravel struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_gravel"`
		Uniequip002Greyy struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_greyy"`
		Uniequip002Gyuki struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_gyuki"`
		Uniequip002Headbr struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_headbr"`
		Uniequip002Helage struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_helage"`
		Uniequip002Hmau struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_hmau"`
		Uniequip002Hsguma struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_hsguma"`
		Uniequip002Huang struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_huang"`
		Uniequip002Humus struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_humus"`
		Uniequip002Ifrit struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_ifrit"`
		Uniequip002Indigo struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_indigo"`
		Uniequip002Irene struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_irene"`
		Uniequip002Iris struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_iris"`
		Uniequip002Jesica struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_jesica"`
		Uniequip002Lionhd struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_lionhd"`
		Uniequip002Lumen struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_lumen"`
		Uniequip002Lunacu struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_lunacu"`
		Uniequip002Mantic struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_mantic"`
		Uniequip002Meteo struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_meteo"`
		Uniequip002Mint struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_mint"`
		Uniequip002Mizuki struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_mizuki"`
		Uniequip002Mm struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_mm"`
		Uniequip002Moeshd struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_moeshd"`
		Uniequip002Mostma struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_mostma"`
		Uniequip002Myrrh struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_myrrh"`
		Uniequip002Myrtle struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_myrtle"`
		Uniequip002Ncdeer struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_ncdeer"`
		Uniequip002Nights struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_nights"`
		Uniequip002Nothin struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_nothin"`
		Uniequip002Otter struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_otter"`
		Uniequip002Panda struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_panda"`
		Uniequip002Phatom struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_phatom"`
		Uniequip002Platnm struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_platnm"`
		Uniequip002Poca struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_poca"`
		Uniequip002Podego struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_podego"`
		Uniequip002Prove struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_prove"`
		Uniequip002Provs struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_provs"`
		Uniequip002Pudd struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_pudd"`
		Uniequip002Red struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_red"`
		Uniequip002Rfrost struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_rfrost"`
		Uniequip002Robrta struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_robrta"`
		Uniequip002Rope struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_rope"`
		Uniequip002Savage struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_savage"`
		Uniequip002Scave struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_scave"`
		Uniequip002Serum struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_serum"`
		Uniequip002Shining struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_shining"`
		Uniequip002Shotst struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_shotst"`
		Uniequip002Shwaz struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_shwaz"`
		Uniequip002Silent struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_silent"`
		Uniequip002Skadi struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_skadi"`
		Uniequip002Skfire struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_skfire"`
		Uniequip002Skgoat struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_skgoat"`
		Uniequip002Slbell struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_slbell"`
		Uniequip002Slchan struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_slchan"`
		Uniequip002Snakek struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_snakek"`
		Uniequip002Snsant struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_snsant"`
		Uniequip002Sqrrel struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_sqrrel"`
		Uniequip002Strong struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_strong"`
		Uniequip002Sunbr struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_sunbr"`
		Uniequip002Susuro struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_susuro"`
		Uniequip002Swire struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_swire"`
		Uniequip002Swllow struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_swllow"`
		Uniequip002Tachak struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_tachak"`
		Uniequip002Texas struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_texas"`
		Uniequip002Texas2 struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_texas2"`
		Uniequip002Tomimi struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_tomimi"`
		Uniequip002Totter struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_totter"`
		Uniequip002Utage struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_utage"`
		Uniequip002Vigna struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_vigna"`
		Uniequip002Vodfox struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_vodfox"`
		Uniequip002Waaifu struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_waaifu"`
		Uniequip002Wildmn struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_wildmn"`
		Uniequip002Yak struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_yak"`
		Uniequip002Yuki struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_002_yuki"`
		Uniequip003Aglina struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_003_aglina"`
		Uniequip003Archet struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_003_archet"`
		Uniequip003Chen struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_003_chen"`
		Uniequip003Ghost2 struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_003_ghost2"`
		Uniequip003Helage struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_003_helage"`
		Uniequip003Hsguma struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_003_hsguma"`
		Uniequip003Lumen struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_003_lumen"`
		Uniequip003Mizuki struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_003_mizuki"`
		Uniequip003Phatom struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_003_phatom"`
		Uniequip003Shining struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_003_shining"`
		Uniequip003Shwaz struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_003_shwaz"`
		Uniequip003Skadi struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			TypeName2    string `json:"typeName2"`
			ShiningColor string `json:"shiningColor"`
		} `json:"uniequip_003_skadi"`
	} `json:"equipmentInfoMap"`
	ManufactureFormulaInfoMap struct {
		Num3 struct {
			ID        string `json:"id"`
			ItemID    string `json:"itemId"`
			Count     int    `json:"count"`
			Weight    int    `json:"weight"`
			CostPoint int    `json:"costPoint"`
		} `json:"3"`
		Num4 struct {
			ID        string `json:"id"`
			ItemID    string `json:"itemId"`
			Count     int    `json:"count"`
			Weight    int    `json:"weight"`
			CostPoint int    `json:"costPoint"`
		} `json:"4"`
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
