package gov

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/forbole/bdjuno/v2/database"

	govsource "github.com/forbole/bdjuno/v2/modules/gov/source"

	"github.com/forbole/juno/v2/modules"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.BlockModule   = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represent x/gov module
type Module struct {
	cdc            codec.Marshaler
	db             *database.Db
	source         govsource.Source
	authModule     AuthModule
	bankModule     BankModule
	distrModule    DistrModule
	mintModule     MintModule
	slashingModule SlashingModule
	stakingModule  StakingModule
}

// NewModule returns a new Module instance
func NewModule(
	source govsource.Source,
	authModule AuthModule,
	bankModule BankModule,
	distrModule DistrModule,
	mintModule MintModule,
	slashingModule SlashingModule,
	stakingModule StakingModule,
	cdc codec.Marshaler,
	db *database.Db,
) *Module {
	return &Module{
		cdc:            cdc,
		source:         source,
		authModule:     authModule,
		bankModule:     bankModule,
		distrModule:    distrModule,
		mintModule:     mintModule,
		slashingModule: slashingModule,
		stakingModule:  stakingModule,
		db:             db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "gov"
}
