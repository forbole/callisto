package top_accounts

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/juno/v4/node"

	"github.com/forbole/juno/v4/modules"
	junomessages "github.com/forbole/juno/v4/modules/messages"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.MessageModule            = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represent x/top_accounts module
type Module struct {
	cdc           codec.Codec
	node          node.Node
	db            *database.Db
	messageParser junomessages.MessageAddressesParser
	authModule    AuthModule
	authSource    AuthSource
	bankModule    BankModule
	distrModule   DistrModule
	stakingModule StakingModule
}

// NewModule returns a new Module instance
func NewModule(
	authModule AuthModule,
	authSource AuthSource,
	bankModule BankModule,
	distrModule DistrModule,
	stakingModule StakingModule,
	messageParser junomessages.MessageAddressesParser,
	cdc codec.Codec,
	node node.Node,
	db *database.Db,
) *Module {
	return &Module{
		cdc:           cdc,
		node:          node,
		authModule:    authModule,
		authSource:    authSource,
		bankModule:    bankModule,
		distrModule:   distrModule,
		messageParser: messageParser,
		stakingModule: stakingModule,
		db:            db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "top_accounts"
}
