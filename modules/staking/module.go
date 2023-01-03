package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v4/modules"

	"github.com/forbole/bdjuno/v3/database"
	stakingsource "github.com/forbole/bdjuno/v3/modules/staking/source"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.BlockModule   = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/staking module
type Module struct {
	cdc            codec.Codec
	db             *database.Db
	source         stakingsource.Source
	slashingModule SlashingModule
}

// NewModule returns a new Module instance
func NewModule(
	source stakingsource.Source, slashingModule SlashingModule,
	cdc codec.Codec, db *database.Db,
) *Module {
	return &Module{
		cdc:            cdc,
		db:             db,
		source:         source,
		slashingModule: slashingModule,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "staking"
}
