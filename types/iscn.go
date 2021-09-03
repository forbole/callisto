package types

import (
	
	iscntypes "github.com/likecoin/likechain/x/iscn/types"
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
	RecordNotes string
	ContentFingerprints []string
	Stakeholders []iscntypes.IscnInput
	ContentMetadata iscntypes.IscnInput
}

func NewRecord(recordNotes string, contentFingerprints []string, 
	stakeholders []iscntypes.IscnInput, contentMetadata iscntypes.IscnInput) Record {
	return Record{
		RecordNotes: recordNotes,
		ContentFingerprints: contentFingerprints,
		Stakeholders: stakeholders,
		ContentMetadata: contentMetadata,
	}
}

// IscnRecord represents the x/iscn records
type IscnRecord struct {
	Owner string
	IscnId string
	LatestVersion uint64
	Ipld string
	Data Record
	Height int64
}

// NewIscnRecord allows to build a new IscnRecord instance
func NewIscnRecord(owner string, iscnId string, latestVersion uint64, ipld string, 
	data Record, height int64) IscnRecord {
	return IscnRecord{
		Owner: owner,
		IscnId: iscnId,
		LatestVersion: latestVersion,
		Ipld: ipld,
		Data: data,
		Height: height,
	}
}

type IscnChangeOwnership struct {
	From string
	IscnId string
	NewOwner string
}

// NewIscnChangeOwnership allows to build a new IscnChangeOwnership instance
func NewIscnChangeOwnership(from string, iscnId string, newOwner string) IscnChangeOwnership {
	return IscnChangeOwnership{
		From: from,
		IscnId: iscnId,
		NewOwner: newOwner,
	}
}
