package overgold

import (
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v3/logging"
	"github.com/forbole/juno/v3/modules"
	jmodules "github.com/forbole/juno/v3/modules"

	"github.com/forbole/juno/v3/node"

	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/database/overgold/chain/last_block"
	"github.com/forbole/bdjuno/v3/modules/overgold/chain/accounts"
	overgoldAccountsSource "github.com/forbole/bdjuno/v3/modules/overgold/chain/accounts/source"
	"github.com/forbole/bdjuno/v3/modules/overgold/chain/assets"
	overgoldAssetsSource "github.com/forbole/bdjuno/v3/modules/overgold/chain/assets/source"
	"github.com/forbole/bdjuno/v3/modules/overgold/chain/banking"
	overgoldBankingSource "github.com/forbole/bdjuno/v3/modules/overgold/chain/banking/source"
	"github.com/forbole/bdjuno/v3/modules/overgold/chain/wallets"
	overgoldWalletsSource "github.com/forbole/bdjuno/v3/modules/overgold/chain/wallets/source"
)

var (
	_ modules.Module        = &module{}
	_ modules.GenesisModule = &module{}
)

type overgoldModule interface {
	jmodules.Module
	jmodules.GenesisModule
	jmodules.MessageModule
}

type module struct {
	cdc             codec.Codec
	db              *database.Db
	lastBlockRepo   last_block.Repository
	logger          logging.Logger
	overgoldModules []overgoldModule
	node            node.Node

	schedulerRun bool
	mutex        sync.RWMutex
}

func NewModule(
	cdc codec.Codec,
	db *database.Db,
	node node.Node,
	logger logging.Logger,

	OvergoldAccountsSource overgoldAccountsSource.Source,
	OvergoldWalletsSource overgoldWalletsSource.Source,
	OvergoldBankingSource overgoldBankingSource.Source,
	OvergoldAssetsSource overgoldAssetsSource.Source,
) *module {
	module := &module{
		cdc:           cdc,
		db:            db,
		lastBlockRepo: *last_block.NewRepository(db.Sqlx),
		node:          node,
		logger:        logger,
		overgoldModules: []overgoldModule{
			accounts.NewModule(OvergoldAccountsSource, cdc, db),
			assets.NewModule(OvergoldAssetsSource, cdc, db),
			banking.NewModule(OvergoldBankingSource, cdc, db),
			wallets.NewModule(OvergoldWalletsSource, cdc, db),
		},
	}

	go module.scheduler()

	return module
}

// Name implements modules.Module
func (m *module) Name() string {
	return "overgold"
}
