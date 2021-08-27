package types

import (
	"time"
)

// RecordRow represents a single row inside the record table
type RecordRow struct {
	OneRowID     		bool		`db:"one_row_id"`
	Ipld 				string 		`db:"ipld"`
	Context 			string 		`db:"context"`
	RecordID 			string 		`db:"record_id"`
	RecordRoute			string 		`db:"record_route"`
	RecordType 			string 		`db:"record_type"`
	ContentFingerprints string 		`db:"content_fingerprints"`
	ContentMetadata 	string 		`db:"content_metadata"`
	RecordNotes 		string 		`db:"record_notes"`
	RecordTimestamp 	time.Time 	`db:"record_timestamp"`
	RecordVersion 		uint64 		`db:"record_version"`
	Stakeholders 		string 		`db:"stakeholders"`
	Height 				int64 		`db:"height"`
}

// NewRecordRow builds a new RecordRow instance
func NewRecordRow(
	ipld string, 
	context string,
	recordID string,
	recordRoute string,
	recordType string,
	contentFingerprints string,
	contentMetadata string,
	recordNotes string,
	recordTimestamp time.Time,
	recordVersion uint64,
	stakeholders string,
	height int64,
	) RecordRow {
	return RecordRow{
		OneRowID: true,
		Ipld: ipld,
		Context: context,
		RecordID: recordID,
		RecordRoute: recordRoute,
		RecordType: recordType,
		ContentFingerprints: contentFingerprints,
		ContentMetadata: contentMetadata,
		RecordNotes: recordNotes,
		RecordTimestamp: recordTimestamp,
		RecordVersion: recordVersion,
		Stakeholders: stakeholders,
		Height: height,
	}
}

// Equals return true if two RecordRow are the same
func (r RecordRow) Equals(another RecordRow) bool {
	return r.Ipld == another.Ipld &&
		r.Context == another.Context &&
		r.RecordID == another.RecordID &&
		r.RecordRoute == another.RecordRoute &&
		r.RecordType == another.RecordType &&
		r.ContentFingerprints == another.ContentFingerprints &&
		r.ContentMetadata == another.ContentMetadata &&
		r.RecordNotes == another.RecordNotes &&
		r.RecordTimestamp.Equal(another.RecordTimestamp) &&
		r.RecordVersion == another.RecordVersion &&
		r.Stakeholders == another.Stakeholders &&
		r.Height == another.Height 
}

// --------------------------------------------------------------------------------------------------------------------

// IscnParamsRow represents a single row inside the iscn_params table
type IscnParamsRow struct {
	OneRowID bool   `db:"one_row_id"`
	Params   string `db:"params"`
	Height   int64  `db:"height"`
}

// NewIscnParamsRow builds a new IscnParamsRow instance
func NewIscnParamsRow(
	params string, height int64,
) IscnParamsRow {
	return IscnParamsRow{
		OneRowID: true,
		Params:   params,
		Height:   height,
	}
}

// Equal reports whether m and n represent the same table rows.
func (m IscnParamsRow) Equal(n IscnParamsRow) bool {
	return m.Params == n.Params &&
		m.Height == n.Height
}
