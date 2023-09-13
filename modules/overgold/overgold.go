package overgold

import (
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v5/logging"
	"github.com/forbole/juno/v5/modules"
	jmodules "github.com/forbole/juno/v5/modules"

	"github.com/forbole/juno/v5/node"

	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/database/overgold/chain/last_block"
	"github.com/forbole/bdjuno/v4/modules/overgold/chain/allowed"
	overgoldAllowedSource "github.com/forbole/bdjuno/v4/modules/overgold/chain/allowed/source"
	"github.com/forbole/bdjuno/v4/modules/overgold/chain/core"
	overgoldCoreSource "github.com/forbole/bdjuno/v4/modules/overgold/chain/core/source"
	"github.com/forbole/bdjuno/v4/modules/overgold/chain/feeexcluder"
	overgoldFeeExcluderSource "github.com/forbole/bdjuno/v4/modules/overgold/chain/feeexcluder/source"
	"github.com/forbole/bdjuno/v4/modules/overgold/chain/referral"
	overgoldReferralSource "github.com/forbole/bdjuno/v4/modules/overgold/chain/referral/source"
	"github.com/forbole/bdjuno/v4/modules/overgold/chain/stake"
	overgoldStakeSource "github.com/forbole/bdjuno/v4/modules/overgold/chain/stake/source"
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

	OverGoldCoreSource overgoldCoreSource.Source,
	OverGoldAllowedSource overgoldAllowedSource.Source,
	OverGoldFeeExcluderSource overgoldFeeExcluderSource.Source,
	OverGoldReferralSource overgoldReferralSource.Source,
	OverGoldStakeSource overgoldStakeSource.Source,
) *module {
	module := &module{
		cdc:           cdc,
		db:            db,
		lastBlockRepo: *last_block.NewRepository(db.Sqlx),
		node:          node,
		logger:        logger,
		overgoldModules: []overgoldModule{
			core.NewModule(OverGoldCoreSource, cdc, db),
			allowed.NewModule(OverGoldAllowedSource, cdc, db),
			feeexcluder.NewModule(OverGoldFeeExcluderSource, cdc, db),
			referral.NewModule(OverGoldReferralSource, cdc, db),
			stake.NewModule(OverGoldStakeSource, cdc, db),
		},
	}

	go module.scheduler()

	return module
}

// Name implements modules.Module
func (m *module) Name() string {
	return "overgold"
}
