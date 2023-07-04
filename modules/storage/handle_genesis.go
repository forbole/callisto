package storage

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/forbole/bdjuno/v4/types"

	storagetypes "github.com/jackalLabs/canine-chain/v3/x/storage/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "storage").Msg("parsing genesis")

	// Read the genesis state
	var genState storagetypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[storagetypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading storage genesis data: %s", err)
	}

	// Save the params
	err = m.db.SaveStorageParams(types.NewStorageParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis storage params: %s", err)
	}

	return nil
}
