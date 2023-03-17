package shield

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/bdjuno/v3/types"
	juno "github.com/forbole/juno/v4/types"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
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

	case *shieldtypes.MsgWithdrawRewards:
		return m.HandleMsgWithdrawRewards(tx, cosmosMsg)

	case *shieldtypes.MsgWithdrawForeignRewards:
		return m.HandleMsgWithdrawForeignRewards(tx, cosmosMsg)

	case *shieldtypes.MsgDepositCollateral:
		return m.HandleMsgDepositCollateral(tx, cosmosMsg)

	case *shieldtypes.MsgWithdrawCollateral:
		return m.HandleMsgWithdrawCollateral(tx, cosmosMsg)

	case *shieldtypes.MsgPurchaseShield:
		return m.HandleMsgPurchaseShield(tx, cosmosMsg)

	case *shieldtypes.MsgUpdateSponsor:
		return m.HandleMsgUpdateSponsor(tx, cosmosMsg)

	case *shieldtypes.MsgStakeForShield:
		return m.HandleMsgStakeForShield(tx, cosmosMsg)

	case *shieldtypes.MsgUnstakeFromShield:
		return m.HandleMsgUnstakeFromShield(tx, cosmosMsg)

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
		poolID, msg.Shield[0].Amount, msg.Deposit, msg.Sponsor,
		msg.SponsorAddr, msg.Description, msg.ShieldLimit, false, tx.Height,
	)

	return m.db.SaveShieldPool(pool)
}

// HandleMsgUpdatePool allows to properly handle a MsgUpdatePool
func (m *Module) HandleMsgUpdatePool(tx *juno.Tx, msg *shieldtypes.MsgUpdatePool) error {

	// Sponsor, sponsor address, and pool pause status will not be updated with ON CONFLICT statement
	pool := types.NewShieldPool(
		msg.PoolId, msg.Shield[0].Amount, msg.ServiceFees, "",
		"", msg.Description, msg.ShieldLimit, false, tx.Height,
	)

	return m.db.SaveShieldPool(pool)
}

// HandleMsgPausePool allows to properly handle a MsgPausePool
func (m *Module) HandleMsgPausePool(tx *juno.Tx, msg *shieldtypes.MsgPausePool) error {
	pause := true
	return m.db.UpdatePoolPauseStatus(msg.PoolId, pause, tx.Height)
}

// HandleMsgResumePool allows to properly handle a MsgResumePool
func (m *Module) HandleMsgResumePool(tx *juno.Tx, msg *shieldtypes.MsgResumePool) error {
	pause := false
	return m.db.UpdatePoolPauseStatus(msg.PoolId, pause, tx.Height)
}

// HandleMsgWithdrawRewards allows to properly handle a MsgWithdrawRewards
func (m *Module) HandleMsgWithdrawRewards(tx *juno.Tx, msg *shieldtypes.MsgWithdrawRewards) error {
	return m.db.WithdrawNativeRewards(msg.From, tx.Height)
}

// HandleMsgWithdrawForeignRewards allows to properly handle a MsgWithdrawForeignRewards
func (m *Module) HandleMsgWithdrawForeignRewards(tx *juno.Tx, msg *shieldtypes.MsgWithdrawForeignRewards) error {
	return m.db.WithdrawForeignRewards(msg.From, tx.Height)
}

// HandleMsgDepositCollateral allows to properly handle a MsgDepositCollateral
func (m *Module) HandleMsgDepositCollateral(tx *juno.Tx, msg *shieldtypes.MsgDepositCollateral) error {
	collateral, err := m.db.GetShieldProviderCollateral(msg.From)
	if err != nil {
		return fmt.Errorf("error while getting shield provider collateral: %s", err)
	}
	updatedCollateral := collateral + msg.Collateral[0].Amount.Int64()

	return m.db.UpdateShieldProviderCollateral(msg.From, updatedCollateral, tx.Height)
}

// HandleMsgWithdrawCollateral allows to properly handle a MsgWithdrawCollateral
func (m *Module) HandleMsgWithdrawCollateral(tx *juno.Tx, msg *shieldtypes.MsgWithdrawCollateral) error {
	collateral, err := m.db.GetShieldProviderCollateral(msg.From)
	if err != nil {
		return fmt.Errorf("error while getting shield provider collateral: %s", err)
	}

	if msg.Collateral[0].Amount.Int64() >= collateral {
		return m.db.UpdateShieldProviderCollateral(msg.From, 0, tx.Height)
	}

	updatedCollateral := collateral - msg.Collateral[0].Amount.Int64()
	return m.db.UpdateShieldProviderCollateral(msg.From, updatedCollateral, tx.Height)
}

// HandleMsgPurchaseShield allows to properly handle a MsgPurchaseShield
func (m *Module) HandleMsgPurchaseShield(tx *juno.Tx, msg *shieldtypes.MsgPurchaseShield) error {
	shield := types.NewShieldPurchase(
		msg.PoolId, msg.From, msg.Shield[0].Amount, msg.Description, tx.Height,
	)

	return m.db.SaveShieldPurchase(shield)
}

// HandleMsgUpdateSponsor allows to properly handle a MsgUpdateSponsor
func (m *Module) HandleMsgUpdateSponsor(tx *juno.Tx, msg *shieldtypes.MsgUpdateSponsor) error {

	return m.db.UpdatePoolSponsor(msg.PoolId, msg.Sponsor, msg.SponsorAddr, tx.Height)
}

// HandleMsgStakeForShield allows to properly handle a MsgStakeForShield
func (m *Module) HandleMsgStakeForShield(tx *juno.Tx, msg *shieldtypes.MsgStakeForShield) error {
	delegation, err := m.db.GetShieldProviderDelegation(msg.From)
	if err != nil {
		return fmt.Errorf("error while getting shield provider delegation: %s", err)
	}
	totalDelegation := delegation + msg.Shield[0].Amount.Int64()
	return m.db.UpdateShieldProviderDelegation(msg.From, totalDelegation, tx.Height)
}

// HandleMsgUnstakeFromShield allows to properly handle a MsgUnstakeFromShield
func (m *Module) HandleMsgUnstakeFromShield(tx *juno.Tx, msg *shieldtypes.MsgUnstakeFromShield) error {
	delegation, err := m.db.GetShieldProviderDelegation(msg.From)
	if err != nil {
		return fmt.Errorf("error while getting shield provider delegation: %s", err)
	}

	if msg.Shield[0].Amount.Int64() >= delegation {
		return m.db.UpdateShieldProviderDelegation(msg.From, 0, tx.Height)
	}

	updatedDelegation := delegation - msg.Shield[0].Amount.Int64()
	return m.db.UpdateShieldProviderDelegation(msg.From, updatedDelegation, tx.Height)
}
