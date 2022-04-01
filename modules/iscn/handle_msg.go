package iscn

import (
	iscntypes "github.com/likecoin/likechain/x/iscn/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v3/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *iscntypes.MsgCreateIscnRecord:
		return m.handleMsgCreateIscnRecord(tx, index, cosmosMsg)

	case *iscntypes.MsgUpdateIscnRecord:
		return m.handleMsgUpdateIscnRecord(tx, cosmosMsg)

	case *iscntypes.MsgChangeIscnRecordOwnership:
		return m.handleMsgChangeIscnRecordOwnership(cosmosMsg)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// handleMsgCreateIscnRecord handles storing iscn records inside the database
func (m *Module) handleMsgCreateIscnRecord(tx *juno.Tx, index int, msg *iscntypes.MsgCreateIscnRecord) error {
	return m.storeIscnRecordFromMessage(tx.Height, tx, index, msg)
}

// handleMsgUpdateIscnRecord handles updating the iscn data inside the database
func (m *Module) handleMsgUpdateIscnRecord(tx *juno.Tx, msg *iscntypes.MsgUpdateIscnRecord) error {
	return m.updateIscnRecordFromMessage(tx.Height, msg)
}

// handleMsgChangeIscnRecordOwnership handles updating the iscn record ownership inside the database
func (m *Module) handleMsgChangeIscnRecordOwnership(msg *iscntypes.MsgChangeIscnRecordOwnership) error {
	return m.updateIscnRecordOwnershipFromMessage(msg)
}
