package pricefeed

import (
	"encoding/json"
	"fmt"

	juno "github.com/desmos-labs/juno/types"

	historyutils "github.com/forbole/bdjuno/modules/history/utils"
	"github.com/forbole/bdjuno/modules/utils"

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
	cfg juno.Config,
	genDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage,
	cdc codec.JSONMarshaler, db *database.Db,
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

	err := db.SaveTokensPrices(prices)
	if err != nil {
		return fmt.Errorf("error while storing token prices: %s", err)
	}

	if utils.IsModuleEnabled(cfg, types.HistoryModuleName) {
		return historyutils.UpdatePriceHistory(prices, db)
	}

	return nil
}
