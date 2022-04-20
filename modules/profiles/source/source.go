package source

import profilestypes "github.com/desmos-labs/desmos/v3/x/profiles/types"

type Source interface {
	GetParams(height int64) (profilestypes.Params, error)
}
