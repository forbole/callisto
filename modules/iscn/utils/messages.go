package utils

import (
	"context"

	"github.com/desmos-labs/juno/client"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"

	iscntypes "github.com/likecoin/likechain/x/iscn/types"
	juno "github.com/desmos-labs/juno/types"
)


// StoreIscnRecordFromMessage handles storing the new iscn record inside the database
// and returns new iscn record instance
func StoreIscnRecordFromMessage(
	height int64, tx *juno.Tx, index int, msg *iscntypes.MsgCreateIscnRecord, iscnClient iscntypes.QueryClient, db *database.Db,
) ( error) {

	event, err := tx.FindEventByType(index, iscntypes.EventTypeIscnRecord)
	if err != nil {
		return err
	}

	id, err := tx.FindAttributeByKey(event, iscntypes.AttributeKeyIscnId)
	if err != nil {
		return err
	}

	// Get the record
	res, err := iscnClient.RecordsById(
		context.Background(),
		&iscntypes.QueryRecordsByIdRequest{IscnId: id},
		client.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}
	iscnRecord := iscntypes.IscnRecord{RecordNotes: msg.Record.RecordNotes, ContentFingerprints: msg.Record.ContentFingerprints, Stakeholders: msg.Record.Stakeholders, ContentMetadata: msg.Record.ContentMetadata}
	newIscnRecord := types.NewIscnRecord(res.Owner, id, res.LatestVersion, res.Records[0].Ipld, iscnRecord, height)
	return db.SaveIscnRecord(newIscnRecord)
}


// UpdateIscnRecordOwnershipFromMessage handles updating new iscn record owner inside the database
func UpdateIscnRecordOwnershipFromMessage(
	height int64, tx *juno.Tx, index int, msg *iscntypes.MsgChangeIscnRecordOwnership, iscnClient iscntypes.QueryClient, db *database.Db,
) ( error) {

	updatedIscnRecord := types.NewIscnChangeOwnership(msg.From, msg.IscnId, msg.NewOwner)
	return db.UpdateIscnRecordOwnership(updatedIscnRecord)
}