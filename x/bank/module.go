package bank

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/client"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	"github.com/go-co-op/gocron"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/database"
)

var _ modules.Module = Module{}

// Module represents the x/bank module
type Module struct{}

// Name implements modules.Module
func (m Module) Name() string {
	return "bank"
}

// RegisterPeriodicOperations implements modules.Module
func (m Module) RegisterPeriodicOperations(*gocron.Scheduler, *codec.Codec, *client.Proxy, db.Database) error {
	return nil
}

// RunAdditionalOperations implements modules.Module
func (m Module) RunAdditionalOperations(*config.Config, *codec.Codec, *client.Proxy, db.Database) error {
	return nil
}

// HandleGenesis implements modules.Module
func (m Module) HandleGenesis(
	*tmtypes.GenesisDoc, map[string]json.RawMessage, *codec.Codec, *client.Proxy, db.Database,
) error {
	return nil
}

// HandleBlock implements modules.Module
func (m Module) HandleBlock(
	*tmctypes.ResultBlock, []types.Tx, *tmctypes.ResultValidators, *codec.Codec, *client.Proxy, db.Database,
) error {
	return nil
}

// HandleTx implements modules.Module
func (m Module) HandleTx(types.Tx, *codec.Codec, *client.Proxy, db.Database) error {
	return nil
}

// HandleMsg implements modules.Module
func (m Module) HandleMsg(index int, sdkMsg sdk.Msg, tx types.Tx, _ *codec.Codec, cp *client.Proxy, db db.Database) error {
	bdDatabase := database.Cast(db)
	return Handler(tx, index, sdkMsg, cp, bdDatabase)
}
