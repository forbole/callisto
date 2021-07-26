package history

import (
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	historyutils "github.com/forbole/bdjuno/modules/history/utils"
	"github.com/forbole/bdjuno/modules/utils"
)

// RegisterPeriodicOps returns the operations that periodically runs and updates the users historic balance
func RegisterPeriodicOps(scheduler *gocron.Scheduler, db *database.Db) error {
	log.Debug().Str("module", "history").Msg("setting up periodic tasks")

	// Update the historic balance of users every 10 minutes
	if _, err := scheduler.Every(30).Minutes().StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return updateHistoricBalances(db) })
	}); err != nil {
		return err
	}

	return nil
}

// updateHistoricBalances updates the historic balances of all the users
func updateHistoricBalances(db *database.Db) error {
	log.Debug().Str("module", "history").Msg("updating historic balances")

	accounts, err := db.GetAccounts()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		err = historyutils.UpdateAccountBalanceHistory(account, db)
		if err != nil {
			return err
		}
	}

	return nil
}
