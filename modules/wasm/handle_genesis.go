package wasm

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v3/types"
	tmtypes "github.com/tendermint/tendermint/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "wasm").Msg("parsing genesis")

	// Read the genesis state
	var genState wasmtypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[wasmtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while unmarshaling wasm genesis state: %s", err)
	}

	err = m.SaveGenesisParams(genState.Params, doc.InitialHeight)
	if err != nil {
		return fmt.Errorf("error while saving genesis wasm params: %s", err)
	}

	return nil
}

func (m *Module) SaveGenesisParams(params wasmtypes.Params, initHeight int64) error {
	err := m.db.SaveWasmParams(
		types.NewWasmParams(&params.CodeUploadAccess, int32(params.InstantiateDefaultPermission), params.MaxWasmCodeSize, initHeight),
	)
	if err != nil {
		return fmt.Errorf("error while saving genesis wasm params: %s", err)
	}
	return nil
}
