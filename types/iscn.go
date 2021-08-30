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
	Records []iscntypes.QueryResponseRecord
	Height int64
}

func NewIscnRecord(
	records []iscntypes.QueryResponseRecord, 
	height int64) IscnRecord {
	return IscnRecord{
		Records: records,
		Height: height,
	}
}

