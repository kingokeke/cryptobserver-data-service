package binanceUtil

type PriceStatStruct struct {
	Symbol    string    `json:"symbol"`
	Price     string    `json:"lastPrice"`
	Volume    string    `json:"quoteVolume"`

	// Highprice          string `json:"highPrice"`
	// Lowprice           string `json:"lowPrice"`
}

type PriceStats []PriceStatStruct
