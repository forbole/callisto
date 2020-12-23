package slashing

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/desmos-labs/juno/client"
	"github.com/forbole/bdjuno/database"
	slashingtypes "github.com/forbole/bdjuno/x/slashing/types"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock represents a method that is called each time a new block is created
func HandleBlock(block *tmctypes.ResultBlock, cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().Str("module", "staking").Msgf("handling block")

	// Update the staking pool
	err := updateSigningInfo(block.Block.Height, cp, db)
	if err != nil {
		return err
	}

	return nil
}

// updateSigningInfo reads from the LCD the current staking pool and stores its value inside the database
func updateSigningInfo(height int64, cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().
		Str("module", "slashing").
		Str("operation", "signing info").
		Msg("getting signing info")

	var pool []slashing.ValidatorSigningInfo
	endpoint := fmt.Sprintf("/slashing/signing_infos?height=%d", height)
	height, err := cp.QueryLCDWithHeight(endpoint, &pool)
	if err != nil {
		log.Err(err).Str("module", "slashing").Msg("error while getting signing info")
		return err
	}

	log.Debug().
		Str("module", "staking").
		Str("operation", "staking_pool").
		Msg("saving staking pool")

	for _, info := range pool {
		err = db.SaveSigningInfos(slashingtypes.NewValidatorSigningInfo(
			info.Address,
			info.StartHeight,
			info.IndexOffset,
			info.JailedUntil,
			info.Tombstoned,
			info.MissedBlocksCounter,
			height, time.Now(),
		))
		if err != nil {
			return err
		}
	}

	return nil
}
