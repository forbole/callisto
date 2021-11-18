package source

import profilestypes "github.com/desmos-labs/desmos/x/profiles/types"

type Source interface {
	GetParams(height int64) (profilestypes.Params, error)
}
