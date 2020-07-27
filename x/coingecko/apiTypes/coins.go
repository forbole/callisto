package coingecko

import "fmt"

type Coin struct {
	Id     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

func (coin Coin) String() string {
	return fmt.Sprintf("%s%s%s", coin.Id, coin.Symbol, coin.Name)
}

func NewCoin(id string, symbol string, name string) Coin {
	return Coin{
		Id:     id,
		Symbol: symbol,
		Name:   name,
	}
}

type Coins []Coin
