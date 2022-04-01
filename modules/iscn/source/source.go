package source

import (
	iscntypes "github.com/likecoin/likechain/x/iscn/types"
)

type Source interface {
	GetParams(height int64) (iscntypes.Params, error)
	GetRecordsByID(height int64, id string) (*iscntypes.QueryRecordsByIdResponse, error)
}
