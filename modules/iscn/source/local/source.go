package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v3/node/local"
	iscntypes "github.com/likecoin/likechain/x/iscn/types"

	iscnsource "github.com/forbole/bdjuno/v2/modules/iscn/source"
)

var (
	_ iscnsource.Source = &Source{}
)

// Source implements iscnsource.Source using a local node
type Source struct {
	*local.Source
	client iscntypes.QueryServer
}

// NewSource allows to build a new Source implementation
func NewSource(source *local.Source, client iscntypes.QueryServer) *Source {
	return &Source{
		Source: source,
		client: client,
	}
}

// GetParams implements iscnsource.Source
func (s *Source) GetParams(height int64) (iscntypes.Params, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return iscntypes.Params{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.client.Params(sdk.WrapSDKContext(ctx), &iscntypes.QueryParamsRequest{})
	if err != nil {
		return iscntypes.Params{}, fmt.Errorf("error while reading iscn params: %s", err)
	}

	return res.Params, nil
}

// GetRecordsByID implements iscnsource.Source
func (s *Source) GetRecordsByID(height int64, id string) (*iscntypes.QueryRecordsByIdResponse, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.client.RecordsById(sdk.WrapSDKContext(ctx), &iscntypes.QueryRecordsByIdRequest{IscnId: id})
	if err != nil {
		return nil, fmt.Errorf("error while reading iscn records by id: %s", err)
	}

	return res, nil
}
