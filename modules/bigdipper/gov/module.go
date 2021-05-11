package gov

import (
	"encoding/json"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"google.golang.org/grpc"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var _ modules.Module = &Module{}

// Module represent x/gov module
type Module struct {
	encodingConfig *params.EncodingConfig
	govClient      govtypes.QueryClient
	authClient     authtypes.QueryClient
	bankClient     banktypes.QueryClient
	db             *bigdipperdb.Db
}

// NewModule returns a new Module instance
func NewModule(encodingConfig *params.EncodingConfig, grpcConnection *grpc.ClientConn, db *bigdipperdb.Db) *Module {
	return &Module{
		encodingConfig: encodingConfig,
		govClient:      govtypes.NewQueryClient(grpcConnection),
		authClient:     authtypes.NewQueryClient(grpcConnection),
		bankClient:     banktypes.NewQueryClient(grpcConnection),
		db:             db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "gov"
}

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(_ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return HandleGenesis(appState, m.govClient, m.authClient, m.bankClient, m.encodingConfig.Marshaler, m.db)
}

// HandleMsg implements modules.Module
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *types.Tx) error {
	return HandleMsg(tx, msg, m.govClient, m.authClient, m.bankClient, m.encodingConfig.Marshaler, m.db)
}
