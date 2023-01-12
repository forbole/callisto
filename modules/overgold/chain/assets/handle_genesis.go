package assets

import (
	"encoding/json"

	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/rs/zerolog/log"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "assets").Msg("parsing genesis")

	// Unmarshal the asset state
	var assetsState assetstypes.GenesisState
	if err := m.cdc.UnmarshalJSON(appState[assetstypes.ModuleName], &assetsState); err != nil {
		return err
	}

	return m.assetRepo.SaveAssets(assetsState.Assets...)
}
