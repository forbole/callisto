package wormhole

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v4/types"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
	wormholetypes "github.com/wormhole-foundation/wormchain/x/wormhole/types"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "wormhole").Msg("parsing genesis")

	// Read the genesis state
	var genState wormholetypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[wormholetypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while unmarshaling wormhole state: %s", err)
	}

	// Save the config
	err = m.db.SaveWormholeConfig(types.NewWormholeConfig(genState.Config, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis wormhole config: %s", err)
	}

	err = m.db.SaveGuardianValidatorList(genState.GuardianValidatorList, doc.InitialHeight)
	if err != nil {
		return fmt.Errorf("error while storing genesis guardian validator list: %s", err)
	}

	err = m.db.SaveGuardianSetList(genState.GuardianSetList, doc.InitialHeight)
	if err != nil {
		return fmt.Errorf("error while storing genesis guardian set list: %s", err)
	}

	err = m.db.SaveValidatorAllowListFromGenesis(genState.AllowedAddresses, doc.InitialHeight)
	if err != nil {
		return fmt.Errorf("error while storing genesis validator allow list: %s", err)
	}

	return nil
}
