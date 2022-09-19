package source

import margintypes "github.com/Sifchain/sifnode/x/margin/types"

type Source interface {
	GetParams(height int64) (*margintypes.Params, error)
}
