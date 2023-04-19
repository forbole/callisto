package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	superfluidtypes "github.com/osmosis-labs/osmosis/v15/x/superfluid/types"
	"github.com/forbole/juno/v4/node/local"

	superfluidsource "github.com/forbole/bdjuno/v4/modules/superfluid/source"
)

var (
	_ superfluidsource.Source = &Source{}
)

// Source implements superfluidsource.Source reading the data from a local node
type Source struct {
	*local.Source
	q superfluidtypes.QueryServer
}

func NewSource(source *local.Source, keeper superfluidtypes.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      keeper,
	}
}

// GetSuperfluidDelegationsByDelegator implements superfluidsource.Source
func (s Source) GetSuperfluidDelegationsByDelegator(address string, height int64) ([]superfluidtypes.SuperfluidDelegationRecord, error) {
		ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}
	res, err := s.q.SuperfluidDelegationsByDelegator(
				sdk.WrapSDKContext(ctx),
		&superfluidtypes.SuperfluidDelegationsByDelegatorRequest{DelegatorAddress: address},
	)
	if err != nil {
		return nil, err
	}

	return res.SuperfluidDelegationRecords, nil
}
