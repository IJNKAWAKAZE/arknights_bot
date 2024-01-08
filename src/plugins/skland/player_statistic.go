package skland

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
