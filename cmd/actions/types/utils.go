package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type GraphQLError struct {
	Message string `json:"message"`
}

func ConvertSdkCoins(s sdk.Coins) []Coin {
	actionCoins := make([]Coin, len(s))
	for index, c := range s {
		actionCoins[index] = Coin{
			Denom:  c.Denom,
			Amount: c.Amount.String(),
		}
	}
	return actionCoins
}

func ConvertSdkDecCoins(s sdk.DecCoins) []Coin {
	actionCoins := make([]Coin, len(s))
	for index, c := range s {
		actionCoins[index] = Coin{
			Denom:  c.Denom,
			Amount: c.Amount.String(),
		}
	}
	return actionCoins
}
