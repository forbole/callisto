package remote

import (
	"fmt"

	oracletypes "github.com/bandprotocol/chain/v2/x/oracle/types"
	"github.com/desmos-labs/juno/v2/node/remote"

	oraclesource "github.com/forbole/bdjuno/v2/modules/oracle/source"
)

var (
	_ oraclesource.Source = &Source{}
)

// Source implements oraclesource.Source based on a remote node
type Source struct {
	*remote.Source
	client oracletypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, client oracletypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: client,
	}
}

// GetParams implements oraclesource.Source
func (s *Source) GetParams(height int64) (oracletypes.Params, error) {
	res, err := s.client.Params(
		s.Ctx,
		&oracletypes.QueryParamsRequest{},
		remote.GetHeightRequestHeader(height),
	)
	if err != nil {
		return oracletypes.Params{}, fmt.Errorf("error while getting oracle params: %s", err)
	}

	return res.Params, nil
}
