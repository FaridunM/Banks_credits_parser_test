package structs

type HumoCreditsUrl struct {
	Consumer      string `json:"consumer"`
	StartBusiness string `json:"startBusiness"`
	Mortgage      string `json:"mortgage"`
	Education     string `json:"education"`
	Agricultural  string `json:"agricultural"`
	Livestock     string `json:"livestock"`
	Orzu          string `json:"orzu"`
}

type HumoCreditProducts struct {
	Consumer      Credit
	StartBusiness Credit
	Mortgage      Credit
	Education     Credit
	Agricultural  Credit
	Livestock     Credit
	Orzu          Credit
}
