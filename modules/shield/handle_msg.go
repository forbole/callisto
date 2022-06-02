package shield

import (
	"fmt"
	"strconv"

	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/bdjuno/v3/types"
	juno "github.com/forbole/juno/v3/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *shieldtypes.MsgCreatePool:
		return m.HandleMsgCreatePool(index, tx, cosmosMsg)

	case *shieldtypes.MsgUpdatePool:
		return m.HandleMsgUpdatePool(tx, cosmosMsg)

	case *shieldtypes.MsgPausePool:
		return m.HandleMsgPausePool(tx, cosmosMsg)

	case *shieldtypes.MsgResumePool:
		return m.HandleMsgResumePool(tx, cosmosMsg)

		// case *shieldtypes.MsgWithdrawRewards:
		// 	return m.HandleMsgWithdrawRewards(tx, cosmosMsg)

		// case *shieldtypes.MsgDepositCollateral:
		// 	return m.HandleMsgDepositCollateral(tx, cosmosMsg)

		// case *shieldtypes.MsgWithdrawCollateral:
		// 	return m.HandleMsgWithdrawCollateral(tx, cosmosMsg)

		// case *shieldtypes.MsgPurchaseShield:
		// 	return m.HandleMsgPurchaseShield(tx, cosmosMsg)

		// case *shieldtypes.MsgUpdateSponsor:
		// 	return m.HandleMsgUpdateSponsor(tx, cosmosMsg)

		// case *shieldtypes.MsgStakeForShield:
		// 	return m.HandleMsgStakeForShield(tx, cosmosMsg)

		// case *shieldtypes.MsgUnstakeFromShield:
		// 	return m.HandleMsgUnstakeFromShield(tx, cosmosMsg)

		// case *shieldtypes.MsgWithdrawReimbursement:
		// 	return m.HandleMsgWithdrawReimbursement(tx, cosmosMsg)

	}

	return nil
}

// HandleMsgCreatePool allows to properly handle a MsgCreatePool
func (m *Module) HandleMsgCreatePool(index int, tx *juno.Tx, msg *shieldtypes.MsgCreatePool) error {

	// Get create pool event
	createPoolEvent, err := tx.FindEventByType(index, "create_pool")
	if err != nil {
		return fmt.Errorf("error while getting create pool event: %s", err)
	}

	// Get pool ID
	poolIDStr, err := tx.FindAttributeByKey(createPoolEvent, shieldtypes.AttributeKeyPoolID)
	if err != nil {
		return fmt.Errorf("error while getting pool ID from MsgCreatePool: %s", err)
	}

	// Convert pool ID
	poolID, err := strconv.ParseUint(poolIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("error while parsing pool ID to uint64: %s", err)
	}

	pool := types.NewShieldPool(
		poolID, msg.From, msg.Shield, msg.Deposit, msg.Sponsor,
		msg.SponsorAddr, msg.Description, msg.ShieldLimit, false, tx.Height,
	)

	return m.db.SaveShieldPool(pool)
}

// HandleMsgUpdatePool allows to properly handle a MsgUpdatePool
func (m *Module) HandleMsgUpdatePool(tx *juno.Tx, msg *shieldtypes.MsgUpdatePool) error {

	// Sponsor, sponsor address, and pool pause status will not be updated with ON CONFLICT statement
	pool := types.NewShieldPool(
		msg.PoolId, msg.From, msg.Shield, msg.ServiceFees, "",
		"", msg.Description, msg.ShieldLimit, false, tx.Height,
	)

	return m.db.SaveShieldPool(pool)
}

// HandleMsgPausePool allows to properly handle a MsgPausePool
func (m *Module) HandleMsgPausePool(tx *juno.Tx, msg *shieldtypes.MsgPausePool) error {
	pause := true
	return m.db.UpdatePoolPauseStatus(msg.PoolId, pause)
}

// HandleMsgResumePool allows to properly handle a MsgResumePool
func (m *Module) HandleMsgResumePool(tx *juno.Tx, msg *shieldtypes.MsgResumePool) error {
	pause := false
	return m.db.UpdatePoolPauseStatus(msg.PoolId, pause)
}

// // HandleMsgWithdrawRewards allows to properly handle a MsgWithdrawRewards
// func (m *Module) HandleMsgWithdrawRewards(tx *juno.Tx, msg *shieldtypes.MsgWithdrawRewards) error {

// 	return nil
// }

// // HandleMsgDepositCollateral allows to properly handle a MsgDepositCollateral
// func (m *Module) HandleMsgDepositCollateral(tx *juno.Tx, msg *shieldtypes.MsgDepositCollateral) error {

// 	return nil
// }

// // HandleMsgWithdrawCollateral allows to properly handle a MsgWithdrawCollateral
// func (m *Module) HandleMsgWithdrawCollateral(tx *juno.Tx, msg *shieldtypes.MsgWithdrawCollateral) error {

// 	return nil
// }

// // HandleMsgPurchaseShield allows to properly handle a MsgPurchaseShield
// func (m *Module) HandleMsgPurchaseShield(tx *juno.Tx, msg *shieldtypes.MsgPurchaseShield) error {

// 	return nil
// }

// // HandleMsgUpdateSponsor allows to properly handle a MsgUpdateSponsor
// func (m *Module) HandleMsgUpdateSponsor(tx *juno.Tx, msg *shieldtypes.MsgUpdateSponsor) error {

// 	return nil
// }

// // HandleMsgStakeForShield allows to properly handle a MsgStakeForShield
// func (m *Module) HandleMsgStakeForShield(tx *juno.Tx, msg *shieldtypes.MsgStakeForShield) error {

// 	return nil
// }

// // HandleMsgUnstakeFromShield allows to properly handle a MsgUnstakeFromShield
// func (m *Module) HandleMsgUnstakeFromShield(tx *juno.Tx, msg *shieldtypes.MsgUnstakeFromShield) error {

// 	return nil
// }

// // HandleMsgWithdrawReimbursement allows to properly handle a MsgWithdrawReimbursement
// func (m *Module) HandleMsgWithdrawReimbursement(tx *juno.Tx, msg *shieldtypes.MsgWithdrawReimbursement) error {

// 	return nil
// }
