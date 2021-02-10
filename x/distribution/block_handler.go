package distribution

import (
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/x/distribution/common"

	"github.com/forbole/bdjuno/database"
)

// HandleBlock represents a method that is called each time a new block is created
func HandleBlock(block *tmctypes.ResultBlock, client distrtypes.QueryClient, db *database.BigDipperDb) error {
	// Update the validator commissions
	go common.UpdateValidatorsCommissionAmounts(block.Block.Height, client, db)

	// Update the delegators commissions amounts
	go common.UpdateDelegatorsRewardsAmounts(block.Block.Height, client, db)

	return nil
}
