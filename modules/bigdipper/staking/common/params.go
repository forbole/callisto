package common

import (
	"context"

	"github.com/forbole/bdjuno/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/rs/zerolog/log"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"
	utils2 "github.com/forbole/bdjuno/modules/common/utils"
)

// UpdateParams updates the staking parameters for the given height,
// storing them inside the database and returning its value
func UpdateParams(
	height int64, client stakingtypes.QueryClient, db *bigdipperdb.Db,
) (*stakingtypes.Params, error) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating params")

	res, err := client.Params(
		context.Background(),
		&stakingtypes.QueryParamsRequest{},
		utils2.GetHeightRequestHeader(height),
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
