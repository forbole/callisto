package auth

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/desmos-labs/juno/modules/messages"
	juno "github.com/desmos-labs/juno/types"
	"google.golang.org/grpc"

	"github.com/desmos-labs/juno/modules"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/database"
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
	authClient     authtypes.QueryClient
	db             *database.BigDipperDb
}

// NewModule builds a new Module instance
func NewModule(
	messagesParser messages.MessageAddressesParser,
	encodingConfig *params.EncodingConfig, grpcConnection *grpc.ClientConn, db *database.BigDipperDb,
) *Module {
	return &Module{
		messagesParser: messagesParser,
		encodingConfig: encodingConfig,
		authClient:     authtypes.NewQueryClient(grpcConnection),
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
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
	return HandleMsg(tx, msg, m.messagesParser, m.authClient, m.encodingConfig.Marshaler, m.db)
}
