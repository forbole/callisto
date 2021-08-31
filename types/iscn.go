package types

import (
	// "time"
	
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


// IscnParams represents the x/iscn records
type IscnRecord struct {
	Owner string
	LatestVersion uint64
	Records []iscntypes.QueryResponseRecord
	Height int64
}

func NewIscnRecord(owner string, latestVersion uint64, 
	records []iscntypes.QueryResponseRecord, 
	height int64) IscnRecord {
	return IscnRecord{
		Owner: owner,
		LatestVersion: latestVersion,
		Records: records,
		Height: height,
	}
}

