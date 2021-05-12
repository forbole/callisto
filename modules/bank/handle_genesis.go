package bank

import (
	"encoding/json"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/database"

	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/types"
)

// HandleGenesis handles the genesis state of the x/bank module in order to store the initial values
// of the different account balances.
func HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, cdc codec.Marshaler, db *database.Db) error {
	log.Debug().Str("module", "bank").Msg("parsing genesis")

	var bankState banktypes.GenesisState
	if err := cdc.UnmarshalJSON(appState[banktypes.ModuleName], &bankState); err != nil {
		return err
	}

	// Store the accounts
	balances := make([]types.AccountBalance, len(bankState.Balances))
	for index, balance := range bankState.Balances {
		balances[index] = types.NewAccountBalance(balance.Address, balance.Coins, doc.InitialHeight)
	}

	return db.SaveAccountBalances(balances)
}
