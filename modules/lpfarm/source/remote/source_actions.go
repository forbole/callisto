package remote

import (
	"fmt"

	lpfarmtypes "github.com/crescent-network/crescent/v5/x/lpfarm/types"
	"github.com/forbole/juno/v4/node/remote"

	sdk "github.com/cosmos/cosmos-sdk/types"
	lpfarmsource "github.com/forbole/bdjuno/v4/modules/lpfarm/source"
)

var (
	_ lpfarmsource.Source = &Source{}
)

// Source implements lpfarmsource.Source using a remote node
type Source struct {
	*remote.Source
	querier lpfarmtypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, querier lpfarmtypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetTotalLPFarmRewards implements lpfarmsource.Source
func (s Source) GetTotalLPFarmRewards(farmer string, height int64) (sdk.DecCoins, error) {
	res, err := s.querier.TotalRewards(
		remote.GetHeightRequestContext(s.Ctx, height),
		&lpfarmtypes.QueryTotalRewardsRequest{Farmer: farmer},
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting total lp farm rewards for farmer %s at height %v: %s", farmer, height, err)
	}

	return res.Rewards, nil
}
