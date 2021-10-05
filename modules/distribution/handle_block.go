package distribution

import (
	"context"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/desmos-labs/juno/client"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/types"

	"github.com/forbole/bdjuno/database"
	distrutils "github.com/forbole/bdjuno/modules/distribution/utils"
	"github.com/forbole/bdjuno/types/config"
)

// HandleBlock represents a method that is called each time a new block is created
func HandleBlock(cfg *config.Config, block *tmctypes.ResultBlock, client distrtypes.QueryClient, db *database.Db) error {
	go updateParams(block.Block.Height, client, db)

	// Update the validator commissions
	go distrutils.UpdateValidatorsCommissionAmounts(cfg, block.Block.Height, client, db)

	// Update the delegators commissions amounts
	go distrutils.UpdateDelegatorsRewardsAmounts(cfg, block.Block.Height, client, db)

	return nil
}

func updateParams(height int64, distrClient distrtypes.QueryClient, db *database.Db) {
	log.Debug().Str("module", "distribution").Int64("height", height).
		Msg("updating params")

	res, err := distrClient.Params(
		context.Background(),
		&distrtypes.QueryParamsRequest{},
		client.GetHeightRequestHeader(height),
	)
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).
			Int64("height", height).
			Msg("error while getting params")
		return
	}

	err = db.SaveDistributionParams(types.NewDistributionParams(res.Params, height))
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).
			Int64("height", height).
			Msg("error while saving params")
		return
	}
}
