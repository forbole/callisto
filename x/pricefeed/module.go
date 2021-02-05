package pricefeed

import (
	"github.com/desmos-labs/juno/modules"
	"github.com/go-co-op/gocron"

	"github.com/forbole/bdjuno/database"
)

var _ modules.Module = &Module{}

// Module represent x/pricefeed module
type Module struct {
	db *database.BigDipperDb
}

// NewModule returns a new Module instance
func NewModule(db *database.BigDipperDb) *Module {
	return &Module{
		db: db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "pricefeed"
}

// RegisterPeriodicOperations implements PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	return RegisterPeriodicOps(scheduler, m.db)
}
