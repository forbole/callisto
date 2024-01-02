package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	inflationtypes "github.com/evmos/evmos/v15/x/inflation/types"
	"github.com/forbole/juno/v5/node/local"

	inflationsource "github.com/forbole/bdjuno/v4/modules/inflation/source"
)

var (
	_ inflationsource.Source = &Source{}
)

// Source implements inflationsource.Source using a local node
type Source struct {
	*local.Source
	querier inflationtypes.QueryServer
}

// NewSource returns a new Source instace
func NewSource(source *local.Source, querier inflationtypes.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// Params implements inflationsource.Source
func (s Source) Params(height int64) (inflationtypes.Params, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return inflationtypes.Params{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.Params(sdk.WrapSDKContext(ctx), &inflationtypes.QueryParamsRequest{})
	if err != nil {
		return inflationtypes.Params{}, nil
	}

	return res.Params, nil
}

// CirculatingSupply implements inflationsource.Source
func (s Source) CirculatingSupply(height int64) (sdk.DecCoin, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return sdk.DecCoin{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.CirculatingSupply(sdk.WrapSDKContext(ctx), &inflationtypes.QueryCirculatingSupplyRequest{})
	if err != nil {
		return sdk.DecCoin{}, err
	}

	return res.CirculatingSupply, nil
}

// EpochMintProvision implements inflationsource.Source
func (s Source) EpochMintProvision(height int64) (sdk.DecCoin, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return sdk.DecCoin{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.EpochMintProvision(sdk.WrapSDKContext(ctx), &inflationtypes.QueryEpochMintProvisionRequest{})
	if err != nil {
		return sdk.DecCoin{}, err
	}

	return res.EpochMintProvision, nil
}

// InflationRate implements inflationsource.Source
func (s Source) InflationRate(height int64) (sdk.Dec, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return sdk.Dec{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.InflationRate(sdk.WrapSDKContext(ctx), &inflationtypes.QueryInflationRateRequest{})
	if err != nil {
		return sdk.Dec{}, err
	}

	return res.InflationRate, nil
}

// InflationPeriod implements inflationsource.Source
func (s Source) InflationPeriod(height int64) (uint64, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return 0, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.Period(sdk.WrapSDKContext(ctx), &inflationtypes.QueryPeriodRequest{})
	if err != nil {
		return 0, err
	}

	return res.Period, nil
}

// SkippedEpochs implements inflationsource.Source
func (s Source) SkippedEpochs(height int64) (uint64, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return 0, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.SkippedEpochs(sdk.WrapSDKContext(ctx), &inflationtypes.QuerySkippedEpochsRequest{})
	if err != nil {
		return 0, err
	}

	return res.SkippedEpochs, nil
}
