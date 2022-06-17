package marker

import (
	"fmt"

	"github.com/forbole/bdjuno/v3/modules/utils"
	"github.com/forbole/bdjuno/v3/types"

	markertypes "github.com/MonikaCat/provenance/x/marker/types"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "marker").Msg("setting up periodic tasks")

	// Setup a cron job to run every hour
	if _, err := scheduler.Every(1).Hour().Do(func() {
		utils.WatchMethod(m.updateAllMarkers)
	}); err != nil {
		return err
	}

	return nil
}

// updateAllMarkers fetches from the REST APIs the latest markers
// and saves them inside the database.
func (m *Module) updateAllMarkers() error {
	log.Debug().
		Str("module", "marker").
		Msg("getting markers data")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get the markers list
	markersList, err := m.source.GetAllMarkers(height)
	if err != nil {
		return fmt.Errorf("error while getting markers list: %s", err)
	}

	var markers []types.Marker
	for _, marker := range markersList {
		var accountI markertypes.MarkerAccountI
		err := m.cdc.UnpackAny(marker, &accountI)
		if err != nil {
			return err
		}

		var supply []types.MarkerSupply
		supplyDenom, supplyAmount := accountI.GetSupplyValues()
		supply = append(supply, types.NewMarkerSupply(supplyDenom, supplyAmount.String()))

		markers = append(markers,
			*types.NewMarker(
				accountI.GetAddress().String(),
				accountI.GetAccessList(),
				accountI.HasGovernanceEnabled(),
				accountI.GetDenom(),
				accountI.GetMarkerType(),
				accountI.GetStatus(),
				supply,
				height))
	}

	return m.db.SaveMarkers(markers, height)
}
