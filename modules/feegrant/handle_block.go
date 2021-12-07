package feegrant

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	feegranttypes "github.com/cosmos/cosmos-sdk/x/feegrant"
	juno "github.com/forbole/juno/v2/types"

	"github.com/forbole/bdjuno/v2/types"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Tx, vals *tmctypes.ResultValidators,
) error {

	// Remove expired fee grant allowances
	err := m.removeExpiredFeeGrantAllowances(block.Block.Height, res.EndBlockEvents)
	if err != nil {
		fmt.Printf("Error when removing expired fee grant allowance, error: %s", err)
	}
	return nil
}

// removeExpiredFeeGrantAllowances removes fee grant allowances in database that have expired
func (m *Module) removeExpiredFeeGrantAllowances(height int64, events []abci.Event) error {
	log.Debug().Str("module", "feegrant").Int64("height", height).
		Msg("updating expired fee grant allowances")

	events = juno.FindEventsByType(events, feegranttypes.EventTypeRevokeFeeGrant)

	for _, event := range events {
		granterAddress, err := juno.FindAttributeByKey(event, feegranttypes.AttributeKeyGranter)
		if err != nil {
			return fmt.Errorf("error while getting fee grant granter address: %s", err)
		}
		granteeAddress, err := juno.FindAttributeByKey(event, feegranttypes.AttributeKeyGrantee)
		if err != nil {
			return fmt.Errorf("error while getting fee grant grantee address: %s", err)
		}
		allowanceToRemove := types.NewGrantRemoval(granteeAddress.String(), granterAddress.String(), height)
		err = m.db.DeleteFeeGrantAllowance(allowanceToRemove)
		if err != nil {
			return fmt.Errorf("error while deleting fee grant allowance: %s", err)

		}
	}
	return nil

}
