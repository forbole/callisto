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
	MarketCap    int64     `json:"market_cap"`
	LastUpdated  time.Time `json:"last_updated"`
}

// NewMarketTicker creates a new instance of MarketTicker
func NewMarketTicker(symbol string, currentPrice float64, marketCap int64, lastUpdated time.Time) MarketTicker {
	return MarketTicker{
		Symbol:       symbol,
		CurrentPrice: currentPrice,
		MarketCap:    marketCap,
		LastUpdated:  lastUpdated,
	}
}

// MarketTickers is an array of MarketTicker
type MarketTickers []MarketTicker
