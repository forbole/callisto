package types

import (
	iscntypes "github.com/likecoin/likecoin-chain/v3/x/iscn/types"
)

// IscnParams represents the x/iscn parameters
type IscnParams struct {
	iscntypes.Params
	Height int64
}

// NewIscnParams allows to build a new IscnParams instance
func NewIscnParams(params iscntypes.Params, height int64) IscnParams {
	return IscnParams{
		Params: params,
		Height: height,
	}
}

type Record struct {
	IscnID              string
	RecordNotes         string
	ContentFingerprints []string
	Stakeholders        []iscntypes.IscnInput
	ContentMetadata     iscntypes.IscnInput
}

func NewRecord(iscnID string, recordNotes string, contentFingerprints []string,
	stakeholders []iscntypes.IscnInput, contentMetadata iscntypes.IscnInput) Record {
	return Record{
		IscnID:              iscnID,
		RecordNotes:         recordNotes,
		ContentFingerprints: contentFingerprints,
		Stakeholders:        stakeholders,
		ContentMetadata:     contentMetadata,
	}
}

// IscnRecord represents the x/iscn records
type IscnRecord struct {
	Owner         string
	IscnID        string
	LatestVersion uint64
	Ipld          string
	Data          Record
	Height        int64
}

// NewIscnRecord allows to build a new IscnRecord instance
func NewIscnRecord(owner string, iscnID string, latestVersion uint64, ipld string,
	data Record, height int64) IscnRecord {
	return IscnRecord{
		Owner:         owner,
		IscnID:        iscnID,
		LatestVersion: latestVersion,
		Ipld:          ipld,
		Data:          data,
		Height:        height,
	}
}

type IscnChangeOwnership struct {
	From     string
	IscnID   string
	NewOwner string
}

// NewIscnChangeOwnership allows to build a new IscnChangeOwnership instance
func NewIscnChangeOwnership(from string, iscnID string, newOwner string) IscnChangeOwnership {
	return IscnChangeOwnership{
		From:     from,
		IscnID:   iscnID,
		NewOwner: newOwner,
	}
}
