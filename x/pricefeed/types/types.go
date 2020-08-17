package types

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
	ID           string  `json:"id"`
	CurrentPrice float64 `json:"current_price"`
	MarketCap    int64   `json:"market_cap"`
}

// NewMarketTicker creates a new instance of MarketTicker
func NewMarketTicker(id string, currentPrice float64, marketCap int64) MarketTicker {
	return MarketTicker{
		ID:           id,
		CurrentPrice: currentPrice,
		MarketCap:    marketCap,
	}
}

// MarketTickers is an array of MarketTicker
type MarketTickers []MarketTicker
