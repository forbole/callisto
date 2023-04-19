package remote

import (
	superfluidtypes "github.com/osmosis-labs/osmosis/v15/x/superfluid/types"

	"github.com/forbole/juno/v4/node/remote"

	superfluidsource "github.com/forbole/bdjuno/v4/modules/superfluid/source"
)

var (
	_ superfluidsource.Source = &Source{}
)

// Source implements superfluidsource.Source querying the data from a remote node
type Source struct {
	*remote.Source
	superfluidClient superfluidtypes.QueryClient
}

// NewSource returns a new Source instace
func NewSource(source *remote.Source, superfluidClient superfluidtypes.QueryClient) *Source {
	return &Source{
		Source:           source,
		superfluidClient: superfluidClient,
	}
}

// GetSuperfluidDelegationsByDelegator implements superfluidsource.Source
func (s Source) GetSuperfluidDelegationsByDelegator(address string, height int64) ([]superfluidtypes.SuperfluidDelegationRecord, error) {
	res, err := s.superfluidClient.SuperfluidDelegationsByDelegator(
		remote.GetHeightRequestContext(s.Ctx, height),
		&superfluidtypes.SuperfluidDelegationsByDelegatorRequest{DelegatorAddress: address},
	)
	if err != nil {
		return nil, err
	}

	return res.SuperfluidDelegationRecords, nil
}
