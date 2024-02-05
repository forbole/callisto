package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/forbole/juno/v5/node/local"

	slashingsource "github.com/forbole/callisto/v4/modules/slashing/source"
)

var (
	_ slashingsource.Source = &Source{}
)

// Source implements slashingsource.Source using a local node
type Source struct {
	*local.Source
	querier slashingtypes.QueryServer
}

// NewSource implements a new Source instance
func NewSource(source *local.Source, querier slashingtypes.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetSigningInfos implements slashingsource.Source
func (s Source) GetSigningInfos(height int64) ([]slashingtypes.ValidatorSigningInfo, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var signingInfos []slashingtypes.ValidatorSigningInfo
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.querier.SigningInfos(
			sdk.WrapSDKContext(ctx),
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
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return slashingtypes.Params{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.Params(sdk.WrapSDKContext(ctx), &slashingtypes.QueryParamsRequest{})
	if err != nil {
		return slashingtypes.Params{}, nil
	}

	return res.Params, nil
}

// GetSigningInfo implements slashingsource.GetSigningInfo
func (s Source) GetSigningInfo(height int64, consAddr sdk.ConsAddress) (slashingtypes.ValidatorSigningInfo, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return slashingtypes.ValidatorSigningInfo{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.SigningInfo(
		sdk.WrapSDKContext(ctx),
		&slashingtypes.QuerySigningInfoRequest{
			ConsAddress: consAddr.String(),
		},
	)

	if err != nil {
		return slashingtypes.ValidatorSigningInfo{}, err
	}

	return res.ValSigningInfo, nil
}
