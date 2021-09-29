package history

import (
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/modules/utils"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "history").Msg("setting up periodic tasks")

	// Update the historic balance of users every 10 minutes
	if _, err := scheduler.Every(30).Minutes().Do(func() {
		utils.WatchMethod(m.updateHistoricBalances)
	}); err != nil {
		return err
	}

	return nil
}

// updateHistoricBalances updates the historic balances of all the users
func (m *Module) updateHistoricBalances() error {
	log.Debug().Str("module", "history").Msg("updating historic balances")

	accounts, err := m.db.GetAccounts()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		err = m.UpdateAccountBalanceHistory(account)
		if err != nil {
			return err
		}
	}

	return nil
}
