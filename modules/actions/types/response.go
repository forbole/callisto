package types

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Coin struct {
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

func ConvertCoins(coins sdk.Coins) []Coin {
	amount := make([]Coin, 0)
	for _, coin := range coins {
		amount = append(amount, Coin{Amount: coin.Amount.String(), Denom: coin.Denom})
	}
	return amount
}

func ConvertDecCoins(coins sdk.DecCoins) []Coin {
	amount := make([]Coin, 0)
	for _, coin := range coins {
		amount = append(amount, Coin{Amount: coin.Amount.String(), Denom: coin.Denom})
	}
	return amount
}

// ========================= Account Balance Response =========================

type Balance struct {
	Coins []Coin `json:"coins"`
}
