package staking

import (
	"context"
	"time"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	juno "github.com/desmos-labs/juno/types"

	"github.com/forbole/bdjuno/modules/common/utils"
	"github.com/forbole/bdjuno/types"
)

// StoreDelegationFromMessage handles a MsgDelegate and saves the delegation inside the database
func StoreDelegationFromMessage(
	height int64, msg *stakingtypes.MsgDelegate, stakingClient stakingtypes.QueryClient, db DB,
) error {
	header := utils.GetHeightRequestHeader(height)
	res, err := stakingClient.Delegation(
		context.Background(),
		&stakingtypes.QueryDelegationRequest{
			DelegatorAddr: msg.DelegatorAddress,
			ValidatorAddr: msg.ValidatorAddress,
		},
		header,
	)
	if err != nil {
		return err
	}

	delegation := ConvertDelegationResponse(height, *res.DelegationResponse)
	return db.SaveDelegations([]types.Delegation{delegation})
}

// StoreRedelegationFromMessage handles a MsgBeginRedelegate by saving the redelegation inside the database,
// and returns the new redelegation instance
func StoreRedelegationFromMessage(
	tx *juno.Tx, index int, msg *stakingtypes.MsgBeginRedelegate, db DB,
) (*types.Redelegation, error) {
	event, err := tx.FindEventByType(index, stakingtypes.EventTypeRedelegate)
	if err != nil {
		return nil, err
	}

	completionTimeStr, err := tx.FindAttributeByKey(event, stakingtypes.AttributeKeyCompletionTime)
	if err != nil {
		return nil, err
	}

	completionTime, err := time.Parse(time.RFC3339, completionTimeStr)
	if err != nil {
		return nil, err
	}

	redelegation := types.NewRedelegation(
		msg.DelegatorAddress,
		msg.ValidatorSrcAddress,
		msg.ValidatorDstAddress,
		msg.Amount,
		completionTime,
		tx.Height,
	)

	return &redelegation, db.SaveRedelegations([]types.Redelegation{redelegation})
}

// StoreUnbondingDelegationFromMessage handles a MsgUndelegate storing the new unbonding delegation inside the database,
// and returns the new unbonding delegation instance
func StoreUnbondingDelegationFromMessage(
	tx *juno.Tx, index int, msg *stakingtypes.MsgUndelegate, db DB,
) (*types.UnbondingDelegation, error) {
	event, err := tx.FindEventByType(index, stakingtypes.EventTypeUnbond)
	if err != nil {
		return nil, err
	}

	completionTimeStr, err := tx.FindAttributeByKey(event, stakingtypes.AttributeKeyCompletionTime)
	if err != nil {
		return nil, err
	}

	completionTime, err := time.Parse(time.RFC3339, completionTimeStr)
	if err != nil {
		return nil, err
	}

	delegation := types.NewUnbondingDelegation(
		msg.DelegatorAddress,
		msg.ValidatorAddress,
		msg.Amount,
		completionTime,
		tx.Height,
	)

	return &delegation, db.SaveUnbondingDelegations([]types.UnbondingDelegation{delegation})
}
