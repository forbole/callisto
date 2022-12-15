package vipcoin

import (
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v2/logging"
	"github.com/forbole/juno/v2/modules"
	jmodules "github.com/forbole/juno/v2/modules"

	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/database/vipcoin/chain/last_block"
	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/accounts"
	overgoldAccountsSource "github.com/forbole/bdjuno/v2/modules/vipcoin/chain/accounts/source"
	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/assets"
	overgoldAssetsSource "github.com/forbole/bdjuno/v2/modules/vipcoin/chain/assets/source"
	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/banking"
	overgoldBankingSource "github.com/forbole/bdjuno/v2/modules/vipcoin/chain/banking/source"
	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/wallets"
	overgoldWalletsSource "github.com/forbole/bdjuno/v2/modules/vipcoin/chain/wallets/source"
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
	cdc             codec.Marshaler
	db              *database.Db
	lastBlockRepo   last_block.Repository
	logger          logging.Logger
	overgoldModules []overgoldModule

	schedulerRun bool
	mutex        sync.RWMutex
}

func NewModule(
	cdc codec.Marshaler,
	db *database.Db,
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
