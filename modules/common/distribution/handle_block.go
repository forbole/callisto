package distribution

import (
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock represents a method that is called each time a new block is created
func HandleBlock(block *tmctypes.ResultBlock, client distrtypes.QueryClient, db DB) error {
	// Update the validator commissions
	go UpdateValidatorsCommissionAmounts(block.Block.Height, client, db)

	// Update the delegators commissions amounts
	go UpdateDelegatorsRewardsAmounts(block.Block.Height, client, db)

	return nil
}
