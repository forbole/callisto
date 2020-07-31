package handlers

import (
	"time"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/forbole/bdjuno/x/gov/types"
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
//refresh the proposal and record the deposit
func HandleMsgDeposit(tx juno.Tx, msg gov.MsgDeposit, db database.BigDipperDb, cp client.ClientProxy) error {
	return nil
}

// HandleMsgVote allows to properly handle a HandleMsgVote
func HandleMsgVote(tx juno.Tx, msg gov.MsgVote, db database.BigDipperDb, cp client.ClientProxy) error {
	//fetch from lcd & store voter in specific time
	var s gov.TallyResult
	height, err := cp.QueryLCDWithHeight(fmt.Sprintf("/gov/proposals/%s/tally",msg.ProposalID), &s)
	if err != nil {
		return err
	}
	//each vote voted
	
	return db.SaveTallyResult(types.NewTallyResult(msg.ProposalID,s.Yes.Int64(),s.Abstain.Int64(),s.No.Int64,s.NoWithVeto.Int64,
	tx.Height,time.Parse(time.RFC3339,tx.Timestamp.String())))
}
