package market

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// updateLeases fetch all the leases with latest statuses and store them into database
func (m *Module) updateLeases(height int64) error {
	log.Debug().Str("module", "market").Int64("height", height).Msg("getting leases")

	leases, err := m.source.GetLeases(height)
	if err != nil {
		return fmt.Errorf("error while getting akash leases: %s", err)
	}

	// Store the leases into the database
	return m.db.SaveLeases(leases, height)
}
