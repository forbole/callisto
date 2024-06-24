package auth

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v6/types"
	"github.com/rs/zerolog/log"

	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	moduleutils "github.com/forbole/callisto/v4/modules/utils"
	"github.com/forbole/callisto/v4/types"
	"github.com/forbole/callisto/v4/utils"
)

var msgFilter = map[string]bool{
	"/cosmos.vesting.v1beta1.MsgCreateVestingAccount": true,
}

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ int, executedMsg juno.Message, tx *juno.Transaction) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg juno.Message, tx *juno.Transaction) error {
	addresses, err := m.messagesParser(tx)
	if err != nil {
		log.Error().Str("module", "auth").Err(err).
			Str("operation", "refresh account").
			Msgf("error while refreshing accounts after message of type %s", msg.GetType())
	}

	if _, ok := msgFilter[msg.GetType()]; !ok {
		return nil
	}

	log.Debug().Str("module", "auth").Str("hash", tx.TxHash).Uint64("height", tx.Height).Msg(fmt.Sprintf("handling auth message %s", msg.GetType()))

	if msg.GetType() == "/cosmos.vesting.v1beta1.MsgCreateVestingAccount" {
		// Store tx timestamp as start_time of the created vesting account
		timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
		if err != nil {
			return fmt.Errorf("error while parsing time: %s", err)
		}

		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &vestingtypes.MsgCreateVestingAccount{})
		err = m.handleMsgCreateVestingAccount(cosmosMsg, timestamp)
		if err != nil {
			return fmt.Errorf("error while handling MsgCreateVestingAccount %s", err)
		}
	}

	return m.RefreshAccounts(int64(tx.Height), moduleutils.FilterNonAccountAddresses(addresses))
}

func (m *Module) handleMsgCreateVestingAccount(msg *vestingtypes.MsgCreateVestingAccount, txTimestamp time.Time) error {

	accAddress, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return fmt.Errorf("error while converting account address %s", err)
	}

	// store account in database
	err = m.db.SaveAccounts([]types.Account{types.NewAccount(accAddress.String())})
	if err != nil {
		return fmt.Errorf("error while storing vesting account: %s", err)
	}

	bva, err := vestingtypes.NewBaseVestingAccount(
		authttypes.NewBaseAccountWithAddress(accAddress), msg.Amount, msg.EndTime,
	)
	if err != nil {
		return fmt.Errorf("error while new base vesting account: %s", err)
	}

	err = m.db.StoreBaseVestingAccountFromMsg(bva, txTimestamp)
	if err != nil {
		return fmt.Errorf("error while storing base vesting account from msg %s", err)
	}
	return nil
}
