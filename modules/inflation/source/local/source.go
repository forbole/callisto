package remote

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	inflationtypes "github.com/e-money/em-ledger/x/inflation/types"
	"github.com/forbole/juno/v2/node/local"

	inflationsource "github.com/forbole/bdjuno/v2/modules/inflation/source"
)

var (
	_ inflationsource.Source = &Source{}
)

// Source implements inflationsource.Source using a remote node
type Source struct {
	*local.Source
	client inflationtypes.QueryServer
}

// NewSource returns a new Source instance
func NewSource(source *local.Source, client inflationtypes.QueryServer) *Source {
	return &Source{
		Source: source,
		client: client,
	}
}

// GetInflation implements inflationsource.Source
func (s *Source) GetInflation(height int64) (inflationtypes.InflationState, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return inflationtypes.InflationState{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.client.Inflation(sdk.WrapSDKContext(ctx), &inflationtypes.QueryInflationRequest{})
	if err != nil {
		return inflationtypes.InflationState{}, fmt.Errorf("error while reading inflation state: %s", err)
	}

	return res.State, nil
}
