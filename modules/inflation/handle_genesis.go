package inflation

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	tmtypes "github.com/tendermint/tendermint/types"

	inflationtypes "github.com/e-money/em-ledger/x/inflation/types"
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"
)

func HandleGenesis(
	genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, cdc codec.Marshaler, db *database.Db,
) error {
	log.Debug().Str("module", "inflation").Msg("parsing genesis")

	// Read the genesis state
	var genState inflationtypes.GenesisState
	err := cdc.UnmarshalJSON(appState[inflationtypes.ModuleName], &genState)
	if err != nil {
		return err
	}

	return db.SaveEmoneyInflation(genState.InflationState, genesisDoc.InitialHeight)
}
