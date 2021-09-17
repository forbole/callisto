package issuer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/modules"
	juno "github.com/desmos-labs/juno/types"
	issuertypes "github.com/e-money/em-ledger/x/issuer/types"

	"github.com/forbole/bdjuno/database"
)

var (
	_ modules.Module        = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represent database/iscn module
type Module struct {
	inflationClient issuertypes.QueryClient
	db              *database.Db
}

// NewModule returns a new Module instance
func NewModule(issuerClient issuertypes.QueryClient, db *database.Db) *Module {
	return &Module{
		inflationClient: issuerClient,
		db:              db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "issuer"
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	return HandleMsg(tx, index, msg, m.inflationClient, m.db)
}
