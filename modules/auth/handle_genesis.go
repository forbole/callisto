package auth

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(_ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "auth").Msg("parsing genesis")

	// Handle Genesis account addresses
	accounts, err := GetGenesisAccounts(appState, m.cdc)
	if err != nil {
		return fmt.Errorf("error while getting genesis accounts: %s", err)
	}

	err = m.db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing genesis accounts: %s", err)
	}

	// Handle Genesis vesting account details
	vestingAccounts, err := GetGenesisVestingAccounts(appState)
	if err != nil {
		return fmt.Errorf("error while getting genesis vesting accounts: %s", err)
	}

	paramsNumber := 5
	err = m.db.SaveVestingAccounts(paramsNumber, vestingAccounts)
	if err != nil {
		return fmt.Errorf("error while storing genesis vesting accounts: %s", err)
	}

	return nil
}
