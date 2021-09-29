package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/juno/modules"

	"github.com/forbole/bdjuno/database"
	stakingsource "github.com/forbole/bdjuno/modules/staking/source"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.BlockModule   = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/staking module
type Module struct {
	cdc           codec.Marshaler
	db            *database.Db
	source        stakingsource.Source
	bankModule    BankModule
	distrModule   DistrModule
	historyModule HistoryModule
}

// NewModule returns a new Module instance
func NewModule(
	source stakingsource.Source,
	bankModule BankModule, distrModule DistrModule, historyModule HistoryModule,
	cdc codec.Marshaler, db *database.Db,
) *Module {
	return &Module{
		cdc:           cdc,
		db:            db,
		source:        source,
		bankModule:    bankModule,
		distrModule:   distrModule,
		historyModule: historyModule,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "staking"
}
