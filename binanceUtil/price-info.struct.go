package binanceUtil

type PriceStatStruct struct {
	Symbol    string    `json:"symbol"`
	Price     string    `json:"lastPrice"`
	Volume    string    `json:"quoteVolume"`
}

type PriceStats []PriceStatStruct
