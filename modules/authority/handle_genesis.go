package authority

import (
	"encoding/json"

	authoritytypes "github.com/e-money/em-ledger/x/authority/types"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/v4/types"
)

// HandleGenesis implements modules.BlockModule
func (m *Module) HandleGenesis(genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "authority/gas price").Msg("parsing genesis")

	// Read the genesis state
	var genState authoritytypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[authoritytypes.ModuleName], &genState)
	if err != nil {
		return err
	}

	return m.db.SaveEMoneyGasPrices(types.NewEMoneyGasPrices(genState.MinGasPrices, genesisDoc.InitialHeight))
}
