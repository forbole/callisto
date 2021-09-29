package gov

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/forbole/bdjuno/v2/database"

	govsource "github.com/forbole/bdjuno/v2/modules/gov/source"

	"github.com/desmos-labs/juno/v2/modules"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.BlockModule   = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represent x/gov module
type Module struct {
	cdc           codec.Marshaler
	db            *database.Db
	source        govsource.Source
	authModule    AuthModule
	bankModule    BankModule
	stakingModule StakingModule
}

// NewModule returns a new Module instance
func NewModule(
	cdc codec.Marshaler, source govsource.Source,
	authModule AuthModule, bankModule BankModule, stakingModule StakingModule,
	db *database.Db,
) *Module {
	return &Module{
		cdc:           cdc,
		db:            db,
		source:        source,
		authModule:    authModule,
		bankModule:    bankModule,
		stakingModule: stakingModule,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "gov"
}
