package bank

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
)

// HandleGenesis handles the genesis state of the x/bank module in order to store the initial values
// of the different accounts.
func HandleGenesis(appState map[string]json.RawMessage, cdc codec.Marshaler, db *database.BigDipperDb) error {
	log.Debug().Str("module", "bank").Msg("parsing genesis")

	var bankState banktypes.GenesisState
	if err := cdc.UnmarshalJSON(appState[banktypes.ModuleName], &bankState); err != nil {
		return err
	}

	// Store the accounts
	for _, balance := range bankState.Balances {
		err := db.SaveAccountBalance(balance.Address, balance.Coins, 1)
		if err != nil {
			return err
		}
	}

	return nil
}
