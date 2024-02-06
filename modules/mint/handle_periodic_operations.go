package mint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/forbole/bdjuno/v4/modules/utils"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "mint").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight
	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.UpdateInflation)
	}); err != nil {
		return err
	}

	return nil
}

// updateInflation fetches from the REST APIs the latest value for the
// inflation, and saves it inside the database.
func (m *Module) UpdateInflation() error {
	log.Debug().
		Str("module", "mint").
		Str("operation", "inflation").
		Msg("getting inflation data")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get the inflation
	resp, err := http.Get("https://api.cosmoshub.forbole.com/cosmos/mint/v1beta1/inflation")
	if err != nil {
		return fmt.Errorf("error while querying API for inflation value: %s", err)
	}

	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error while reading response body: %s", err)
	}

	var inf minttypes.QueryInflationResponse
	err = json.Unmarshal(bz, &inf)
	if err != nil {
		return fmt.Errorf("error while unmarshaling response body: %s", err)
	}

	return m.db.SaveInflation(inf.Inflation, height)
}
