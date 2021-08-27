package utils

import (
	"context"

	"github.com/desmos-labs/juno/client"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"

	"github.com/cosmos/cosmos-sdk/codec"

	iscntypes "github.com/likecoin/likechain/x/iscn/types"
	juno "github.com/desmos-labs/juno/types"
)


// StoreIscnRecordFromMessage handles storing the new iscn record inside the database
// and returns new iscn record instance
func StoreIscnRecordFromMessage(
	height int64, tx *juno.Tx, index int, msg *iscntypes.MsgCreateIscnRecord, iscnClient iscntypes.QueryClient, cdc codec.Marshaler, db *database.Db,
) ( *types.IscnRecord, error) {
	header := client.GetHeightRequestHeader(height)

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

	record := res.Records


	// Store the record info
	recordObject := types.NewIscnRecord(
		record.Ipld,
		record.Context,
		record.RecordID,
		record.RecordRoute,
		record.RecordType,
		record.ContentFingerprints,
		record.ContentMetadata,
		record.RecordNotes,
		record.RecordTimestamp,
		record.RecordVersion,
		record.Stakeholders,
		record.Height,
	)

	iscnRecord := types.NewIscnRecord(recordObject)
	return db.SaveNewIscnRecord([]types.IscnRecord{iscnRecord})



}
