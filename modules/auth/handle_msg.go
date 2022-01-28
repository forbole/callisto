package auth

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v2/types"
	"github.com/gogo/protobuf/proto"
	"github.com/rs/zerolog/log"

	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	"github.com/forbole/bdjuno/v2/modules/utils"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
	addresses, err := m.messagesParser(m.cdc, msg)
	if err != nil {
		log.Error().Str("module", "auth").Err(err).
			Str("operation", "refresh account").
			Msgf("error while refreshing accounts after message of type %s", proto.MessageName(msg))
	}

	switch cosmosMsg := msg.(type) {
	case *vestingtypes.MsgCreateVestingAccount:
		err = m.handleMsgCreateVestingAccount(cosmosMsg)
		if err != nil {
			return fmt.Errorf("error while handling to MsgCreateVestingAccount %s", err)
		}
	case *vestingtypes.MsgCreatePeriodicVestingAccount:
		err = m.handleMsgCreatePeriodicVestingAccount(cosmosMsg)
		if err != nil {
			return fmt.Errorf("error while handling to MsgCreatePeriodicVestingAccount %s", err)
		}
	}

	return m.RefreshAccounts(tx.Height, utils.FilterNonAccountAddresses(addresses))
}

func (m *Module) handleMsgCreateVestingAccount(msg *vestingtypes.MsgCreateVestingAccount) error {

	accAddress, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return fmt.Errorf("error while converting account address %s", err)
	}

	vestingAccount := vestingtypes.NewBaseVestingAccount(
		authttypes.NewBaseAccountWithAddress(accAddress), msg.Amount, msg.EndTime,
	)

	err = m.db.StoreBaseVestingAccountFromMsg(vestingAccount)
	if err != nil {
		return fmt.Errorf("error while storing to base vesting account from msg %s", err)
	}
	return nil
}

func (m *Module) handleMsgCreatePeriodicVestingAccount(msg *vestingtypes.MsgCreatePeriodicVestingAccount) error {

	accAddress, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return fmt.Errorf("error while converting account address %s", err)
	}

	vestingAccount := vestingtypes.NewPeriodicVestingAccount(
		authttypes.NewBaseAccountWithAddress(accAddress), nil, msg.StartTime, msg.VestingPeriods,
	)

	err = m.db.StorePeriodicVestingAccountFromMsg(vestingAccount)
	if err != nil {
		return fmt.Errorf("error while storing to periodic vesting account from msg %s", err)
	}

	return nil
}
