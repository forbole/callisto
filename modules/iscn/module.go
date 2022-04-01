package iscn

import (
	"github.com/forbole/juno/v3/modules"

	"github.com/forbole/bdjuno/v2/database"
	iscnsource "github.com/forbole/bdjuno/v2/modules/iscn/source"
)

var (
	_ modules.Module      = &Module{}
	_ modules.BlockModule = &Module{}
)

// Module represent database/iscn module
type Module struct {
	db     *database.Db
	source iscnsource.Source
}

// NewModule returns a new Module instance
func NewModule(source iscnsource.Source, db *database.Db) *Module {
	return &Module{
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "iscn"
}
