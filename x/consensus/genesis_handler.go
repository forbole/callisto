package consensus

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/juno/parse/worker"
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

func GenesisHandler(codec *codec.Codec, genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, w worker.Worker) error {
	log.Debug().Str("module", "gov").Msg("parsing genesis")

	bigDipperDb, ok := w.Db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("given database instance is not a BigDipperDb")
	}

	if err := saveGenesisHeight(genesisDoc.GenesisTime, bigDipperDb); err != nil {
		return err
	}
	return nil
}

func saveGenesisHeight(genesisTime time.Time, db database.BigDipperDb) error {
	db.SaveGenesisHeight(genesisTime)
}
