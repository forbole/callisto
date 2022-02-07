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

	if cosmosMsg, ok := msg.(*vestingtypes.MsgCreateVestingAccount); ok {
		err = m.handleMsgCreateVestingAccount(cosmosMsg)
		if err != nil {
			return fmt.Errorf("error while handling MsgCreateVestingAccount %s", err)
		}
	}

	return m.RefreshAccounts(tx.Height, utils.FilterNonAccountAddresses(addresses))
}

func (m *Module) handleMsgCreateVestingAccount(msg *vestingtypes.MsgCreateVestingAccount) error {
	va, err := convertBaseVestingAccountFromMsg(msg)
	if err != nil {
		return fmt.Errorf("error while converting MsgCreateVestingAccount to base vesting account %s", err)
	}

	err = m.db.StoreVestingAccountFromMsg(va)
	if err != nil {
		return fmt.Errorf("error while storing base vesting account from msg %s", err)
	}

	return nil
}

func convertBaseVestingAccountFromMsg(msg *vestingtypes.MsgCreateVestingAccount) (*vestingtypes.BaseVestingAccount, error) {

	accAddress, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return &vestingtypes.BaseVestingAccount{}, fmt.Errorf("error while converting account address %s", err)
	}
	account := authttypes.NewBaseAccountWithAddress(accAddress)

	return vestingtypes.NewBaseVestingAccount(account, msg.Amount, msg.EndTime), nil
}
