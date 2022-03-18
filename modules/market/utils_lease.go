package market

import (
	"fmt"

	markettypes "github.com/ovrclk/akash/x/market/types/v1beta2"
)

func (m *Module) getLeases(height int64) ([]markettypes.QueryLeaseResponse, error) {

	leaseResponse, err := m.source.GetLeases(height)
	if err != nil {
		return nil, fmt.Errorf("error while getting lease responses: %s", err)
	}
	return leaseResponse, nil
}
