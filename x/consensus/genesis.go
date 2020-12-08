package consensus

import (
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

func HandleGenesis(genesisDoc *tmtypes.GenesisDoc, db *database.BigDipperDb) error {
	log.Debug().Str("module", "consensus").Msg("parsing genesis")

	// Save the genesis time
	err := db.SaveGenesisTime(genesisDoc.GenesisTime)
	if err != nil {
		return err
	}

	return nil
}
