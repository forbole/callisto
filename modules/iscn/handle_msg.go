package iscn

import (
	"github.com/forbole/bdjuno/database"
	iscnutils "github.com/forbole/bdjuno/modules/iscn/utils"

	iscntypes "github.com/likecoin/likechain/x/iscn/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/desmos-labs/juno/types"
	
)

// HandleMsg allows to handle the different utils related to the iscn record module
func HandleMsg(
	tx *juno.Tx, index int, msg sdk.Msg, iscnClient iscntypes.QueryClient, db *database.Db,
) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *iscntypes.MsgCreateIscnRecord:
		return handleMsgCreateIscnRecord(tx, index, cosmosMsg, iscnClient, db)

	case *iscntypes.MsgUpdateIscnRecord:
		return handleMsgUpdateIscnRecord(tx, index, cosmosMsg, iscnClient, db)

	case *iscntypes.MsgChangeIscnRecordOwnership:
		return handleMsgChangeIscnRecordOwnership(tx, index, cosmosMsg, iscnClient, db)
	}


	return nil
}

// ---------------------------------------------------------------------------------------------------------------------


// handleMsgCreateIscnRecord handles storing iscn records inside the database
func handleMsgCreateIscnRecord(
	tx *juno.Tx, index int, msg *iscntypes.MsgCreateIscnRecord,
	iscnClient iscntypes.QueryClient, db *database.Db,
) error {
	err := iscnutils.StoreIscnRecordFromMessage(tx.Height, tx, index, msg, iscnClient, db) 
	if err != nil {
		return err
	}

	return err 
}

// handleMsgUpdateIscnRecord handles updating the iscn data inside the database
func handleMsgUpdateIscnRecord(
	tx *juno.Tx, index int, msg *iscntypes.MsgUpdateIscnRecord,
	iscnClient iscntypes.QueryClient, db *database.Db,
) error {
	// to do 
	return nil 
}


// handleMsgChangeIscnRecordOwnership handles updating the iscn record ownership inside the database
func handleMsgChangeIscnRecordOwnership(
	tx *juno.Tx, index int, msg *iscntypes.MsgChangeIscnRecordOwnership,
	iscnClient iscntypes.QueryClient, db *database.Db,
) error {
	// err := iscnutils.UpdateIscnRecordOwnershipFromMessage(tx.Height, tx, index, msg, iscnClient, db) 
	// if err != nil {
	// 	return err
	// }

	return nil 
}
