package source

import (
	iscntypes "github.com/likecoin/likecoin-chain/v4/x/iscn/types"
)

type Source interface {
	GetParams(height int64) (iscntypes.Params, error)
	GetRecordsByID(height int64, id string) (*iscntypes.QueryRecordsByIdResponse, error)
}
