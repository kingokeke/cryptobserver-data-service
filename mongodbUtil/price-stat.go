package mongodbUtil

import "time"

type PriceStatStruct struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"lastPrice"`
	Volume    float64 `json:"quoteVolume"`
	Timestamp time.Time   `bson:"timestamp" json:"updated_at,omitempty"`
}
