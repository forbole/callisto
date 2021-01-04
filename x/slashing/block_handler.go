package slashing

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/desmos-labs/juno/client"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/database"
	slashingtypes "github.com/forbole/bdjuno/x/slashing/types"
)

// HandleBlock represents a method that is called each time a new block is created
func HandleBlock(block *tmctypes.ResultBlock, cp *client.Proxy, db *database.BigDipperDb) error {
	// Update the staking pool
	err := updateSigningInfo(block.Block.Height, cp, db)
	if err != nil {
		log.Error().Str("module", "slashing").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating signing info")
	}

	return nil
}

// updateSigningInfo reads from the LCD the current staking pool and stores its value inside the database
func updateSigningInfo(height int64, cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().Str("module", "slashing").Int64("height", height).
		Str("operation", "signing info").Msg("getting signing info")

	var pool []slashing.ValidatorSigningInfo
	endpoint := fmt.Sprintf("/slashing/signing_infos?height=%d", height)
	height, err := cp.QueryLCDWithHeight(endpoint, &pool)
	if err != nil {
		log.Err(err).Str("module", "slashing").Msg("error while getting signing info")
		return err
	}

	log.Debug().Str("module", "slashing").Int64("height", height).
		Str("operation", "signing info").Msg("saving signing info")

	infos := make([]slashingtypes.ValidatorSigningInfo, len(pool))
	for index, info := range pool {
		infos[index] = slashingtypes.NewValidatorSigningInfo(
			info.Address.String(),
			info.StartHeight,
			info.IndexOffset,
			info.JailedUntil,
			info.Tombstoned,
			info.MissedBlocksCounter,
			height,
		)
	}

	return db.SaveValidatorsSigningInfos(infos)
}
