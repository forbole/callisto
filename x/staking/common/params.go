package common

import (
	"context"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"
	"github.com/forbole/bdjuno/x/utils"
)

// UpdateParams updates the staking parameters for the given height,
// storing them inside the database and returning its value
func UpdateParams(
	height int64, client stakingtypes.QueryClient, db *database.BigDipperDb,
) (*stakingtypes.Params, error) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating params")

	res, err := client.Params(
		context.Background(),
		&stakingtypes.QueryParamsRequest{},
		utils.GetHeightRequestHeader(height),
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
