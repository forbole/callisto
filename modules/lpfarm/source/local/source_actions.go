package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	lpfarmtypes "github.com/crescent-network/crescent/v5/x/lpfarm/types"

	"github.com/forbole/juno/v4/node/local"

	lpfarmsource "github.com/forbole/bdjuno/v4/modules/lpfarm/source"
)

var (
	_ lpfarmsource.Source = &Source{}
)

// Source implements lpfarmsource.Source using a local node
type Source struct {
	*local.Source
	querier lpfarmtypes.QueryServer
}

// NewSource returns a new Source instace
func NewSource(source *local.Source, querier lpfarmtypes.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetTotalLPFarmRewards implements lpfarmsource.Source
func (s Source) GetTotalLPFarmRewards(farmer string, height int64) (sdk.DecCoins, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return sdk.DecCoins{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.TotalRewards(
		sdk.WrapSDKContext(ctx),
		&lpfarmtypes.QueryTotalRewardsRequest{Farmer: farmer},
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting total lp farm rewards for farmer %s at height %v: %s", farmer, height, err)
	}

	return res.Rewards, nil
}
