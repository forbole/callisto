package handlers

import (
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/desmos-labs/juno/parse/client"
	juno "github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
)

// HandleMsgSubmitProposal allows to properly handle a HandleMsgSubmitProposal
func HandleMsgSubmitProposal(tx juno.Tx, msg gov.MsgSubmitProposal, db database.BigDipperDb, cp client.ClientProxy) error {
	return nil
}

func HandleMsgDeposit(tx juno.Tx, msg gov.MsgDeposit, db database.BigDipperDb, cp client.ClientProxy) error {
	return nil
}

func HandleMsgVote(tx juno.Tx, msg gov.MsgVote, db database.BigDipperDb, cp client.ClientProxy) error {
	return nil
}
