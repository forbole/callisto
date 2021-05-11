package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	juno "github.com/desmos-labs/juno/types"

	forbolexdb "github.com/forbole/bdjuno/database/forbolex"
	"github.com/forbole/bdjuno/modules/common/staking"
)

// HandleMsg allows to handle the different utils related to the staking module
func HandleMsg(
	tx *juno.Tx, index int, msg sdk.Msg, stakingClient stakingtypes.QueryClient, cdc codec.Marshaler, db *forbolexdb.Db,
) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {

	case *stakingtypes.MsgCreateValidator:
		return staking.StoreValidatorFromMsgCreateValidator(tx.Height, cosmosMsg, cdc, db)

	case *stakingtypes.MsgDelegate:
		return staking.StoreDelegationFromMessage(tx.Height, cosmosMsg, stakingClient, db)

	case *stakingtypes.MsgBeginRedelegate:
		return handleMsgBeginRedelegate(tx, index, cosmosMsg, stakingClient, db)

	case *stakingtypes.MsgUndelegate:
		return handleMsgUndelegate(tx, index, cosmosMsg, stakingClient, db)
	}

	return nil
}

func handleMsgBeginRedelegate(
	tx *juno.Tx, index int, msg *stakingtypes.MsgBeginRedelegate,
	client stakingtypes.QueryClient, db *forbolexdb.Db,
) error {
	_, err := staking.StoreRedelegationFromMessage(tx, index, msg, db)
	if err != nil {
		return err
	}

	// Update the current delegations
	return staking.UpdateDelegations(tx.Height, msg.DelegatorAddress, client, db)
}

func handleMsgUndelegate(
	tx *juno.Tx, index int, msg *stakingtypes.MsgUndelegate,
	client stakingtypes.QueryClient, db *forbolexdb.Db,
) error {
	_, err := staking.StoreUnbondingDelegationFromMessage(tx, index, msg, db)
	if err != nil {
		return err
	}

	// Update the current delegations
	return staking.UpdateDelegations(tx.Height, msg.DelegatorAddress, client, db)
}
