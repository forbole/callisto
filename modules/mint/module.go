package mint

import (
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/desmos-labs/juno/modules"
	"github.com/go-co-op/gocron"
	"google.golang.org/grpc"

	"github.com/forbole/bdjuno/database"
)

var _ modules.Module = &Module{}

// Module represent database/mint module
type Module struct {
	mintClient minttypes.QueryClient
	db         *database.Db
}

// NewModule returns a new Module instance
func NewModule(grpcConnection *grpc.ClientConn, db *database.Db) *Module {
	return &Module{
		mintClient: minttypes.NewQueryClient(grpcConnection),
		db:         db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "mint"
}

// RegisterPeriodicOperations implements PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	return RegisterPeriodicOps(scheduler, m.mintClient, m.db)
}
