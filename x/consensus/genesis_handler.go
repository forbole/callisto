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

func GenesisHandler(_ *codec.Codec, genesisDoc *tmtypes.GenesisDoc, _ map[string]json.RawMessage, w worker.Worker) error {
	log.Debug().Str("module", "consensus").Msg("parsing genesis")

	bigDipperDb, ok := w.Db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("given database instance is not a BigDipperDb")
	}

	if err := saveGenesisTime(genesisDoc.GenesisTime, bigDipperDb); err != nil {
		return err
	}
	return nil
}

func saveGenesisTime(genesisTime time.Time, db database.BigDipperDb) error {
	return db.SaveGenesisTime(genesisTime)
}
