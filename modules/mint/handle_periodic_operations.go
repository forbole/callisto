package mint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/forbole/bdjuno/v4/modules/utils"
	"github.com/forbole/bdjuno/v4/types"

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

	inflation, err := queryLatestInflation()
	if err != nil {
		return err
	}

	return m.db.SaveInflation(inflation, height)

}

// queryLatestInflation queries cheqd latest inflation value
func queryLatestInflation() (string, error) {
	resp, err := http.Get("https://api.cheqd.net/cosmos/mint/v1beta1/inflation")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading inflation response body: %s", err)
	}

	var inflation types.InflationCheqd
	err = json.Unmarshal(bz, &inflation)
	if err != nil {
		return "", fmt.Errorf("error while unmarshaling inflation response body: %s", err)
	}

	return inflation.Inflation, nil
}
