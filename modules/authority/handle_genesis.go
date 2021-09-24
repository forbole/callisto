package authority

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	authoritytypes "github.com/e-money/em-ledger/x/authority/types"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

func HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, cdc codec.Marshaler, db *database.Db) error {
	log.Debug().Str("module", "authority/gas price").Msg("parsing genesis")

	// Read the genesis state
	var genState authoritytypes.GenesisState
	err := cdc.UnmarshalJSON(appState[authoritytypes.ModuleName], &genState)
	if err != nil {
		return err
	}

	newEmoneyGasPrice := types.NewEmoneyGasPrice(genState.AuthorityKey, genState.MinGasPrices, doc.InitialHeight)

	return db.SaveEmoneyGasPrices(newEmoneyGasPrice)
}
