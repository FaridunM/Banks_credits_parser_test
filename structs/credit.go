package structs

type CreditsUrl struct {
	Humo    HumoCreditsUrl           `json:"humo"`
	Eskhata EskhataCreditsUrlAnchors `json:"eskhata"`
}

type Credit struct {
	Title              string              `json:"title,omitempty"`
	Description        string              `json:"description,omitempty"`
	Purpose            string              `json:"purpose,omitempty"`
	LandingMethodology string              `json:"landingMethodology,omitempty"`
	Collateral         []string            `json:"collateral,omitempty"`
	InitialPayment     string              `json:"initialPayment,omitempty"`
	Amount             []string            `json:"amount,omitempty"`
	Currency           string              `json:"currency,omitempty"`
	InterestRate       []string            `json:"interestRate,omitempty"`
	GracePeriod        string              `json:"gracePeriod,omitempty"`
	LoanPeriod         string              `json:"loanPeriod,omitempty"`
	LoanTerms          []string            `json:"loanTerms,omitempty"`
	NeedDocuments      []string            `json:"needDocuments,omitempty"`
	Url                string              `json:"url,omitempty"`
	AdditionalFields   map[string][]string `json:"additionalFields,omitempty"`
}
