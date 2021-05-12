package distribution

import (
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/database"
	distrutils "github.com/forbole/bdjuno/modules/distribution/utils"
)

// HandleBlock represents a method that is called each time a new block is created
func HandleBlock(block *tmctypes.ResultBlock, client distrtypes.QueryClient, db *database.Db) error {
	// Update the validator commissions
	go distrutils.UpdateValidatorsCommissionAmounts(block.Block.Height, client, db)

	// Update the delegators commissions amounts
	go distrutils.UpdateDelegatorsRewardsAmounts(block.Block.Height, client, db)

	return nil
}
