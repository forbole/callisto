package staking

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/desmos-labs/juno/parse/worker"
	"github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"

	"fmt"
)

func MsgHandler(tx types.Tx, index int, msg sdk.Msg, w worker.Worker) error {
	log.Info().Str("tx_hash", tx.TxHash).Int("msg_index", index).Str("msg_type", msg.Type()).Msg("found message")

	if len(tx.Logs) == 0 {
		log.Info().Msg(fmt.Sprintf("Skipping message at index %d of tx hash %s as it was not successull",
			index, tx.TxHash))
		return nil
	}
	db, ok := w.Db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("invalid database")
	}
	switch stakingMsg := msg.(type) {
	case staking.MsgEditValidator:
		// TODO: Handle message here
		//store commission rate
		StoreModifiedCommission(stakingMsg, w.ClientProxy, db)
	case staking.MsgDelegate:
		StoreDelegation(stakingMsg, w.ClientProxy, timestamp, tx.Height, db)
	}

	return nil
}

//once there is delegation, store current volting power
//if delegators delegation address is the self delegation address
//store the /staking/delegators/{delegatorAddr}/delegations/{validatorAddr}
//else store in other table and see what happen
func StoreDelegation(msg staking.MsgDelegate, cp client.ClientProxy, time time.Time, height int64, db database.BigDipperDb) error {
	validatorAddress := msg.ValidatorAddress
	deligatorAddress := msg.DelegatorAddress
	if found, _ := db.HasValidator(validatorAddress.String()); !found {
		return nil
	}
	if found, _ := db.HasValidator(deligatorAddress.String()); !found {
		return nil
	}
	var delegation staking.Delegation
	endpoint := fmt.Sprintf("/staking/delegators/%s/delegations/%s", deligatorAddress.String(), validatorAddress.String())
	height, ok := cp.QueryLCDWithHeight(endpoint, &delegation)
	if ok != nil {
		return nil
	}
	//check if the delegation is self delegation
	selfAddress := sdk.AccAddress(validatorAddress.Bytes())
	if deligatorAddress.Equals(selfAddress) {
		db.SaveSelfDelegation(delegation, time, height)
	} else {
		//If that is the message that delegated to other account
	}
	return nil
}

func StoreModifiedCommission(msg staking.MsgEditValidator, cp client.ClientProxy, db database.BigDipperDb) error {
	//should I take from REST or store the message?
	//store the message
	address := msg.ValidatorAddress
	if found, _ := db.HasValidator(address.String()); !found {
		return nil
	}

	var validator staking.Validator
	endpoint := fmt.Sprintf("/staking/validators/$s", address.String())
	height, ok := cp.QueryLCDWithHeight(endpoint, &validator)
	if ok != nil {
		return ok
	}

	db.SaveVaildatorComission(validator, height)
	return nil
}
