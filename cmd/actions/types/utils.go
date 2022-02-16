package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type GraphQLError struct {
	Message string `json:"message"`
}

func ConvertSdkCoins(sdkCoin sdk.Coins) []Coin {
	actionCoins := make([]Coin, len(sdkCoin))
	for index, s := range sdkCoin {
		actionCoins[index] = Coin{
			Denom:  s.Denom,
			Amount: s.Amount.String(),
		}
	}
	return actionCoins
}
