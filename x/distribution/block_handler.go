package distribution

import (
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/database"
)

// HandleBlock represents a method that is called each time a new block is created
func HandleBlock(block *tmctypes.ResultBlock, client distrtypes.QueryClient, db *database.BigDipperDb) error {

	// Update the validator commissions
	err := updateValidatorsCommissionAmounts(block.Block.Height, client, db)
	if err != nil {
		log.Error().Str("module", "distribution").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating validators commissions")
	}

	// Update the delegators commissions amounts
	err = updateDelegatorsCommissionsAmounts(block.Block.Height, client, db)
	if err != nil {
		log.Error().Str("module", "distribution").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating delegators commissions amounts")
	}

	return nil
}
