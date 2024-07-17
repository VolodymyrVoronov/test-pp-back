package models

type StockData struct {
	Date  string  `json:"date"`
	Open  float64 `json:"open"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Close float64 `json:"close"`
}

type StockRequest struct {
	Stocks []StockData `json:"stocks"`
}

type Prediction struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
}
