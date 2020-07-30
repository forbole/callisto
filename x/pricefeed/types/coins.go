package pricefeed

// MarketTicker contains the current market data for a single coin
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
