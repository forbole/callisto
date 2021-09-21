package bank

import (
	"encoding/json"
	"fmt"

	authutils "github.com/forbole/bdjuno/modules/auth/utils"

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
		return fmt.Errorf("error while unmarhshaling bank state: %s", err)
	}

	// Store the balances
	accounts, err := authutils.GetGenesisAccounts(appState, cdc)
	if err != nil {
		return fmt.Errorf("error while getting genesis account: %s", err)
	}
	accountsMap := getAccountsMap(accounts)

	var balances []types.AccountBalance
	for _, balance := range bankState.Balances {
		_, ok := accountsMap[balance.Address]
		if !ok {
			continue
		}

		balances = append(balances, types.NewAccountBalance(balance.Address, balance.Coins, doc.InitialHeight))
	}

	return db.SaveAccountBalances(balances)
}

func getAccountsMap(accounts []types.Account) map[string]bool {
	var accountsMap = make(map[string]bool, len(accounts))
	for _, account := range accounts {
		accountsMap[account.Address] = true
	}
	return accountsMap
}
