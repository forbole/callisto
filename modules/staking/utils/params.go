package utils

import (
	"context"

	"github.com/desmos-labs/juno/client"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/rs/zerolog/log"
)

// UpdateParams updates the staking parameters for the given height,
// storing them inside the database and returning its value
func UpdateParams(
	height int64, stakingClient stakingtypes.QueryClient, db *database.Db,
) (*stakingtypes.Params, error) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating params")

	res, err := stakingClient.Params(
		context.Background(),
		&stakingtypes.QueryParamsRequest{},
		client.GetHeightRequestHeader(height),
	)
	if err != nil {
		return nil, err
	}

	err = db.SaveStakingParams(types.NewStakingParams(res.Params.BondDenom))
	if err != nil {
		return nil, err
	}

	return &res.Params, nil
}
