package remote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	inflationtypes "github.com/evmos/evmos/v15/x/inflation/types"
	"github.com/forbole/juno/v5/node/remote"

	inflationsource "github.com/forbole/bdjuno/v4/modules/inflation/source"
)

var (
	_ inflationsource.Source = &Source{}
)

// Source implements inflationsource.Source using a remote node
type Source struct {
	*remote.Source
	querier inflationtypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, querier inflationtypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// Params implements inflationsource.Source
func (s Source) Params(height int64) (inflationtypes.Params, error) {
	res, err := s.querier.Params(remote.GetHeightRequestContext(s.Ctx, height), &inflationtypes.QueryParamsRequest{})
	if err != nil {
		return inflationtypes.Params{}, nil
	}

	return res.Params, nil
}

// CirculatingSupply implements inflationsource.Source
func (s Source) CirculatingSupply(height int64) (sdk.DecCoin, error) {
	res, err := s.querier.CirculatingSupply(remote.GetHeightRequestContext(s.Ctx, height), &inflationtypes.QueryCirculatingSupplyRequest{})
	if err != nil {
		return sdk.DecCoin{}, err
	}

	return res.CirculatingSupply, nil
}

// EpochMintProvision implements inflationsource.Source
func (s Source) EpochMintProvision(height int64) (sdk.DecCoin, error) {
	res, err := s.querier.EpochMintProvision(remote.GetHeightRequestContext(s.Ctx, height), &inflationtypes.QueryEpochMintProvisionRequest{})
	if err != nil {
		return sdk.DecCoin{}, err
	}

	return res.EpochMintProvision, nil
}

// InflationRate implements inflationsource.Source
func (s Source) InflationRate(height int64) (sdk.Dec, error) {
	res, err := s.querier.InflationRate(remote.GetHeightRequestContext(s.Ctx, height), &inflationtypes.QueryInflationRateRequest{})
	if err != nil {
		return sdk.Dec{}, err
	}

	return res.InflationRate, nil
}

// InflationPeriod implements inflationsource.Source
func (s Source) InflationPeriod(height int64) (uint64, error) {
	res, err := s.querier.Period(remote.GetHeightRequestContext(s.Ctx, height), &inflationtypes.QueryPeriodRequest{})
	if err != nil {
		return 0, err
	}

	return res.Period, nil
}

// SkippedEpochs implements inflationsource.Source
func (s Source) SkippedEpochs(height int64) (uint64, error) {
	res, err := s.querier.SkippedEpochs(remote.GetHeightRequestContext(s.Ctx, height), &inflationtypes.QuerySkippedEpochsRequest{})
	if err != nil {
		return 0, err
	}

	return res.SkippedEpochs, nil
}
