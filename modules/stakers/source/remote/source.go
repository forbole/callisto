package remote

import (
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/forbole/juno/v4/node/remote"

	stakerssource "github.com/forbole/bdjuno/v4/modules/stakers/source"
)

var (
	_ stakerssource.Source = &Source{}
)

// Source implements stakerssource.Source using a remote node
type Source struct {
	*remote.Source
	querier stakerstypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, querier stakerstypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// Params implements stakerssource.Source
func (s Source) Params(height int64) (stakerstypes.Params, error) {
	res, err := s.querier.Params(remote.GetHeightRequestContext(s.Ctx, height), &stakerstypes.QueryParamsRequest{})
	if err != nil {
		return stakerstypes.Params{}, nil
	}

	return res.Params, nil
}
