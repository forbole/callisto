package feeexcluder

import (
	"encoding/json"

	feeexcluder "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(_ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", feeexcluder.ModuleName).Msg("parsing genesis")

	// Unmarshal the bank state
	var genesisState feeexcluder.GenesisState
	if err := m.cdc.UnmarshalJSON(appState[feeexcluder.ModuleName], &genesisState); err != nil {
		return err
	}

	return m.feeexcluderRepo.InsertToGenesisState(genesisState)
}
