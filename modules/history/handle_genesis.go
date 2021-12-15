package history

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, _ map[string]json.RawMessage) error {
	log.Debug().Str("module", "history").Msg("parsing genesis")

	accounts, err := m.db.GetAccounts()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		err = m.UpdateAccountBalanceHistoryWithTime(account, doc.GenesisTime)
		if err != nil {
			return err
		}
	}

	return nil
}
