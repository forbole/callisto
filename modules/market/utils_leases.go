package market

import (
	"fmt"

	"github.com/forbole/bdjuno/v3/types"
	"github.com/rs/zerolog/log"

	markettypes "github.com/ovrclk/akash/x/market/types/v1beta2"
)

// updateLeases fetch all the leases with latest statuses and store them into database
func (m *Module) updateLeases(height int64) error {
	log.Debug().Str("module", "market").Int64("height", height).Msg("getting leases")

	leasesResponse, err := m.source.GetActiveLeases(height)
	if err != nil {
		return fmt.Errorf("error while getting akash leases: %s", err)
	}

	leases := m.convertLeaseFromResponse(leasesResponse, height)

	// Store the leases into the database
	return m.db.SaveLeases(leases, height)
}

func (m *Module) convertLeaseFromResponse(res []markettypes.QueryLeaseResponse, height int64) []*types.MarketLease {
	leases := make([]*types.MarketLease, len(res))
	for i, r := range res {
		leases[i] = types.NewMarketLease(r, height)
	}

	return leases
}
