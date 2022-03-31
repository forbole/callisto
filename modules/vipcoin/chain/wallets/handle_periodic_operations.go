package wallets

import (
	"github.com/go-co-op/gocron"
)

// RegisterPeriodicOperations implements modules.Module
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	return nil
}

// updateSupply updates the supply of all the tokens
func (m *Module) updateSupply() error {
	return nil
}
