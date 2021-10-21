package remote

import (
	"fmt"

	"github.com/forbole/juno/v2/node/remote"
	iscntypes "github.com/likecoin/likechain/x/iscn/types"

	iscnsource "github.com/forbole/bdjuno/v2/modules/iscn/source"
)

var (
	_ iscnsource.Source = &Source{}
)

// Source implements iscnsource.Source using a remote node
type Source struct {
	*remote.Source
	client iscntypes.QueryClient
}

// NewSource allows to build a new Source implementation
func NewSource(source *remote.Source, client iscntypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: client,
	}
}

// GetParams implements iscnsource.Source
func (s *Source) GetParams(height int64) (iscntypes.Params, error) {
	res, err := s.client.Params(s.Ctx, &iscntypes.QueryParamsRequest{}, remote.GetHeightRequestHeader(height))
	if err != nil {
		return iscntypes.Params{}, fmt.Errorf("error while querying iscn params: %s", err)
	}

	return res.Params, nil
}

// GetRecordsByID implements iscnsource.Source
func (s *Source) GetRecordsByID(height int64, id string) (*iscntypes.QueryRecordsByIdResponse, error) {
	res, err := s.client.RecordsById(
		s.Ctx,
		&iscntypes.QueryRecordsByIdRequest{IscnId: id},
		remote.GetHeightRequestHeader(height),
	)
	if err != nil {
		return nil, fmt.Errorf("error while querying iscn records by id: %s", err)
	}

	return res, nil
}
