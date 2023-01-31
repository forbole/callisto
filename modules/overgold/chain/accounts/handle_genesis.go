package accounts

import (
	"encoding/json"

	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "accounts").Msg("parsing genesis")

	// Unmarshal the bank state
	var accountsState accountstypes.GenesisState
	if err := m.cdc.UnmarshalJSON(appState[accountstypes.ModuleName], &accountsState); err != nil {
		return err
	}

	return m.accountRepo.SaveAccounts(accountsState.Accounts...)
}