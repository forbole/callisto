package staking

import (
	"fmt"
	"time"

	"github.com/forbole/bdjuno/v2/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	juno "github.com/forbole/juno/v2/types"
)

// storeDelegationFromMessage handles a MsgDelegate and saves the delegation inside the database
func (m *Module) storeDelegationFromMessage(height int64, msg *stakingtypes.MsgDelegate) error {
	delegation, err := m.source.GetDelegation(height, msg.DelegatorAddress, msg.ValidatorAddress)
	if err != nil {
		return err
	}

	return m.db.SaveDelegations([]types.Delegation{
		convertDelegationResponse(height, delegation),
	})
}

// storeRedelegationFromMessage handles a MsgBeginRedelegate by saving the redelegation inside the database,
// and returns the new redelegation instance
func (m *Module) storeRedelegationFromMessage(
	tx *juno.Tx, index int, msg *stakingtypes.MsgBeginRedelegate,
) (*types.Redelegation, error) {
	event, err := tx.FindEventByType(index, stakingtypes.EventTypeRedelegate)
	if err != nil {
		return nil, fmt.Errorf("error while searching for EventTypeRedelegate: %s", err)
	}

	completionTimeStr, err := tx.FindAttributeByKey(event, stakingtypes.AttributeKeyCompletionTime)
	if err != nil {
		return nil, fmt.Errorf("error while searching for AttributeKeyCompletionTime: %s", err)
	}

	completionTime, err := time.Parse(time.RFC3339, completionTimeStr)
	if err != nil {
		return nil, fmt.Errorf("error while covnerting redelegation completion time: %s", err)
	}

	redelegation := types.NewRedelegation(
		msg.DelegatorAddress,
		msg.ValidatorSrcAddress,
		msg.ValidatorDstAddress,
		msg.Amount,
		completionTime,
		tx.Height,
	)

	return &redelegation, m.db.SaveRedelegations([]types.Redelegation{redelegation})
}

// storeUnbondingDelegationFromMessage handles a MsgUndelegate storing the new unbonding delegation inside the database,
// and returns the new unbonding delegation instance
func (m *Module) storeUnbondingDelegationFromMessage(
	tx *juno.Tx, index int, msg *stakingtypes.MsgUndelegate,
) (*types.UnbondingDelegation, error) {
	event, err := tx.FindEventByType(index, stakingtypes.EventTypeUnbond)
	if err != nil {
		return nil, fmt.Errorf("error while searching for EventTypeUnbond: %s", err)
	}

	completionTimeStr, err := tx.FindAttributeByKey(event, stakingtypes.AttributeKeyCompletionTime)
	if err != nil {
		return nil, fmt.Errorf("error while searching for AttributeKeyCompletionTime: %s", err)
	}

	completionTime, err := time.Parse(time.RFC3339, completionTimeStr)
	if err != nil {
		return nil, fmt.Errorf("error while covnerting unbonding delegation completion time: %s", err)
	}

	delegation := types.NewUnbondingDelegation(
		msg.DelegatorAddress,
		msg.ValidatorAddress,
		msg.Amount,
		completionTime,
		tx.Height,
	)

	return &delegation, m.db.SaveUnbondingDelegations([]types.UnbondingDelegation{delegation})
}
