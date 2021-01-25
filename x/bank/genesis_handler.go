package bank

import (
	"encoding/json"
	"fmt"

	bbanktypes "github.com/forbole/bdjuno/x/bank/types"

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
	balances := make([]bbanktypes.AccountBalance, len(bankState.Balances))
	for index, balance := range bankState.Balances {
		balances[index] = bbanktypes.NewAccountBalance(balance.Address, balance.Coins, 1)
	}

	err := db.SaveAccountBalances(balances)
	if err != nil {
		return fmt.Errorf("error while storing genesis balances: %s", err)
	}

	return nil
}
