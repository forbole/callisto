package distribution

import (
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"
	bdistrcommon "github.com/forbole/bdjuno/modules/bigdipper/distribution/common"
)

// HandleBlock represents a method that is called each time a new block is created
func HandleBlock(block *tmctypes.ResultBlock, client distrtypes.QueryClient, db *bigdipperdb.Db) error {
	// Update the validator commissions
	go bdistrcommon.UpdateValidatorsCommissionAmounts(block.Block.Height, client, db)

	// Update the delegators commissions amounts
	go bdistrcommon.UpdateDelegatorsRewardsAmounts(block.Block.Height, client, db)

	return nil
}
