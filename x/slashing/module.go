package slashing

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

// Module represent x/slashing module
type Module struct{}

// Name implements modules.Module
func (m Module) Name() string {
	return "slashing"
}

// RunAdditionalOperations implements modules.Module
func (m Module) RunAdditionalOperations(cfg *config.Config, cdc *codec.Codec, cp *client.Proxy, db db.Database) error {
	return nil
}

// RegisterPeriodicOperations implements modules.Module
func (m Module) RegisterPeriodicOperations(
	scheduler *gocron.Scheduler, cdc *codec.Codec, cp *client.Proxy, db db.Database,
) error {
	return nil
}

// HandleGenesis implements modules.Module
func (m Module) HandleGenesis(
	doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, cdc *codec.Codec, cp *client.Proxy, db db.Database,
) error {
	return nil
}

// HandleBlock implements modules.Module
func (m Module) HandleBlock(
	block *tmctypes.ResultBlock, txs []types.Tx, vals *tmctypes.ResultValidators,
	cdc *codec.Codec, cp *client.Proxy, db db.Database,
) error {
	bdDatabase := database.Cast(db)
	return HandleBlock(block, cp, bdDatabase)
}

// HandleTx implements modules.Module
func (m Module) HandleTx(tx types.Tx, cdc *codec.Codec, cp *client.Proxy, db db.Database) error {
	return nil
}

// HandleMsg implements modules.Module
func (m Module) HandleMsg(
	index int, msg sdk.Msg, tx types.Tx, cdc *codec.Codec, cp *client.Proxy, db db.Database,
) error {
	return nil
}
