package history

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/modules/messages"
	juno "github.com/desmos-labs/juno/types"
	"github.com/go-co-op/gocron"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
	_ modules.GenesisModule            = &Module{}
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
	return types.HistoryModuleName
}

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	return RegisterPeriodicOps(scheduler, m.db)
}

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, _ map[string]json.RawMessage) error {
	return HandleGenesis(doc.GenesisTime, m.db)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
	return HandleMsg(tx, msg, m.messagesParser, m.encodingConfig.Marshaler, m.db)
}
