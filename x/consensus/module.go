package consensus

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/client"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
	"github.com/go-co-op/gocron"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var _ modules.Module = Module{}

// Module implements the consensus operations
type Module struct{}

// Name implements modules.Module
func (m Module) Name() string {
	return "consensus"
}

// RegisterPeriodicOperations implements modules.Module
func (m Module) RegisterPeriodicOperations(
	scheduler *gocron.Scheduler, _ *codec.Codec, cp *client.Proxy, db db.Database,
) error {
	bdDatabase := database.Cast(db)
	return Register(scheduler, cp, bdDatabase)
}

// RunAdditionalOperations implements modules.Module
func (m Module) RunAdditionalOperations(_ *config.Config, _ *codec.Codec, cp *client.Proxy, db db.Database) error {
	bdDatabase := database.Cast(db)
	return ListenOperation(cp, bdDatabase)
}

// HandleGenesis implements modules.Module
func (m Module) HandleGenesis(
	doc *tmtypes.GenesisDoc, _ map[string]json.RawMessage, _ *codec.Codec, _ *client.Proxy, db db.Database,
) error {
	bdDatabase := database.Cast(db)
	return HandleGenesis(doc, bdDatabase)
}

// HandleBlock implements modules.Module
func (m Module) HandleBlock(
	b *tmctypes.ResultBlock, _ []types.Tx, _ *tmctypes.ResultValidators, _ *codec.Codec, _ *client.Proxy, db db.Database,
) error {
	bdDatabase := database.Cast(db)
	return HandleBlock(b, bdDatabase)
}

// HandleTx implements modules.Module
func (m Module) HandleTx(types.Tx, *codec.Codec, *client.Proxy, db.Database) error {
	return nil
}

// HandleMsg implements modules.Module
func (m Module) HandleMsg(int, sdk.Msg, types.Tx, *codec.Codec, *client.Proxy, db.Database) error {
	return nil
}
