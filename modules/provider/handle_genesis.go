package provider

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v4/types"
	"github.com/rs/zerolog/log"

	providertypes "github.com/akash-network/node/x/provider/types/v1beta2"
	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "provider").Msg("parsing genesis")

	// Read the genesis state
	var genState providertypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[providertypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while unmarshaling provider state: %s", err)
	}

	err = m.saveGenesisProviders(genState.Providers, doc.InitialHeight)
	if err != nil {
		return fmt.Errorf("error while saving genesis providers: %s", err)
	}

	return nil
}

func (m *Module) saveGenesisProviders(genProviders []providertypes.Provider, height int64) error {
	providers := make([]*types.Provider, len(genProviders))
	for index, info := range genProviders {
		providers[index] = types.NewProvider(info, height)
	}

	err := m.db.SaveProviders(providers, height)
	if err != nil {
		return err
	}
	return nil
}
