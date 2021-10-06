package source

import (
	oracletypes "github.com/bandprotocol/chain/v2/x/oracle/types"
)

type Source interface {
	GetParams(height int64) (oracletypes.Params, error)
}
