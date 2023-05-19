package auth

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(_ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "auth").Msg("parsing genesis")

	accounts, err := GetGenesisAccounts(appState, m.cdc)
	if err != nil {
		return fmt.Errorf("error while getting genesis accounts: %s", err)
	}
	err = m.db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing genesis accounts: %s", err)
	}

	vestingAccounts, err := GetGenesisVestingAccounts(appState, m.cdc)
	if err != nil {
		return fmt.Errorf("error while getting genesis vesting accounts: %s", err)
	}
	err = m.db.SaveVestingAccounts(vestingAccounts)
	if err != nil {
		return fmt.Errorf("error while storing genesis vesting accounts: %s", err)
	}

	return nil
}
