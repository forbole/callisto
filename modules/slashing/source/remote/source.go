package remote

import (
	"github.com/cosmos/cosmos-sdk/types/query"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/desmos-labs/juno/v2/node/remote"

	slashingsource "github.com/forbole/bdjuno/modules/slashing/source"
)

var (
	_ slashingsource.Source = &Source{}
)

// Source implements slashingsource.Source using a remote node
type Source struct {
	*remote.Source
	querier slashingtypes.QueryClient
}

// NewSource returns a new Source implementation
func NewSource(source *remote.Source, querier slashingtypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetSigningInfos implements slashingsource.Source
func (s Source) GetSigningInfos(height int64) ([]slashingtypes.ValidatorSigningInfo, error) {
	header := remote.GetHeightRequestHeader(height)

	var signingInfos []slashingtypes.ValidatorSigningInfo
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.querier.SigningInfos(
			s.Ctx,
			&slashingtypes.QuerySigningInfosRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 1000, // Query 1000 signing infos at a time
				},
			},
			header,
		)
		if err != nil {
			return nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		signingInfos = append(signingInfos, res.Info...)
	}

	return signingInfos, nil
}

// GetParams implements slashingsource.Source
func (s Source) GetParams(height int64) (slashingtypes.Params, error) {
	res, err := s.querier.Params(s.Ctx, &slashingtypes.QueryParamsRequest{}, remote.GetHeightRequestHeader(height))
	if err != nil {
		return slashingtypes.Params{}, nil
	}

	return res.Params, nil
}
