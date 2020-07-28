package coingecko

// Coin this represent coingecko api's coins/list attributes
type Coin struct {
	Id     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

// NewCoin return a new coingecko coin instance
func NewCoin(id string, symbol string, name string) Coin {
	return Coin{
		Id:     id,
		Symbol: symbol,
		Name:   name,
	}
}

// Coins is an array of Coin
type Coins []Coin

//____________________________________________________

// Market represent some of the attributes in coingecko api's coins/market
type Market struct {
	Id           string  `json:"id"`
	CurrentPrice float64 `json:"current_price"`
	MarketCap    int64   `json:"market_cap"`
}

// NewMarket return an instance of Market
func NewMarket(id string, currentPrice float64, marketCap int64) Market {
	return Market{
		Id:           id,
		CurrentPrice: currentPrice,
		MarketCap:    marketCap,
	}
}

// Markets is an array of Market
type Markets []Market
