package coingecko

import "time"

// Token contains the information of a single token
type Token struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

// Tokens represents a list of Token objects
type Tokens []Token

// MarketTicker contains the current market data for a single token
type MarketTicker struct {
	Symbol       string    `json:"symbol"`
	CurrentPrice float64   `json:"current_price"`
	MarketCap    float64   `json:"market_cap"`
	LastUpdated  time.Time `json:"last_updated"`
}

// MarketTickers is an array of MarketTicker
type MarketTickers []MarketTicker
