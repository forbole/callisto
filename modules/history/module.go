package history

import (
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/modules/messages"
	juno "github.com/desmos-labs/juno/types"
	"github.com/go-co-op/gocron"

	"github.com/forbole/bdjuno/database"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
	_ modules.MessageModule            = &Module{}
)

// Module represents the module that allows to store historic information
type Module struct {
	messagesParser messages.MessageAddressesParser
	encodingConfig *params.EncodingConfig
	db             *database.Db
}

// NewModule allows to build a new Module instance
func NewModule(messagesParser messages.MessageAddressesParser, encodingConfig *params.EncodingConfig, db *database.Db) *Module {
	return &Module{
		messagesParser: messagesParser,
		encodingConfig: encodingConfig,
		db:             db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "history"
}

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	return RegisterPeriodicOps(scheduler, m.db)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, _ *juno.Tx) error {
	return HandleMsg(msg, m.messagesParser, m.encodingConfig.Marshaler, m.db)
}
