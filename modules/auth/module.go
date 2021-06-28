package auth

import (
	"encoding/json"

	"github.com/forbole/bdjuno/database"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/modules/messages"
	juno "github.com/desmos-labs/juno/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	_ modules.Module         = &Module{}
	_ modules.GenesisModule  = &Module{}
	_ modules.MessageModule  = &Module{}
	_ modules.FastSyncModule = &Module{}
)

// Module represents the x/auth module
type Module struct {
	messagesParser messages.MessageAddressesParser
	encodingConfig *params.EncodingConfig
	authClient     authttypes.QueryClient
	db             *database.Db
}

// NewModule builds a new Module instance
func NewModule(
	messagesParser messages.MessageAddressesParser,
	authClient authttypes.QueryClient,
	encodingConfig *params.EncodingConfig, db *database.Db,
) *Module {
	return &Module{
		messagesParser: messagesParser,
		encodingConfig: encodingConfig,
		authClient:     authClient,
		db:             db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "auth"
}

// DownloadState implements modules.FastSyncModule
func (m *Module) DownloadState(height int64) error {
	return FastSync(height, m.authClient, m.db)
}

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(_ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return Handler(appState, m.encodingConfig.Marshaler, m.db)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, _ *juno.Tx) error {
	return HandleMsg(msg, m.messagesParser, m.encodingConfig.Marshaler, m.db)
}
