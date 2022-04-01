package remote

import (
	"fmt"

	"github.com/forbole/juno/v3/node/remote"
	iscntypes "github.com/likecoin/likechain/x/iscn/types"

	iscnsource "github.com/forbole/bdjuno/v2/modules/iscn/source"
	"github.com/forbole/bdjuno/v2/utils"
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
	res, err := s.client.Params(utils.GetHeightRequestContext(s.Ctx, height), &iscntypes.QueryParamsRequest{})
	if err != nil {
		return iscntypes.Params{}, fmt.Errorf("error while querying iscn params: %s", err)
	}

	return res.Params, nil
}

// GetRecordsByID implements iscnsource.Source
func (s *Source) GetRecordsByID(height int64, id string) (*iscntypes.QueryRecordsByIdResponse, error) {
	res, err := s.client.RecordsById(
		utils.GetHeightRequestContext(s.Ctx, height),
		&iscntypes.QueryRecordsByIdRequest{IscnId: id},
	)
	if err != nil {
		return nil, fmt.Errorf("error while querying iscn records by id: %s", err)
	}

	return res, nil
}
