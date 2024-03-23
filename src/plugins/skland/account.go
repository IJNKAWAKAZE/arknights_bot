package skland

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
