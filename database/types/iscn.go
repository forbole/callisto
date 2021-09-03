package types

import (
	iscntypes "github.com/likecoin/likechain/x/iscn/types"
)

type RecordRow struct {
	OneRowID     	bool      `db:"one_row_id"`
	Owner 			string 	  `db:"owner_address"`
	IscnId 			string	  `db:"iscn_id"`
	LatestVersion 	uint64 	  `db:"latest_version"`
	Ipld			string 	  `db:"ipld"`
	Data 			iscntypes.IscnRecord `db:"iscn_data"`
	Height       	int64     `db:"height"`
}


// NewRecordRow builds a new RecordRow instance
func NewRecordRow(owner string, iscnId string, latestVersion uint64, ipld string, data iscntypes.IscnRecord, height int64) RecordRow {
	return RecordRow{
		OneRowID: true,
		Owner: owner,
		IscnId: iscnId,
		LatestVersion: latestVersion,
		Ipld: ipld,
		Data: data,
		Height:   height,
	}
}

// Equal reports whether i and j represent the same table rows.
func (i RecordRow) Equal(j RecordRow) bool {
	return i.Owner == j.Owner && 
	i.LatestVersion == j.LatestVersion && 
	i.IscnId == j.IscnId &&
	i.Ipld == j.Ipld && 
	i.Height == j.Height
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