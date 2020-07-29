package pricefeed

// Coin this represent coingecko api's coins/list attributes
type Coin struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

// NewCoin return a new pricefeed coin instance
func NewCoin(id string, symbol string, name string) Coin {
	return Coin{
		ID:     id,
		Symbol: symbol,
		Name:   name,
	}
}

// Coins is an array of Coin
type Coins []Coin

//____________________________________________________

// MarketTicker represent some of the attributes in coingecko api's coins/market
type MarketTicker struct {
	ID           string  `json:"id"`
	CurrentPrice float64 `json:"current_price"`
	MarketCap    int64   `json:"market_cap"`
}

// NewMarket return an instance of MarketTicker
func NewMarket(id string, currentPrice float64, marketCap int64) MarketTicker {
	return MarketTicker{
		ID:           id,
		CurrentPrice: currentPrice,
		MarketCap:    marketCap,
	}
}

// Pricefeeds is an array of MarketTicker
type Pricefeeds []MarketTicker
