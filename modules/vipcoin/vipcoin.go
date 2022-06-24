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
	vipcoinaccountssource "github.com/forbole/bdjuno/v2/modules/vipcoin/chain/accounts/source"
	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/assets"
	vipcoinassetssource "github.com/forbole/bdjuno/v2/modules/vipcoin/chain/assets/source"
	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/banking"
	vipcoinbankingsource "github.com/forbole/bdjuno/v2/modules/vipcoin/chain/banking/source"
	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/wallets"
	vipcoinwalletssource "github.com/forbole/bdjuno/v2/modules/vipcoin/chain/wallets/source"
)

var (
	_ modules.Module        = &module{}
	_ modules.GenesisModule = &module{}
)

type vipcoinModule interface {
	jmodules.Module
	jmodules.GenesisModule
	jmodules.MessageModule
}

type module struct {
	cdc            codec.Marshaler
	db             *database.Db
	lastBlockRepo  last_block.Repository
	logger         logging.Logger
	vipcoinModules []vipcoinModule

	schedulerRun bool
	mutex        sync.RWMutex
}

func NewModule(
	cdc codec.Marshaler,
	db *database.Db,
	logger logging.Logger,

	VipcoinAccountsSource vipcoinaccountssource.Source,
	VipcoinWalletsSource vipcoinwalletssource.Source,
	VipcoinBankingSource vipcoinbankingsource.Source,
	VipcoinAssetsSource vipcoinassetssource.Source,
) *module {
	module := &module{
		cdc:           cdc,
		db:            db,
		lastBlockRepo: *last_block.NewRepository(db.Sqlx),
		logger:        logger,
		vipcoinModules: []vipcoinModule{
			accounts.NewModule(VipcoinAccountsSource, cdc, db),
			assets.NewModule(VipcoinAssetsSource, cdc, db),
			banking.NewModule(VipcoinBankingSource, cdc, db),
			wallets.NewModule(VipcoinWalletsSource, cdc, db),
		},
	}

	go module.scheduler()

	return module
}

// Name implements modules.Module
func (m *module) Name() string {
	return "vipcoin"
}
