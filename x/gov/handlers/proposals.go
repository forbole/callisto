package handlers

import (
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/desmos-labs/juno/parse/client"
	juno "github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
)

// HandleMsgSubmitProposal allows to properly handle a HandleMsgSubmitProposal
func HandleMsgSubmitProposal(tx juno.Tx, msg gov.MsgSubmitProposal, db database.BigDipperDb, cp client.ClientProxy) error {
	//get proposal messages
	
	return nil
}

// HandleMsgDeposit allows to properly handle a HandleMsgDeposit
func HandleMsgDeposit(tx juno.Tx, msg gov.MsgDeposit, db database.BigDipperDb, cp client.ClientProxy) error {
	return nil
}

// HandleMsgVote allows to properly handle a HandleMsgVote
func HandleMsgVote(tx juno.Tx, msg gov.MsgVote, db database.BigDipperDb, cp client.ClientProxy) error {
	return nil
}
