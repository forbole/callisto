package pricefeed

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/bdjuno/database"

	"github.com/forbole/bdjuno/types"

	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleGenesis handles the genes for the pricefeed module.
// It creates the proper rows of all tokens and sets their default prices to 0
func HandleGenesis(
	genDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, cdc codec.JSONMarshaler, db *database.Db,
) error {
	log.Debug().Str("module", "pricefeed").Msg("parsing genesis")

	var bankState banktypes.GenesisState
	cdc.MustUnmarshalJSON(appState[banktypes.ModuleName], &bankState)

	var supply = sdk.Coins{}
	for _, balance := range bankState.Balances {
		supply = supply.Add(balance.Coins...)
	}

	var prices []types.TokenPrice
	for _, coin := range supply {
		// Save the coin as a token with its units
		err := db.SaveToken(types.NewToken(coin.Denom, []types.TokenUnit{
			types.NewTokenUnit(coin.Denom, 0, nil),
		}))
		if err != nil {
			return err
		}

		// Create the price entry
		prices = append(prices, types.NewTokenPrice(coin.Denom, 0, 0, genDoc.GenesisTime))
	}

	return db.SaveTokensPrices(prices)
}
