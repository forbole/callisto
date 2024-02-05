package remote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/query"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/forbole/juno/v5/node/remote"

	slashingsource "github.com/forbole/callisto/v4/modules/slashing/source"
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
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var signingInfos []slashingtypes.ValidatorSigningInfo
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.querier.SigningInfos(
			ctx,
			&slashingtypes.QuerySigningInfosRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 1000, // Query 1000 signing infos at a time
				},
			},
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
	res, err := s.querier.Params(remote.GetHeightRequestContext(s.Ctx, height), &slashingtypes.QueryParamsRequest{})
	if err != nil {
		return slashingtypes.Params{}, nil
	}

	return res.Params, nil
}

// GetSigningInfo implements slashingsource.GetSigningInfo
func (s Source) GetSigningInfo(height int64, consAddr sdk.ConsAddress) (slashingtypes.ValidatorSigningInfo, error) {
	res, err := s.querier.SigningInfo(
		s.Ctx,
		&slashingtypes.QuerySigningInfoRequest{
			ConsAddress: consAddr.String(),
		},
	)

	if err != nil {
		return slashingtypes.ValidatorSigningInfo{}, err
	}
	return res.ValSigningInfo, nil
}
