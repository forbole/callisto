package iscn

import (
	juno "github.com/forbole/juno/v3/types"
	iscntypes "github.com/likecoin/likechain/x/iscn/types"

	"github.com/forbole/bdjuno/v2/types"
)

// storeIscnRecordFromMessage handles storing new iscn record inside the database
func (m *Module) storeIscnRecordFromMessage(height int64, tx *juno.Tx, index int, msg *iscntypes.MsgCreateIscnRecord) error {
	event, err := tx.FindEventByType(index, iscntypes.EventTypeIscnRecord)
	if err != nil {
		return err
	}

	id, err := tx.FindAttributeByKey(event, iscntypes.AttributeKeyIscnId)
	if err != nil {
		return err
	}

	// Get the record
	res, err := m.source.GetRecordsByID(height, id)
	if err != nil {
		return err
	}

	return m.db.SaveIscnRecord(types.NewIscnRecord(
		res.Owner,
		id,
		res.LatestVersion,
		res.Records[0].Ipld,
		types.NewRecord(
			id,
			msg.Record.RecordNotes,
			msg.Record.ContentFingerprints,
			msg.Record.Stakeholders,
			msg.Record.ContentMetadata,
		),
		height,
	))
}

// updateIscnRecordFromMessage handles updating the existing iscn record inside the database
func (m *Module) updateIscnRecordFromMessage(height int64, msg *iscntypes.MsgUpdateIscnRecord) error {

	id := msg.IscnId

	// Get the record
	res, err := m.source.GetRecordsByID(height, msg.IscnId)
	if err != nil {
		return err
	}

	return m.db.UpdateIscnRecord(types.NewIscnRecord(
		res.Owner,
		id,
		res.LatestVersion,
		res.Records[0].Ipld,
		types.NewRecord(
			id,
			msg.Record.RecordNotes,
			msg.Record.ContentFingerprints,
			msg.Record.Stakeholders,
			msg.Record.ContentMetadata,
		),
		height,
	))
}

// updateIscnRecordOwnershipFromMessage handles updating ownership of the existing iscn record inside the database
func (m *Module) updateIscnRecordOwnershipFromMessage(msg *iscntypes.MsgChangeIscnRecordOwnership) error {
	return m.db.UpdateIscnRecordOwnership(types.NewIscnChangeOwnership(msg.From, msg.IscnId, msg.NewOwner))
}
