package structs

type EskhataCreditsUrlAnchors struct {
	CommonUrl            string `json:"commonUrl"`
	MultiPurposeConsumer string `json:"multiPurposeConsumer"`
	Express              string `json:"express"`
	HunarkhoiMardumi     string `json:"hunarkhoiMardumi"`
	Imkon                string `json:"imkon"`
	Car                  string `json:"car"`
	Manzili              string `json:"manzili"`
}

type EskhataCreditProducts struct {
	MultiPurposeConsumer Credit `json:"multiPurposeConsumer"`
	Express              Credit `json:"express"`
	HunarkhoiMardumi     Credit `json:"hunarkhoiMardumi"`
	Imkon                Credit `json:"imkon"`
	Car                  Credit `json:"car"`
	Manzili              Credit `json:"manzili"`
}
