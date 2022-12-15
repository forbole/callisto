package topaccounts

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/forbole/bdjuno/v3/database"

	govsource "github.com/forbole/bdjuno/v3/modules/gov/source"

	"github.com/forbole/juno/v3/modules"
)

var (
	_ modules.Module        = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represent x/gov module
type Module struct {
	cdc           codec.Codec
	db            *database.Db
	source        govsource.Source
	bankModule    BankModule
	distrModule   DistrModule
	stakingModule StakingModule
}

// NewModule returns a new Module instance
func NewModule(
	source govsource.Source,
	bankModule BankModule,
	distrModule DistrModule,
	stakingModule StakingModule,
	cdc codec.Codec,
	db *database.Db,
) *Module {
	return &Module{
		cdc:           cdc,
		source:        source,
		bankModule:    bankModule,
		distrModule:   distrModule,
		stakingModule: stakingModule,
		db:            db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "gov"
}
