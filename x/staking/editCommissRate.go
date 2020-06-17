package staking

import (
	
	"github.com/desmos-labs/juno/parse/worker"
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"
	"github.com/desmos-labs/juno/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
			StoreModifiedVaildator(stakingMsg,db) 
	}

	return nil
}


func StoreModifiedVaildator(msg staking.MsgEditValidator,db database.BigDipperDb){
	//should I take from REST or store the message?
	//store the message
	vc := database.ValidatorCommission{
		ValidatorAddress : msg.ValidatorAddr,
		Commission       : msg.CommissionRate,
	}
	db.SaveVaildatorComission(vc)
}