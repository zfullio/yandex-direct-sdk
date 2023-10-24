package anti_fraud

//Запросы слать на: http: //api.direct.yandex.com/json/v5/conversionscores

type Req struct {
	Method string `json:"method"`
	Params Params `json:"params"`
}

type Params struct {
	FieldNames        []string          `json:"FieldNames"`
	SelectionCriteria SelectionCriteria `json:"SelectionCriteria"`
}

type SelectionCriteria struct {
	Requests []Requests `json:"Requests"`
}

type Requests struct {
	Yclid string `json:"Yclid"`
	Email string `json:"Email,omitempty"`
	Phone string `json:"Phone,omitempty"`
}

type Response struct {
	Result struct {
		ConversionScores []struct {
			Yclid string      `json:"Yclid"`
			Email interface{} `json:"Email"`
			Phone interface{} `json:"Phone"`
			Score int         `json:"Score"`
		} `json:"ConversionScores"`
	} `json:"result"`
}
