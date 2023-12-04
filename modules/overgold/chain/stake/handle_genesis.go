package stake

import (
	"encoding/json"

	stake "git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(_ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", stake.ModuleName).Msg("parsing genesis")

	// Unmarshal the stake state
	var state stake.GenesisState
	if err := m.cdc.UnmarshalJSON(appState[stake.ModuleName], &state); err != nil {
		return err
	}

	return nil
}
