package remote

import (
	"fmt"

	inflationtypes "github.com/MonikaCat/em-ledger/x/inflation/types"
	"github.com/forbole/juno/v5/node/remote"

	inflationsource "github.com/forbole/bdjuno/v4/modules/inflation/source"
	"github.com/forbole/bdjuno/v4/utils"
)

var (
	_ inflationsource.Source = &Source{}
)

// Source implements inflationsource.Source using a remote node
type Source struct {
	*remote.Source
	client inflationtypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, client inflationtypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: client,
	}
}

// GetInflation implements inflationsource.Source
func (s *Source) GetInflation(height int64) (inflationtypes.InflationState, error) {
	res, err := s.client.Inflation(utils.GetHeightRequestContext(s.Ctx, height), &inflationtypes.QueryInflationRequest{})
	if err != nil {
		return inflationtypes.InflationState{}, fmt.Errorf("error while querying inflation state: %s", err)
	}

	return res.State, nil
}
