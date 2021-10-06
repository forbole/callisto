package source

import (
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

type Source interface {
	GetSigningInfos(height int64) ([]slashingtypes.ValidatorSigningInfo, error)
	GetParams(height int64) (slashingtypes.Params, error)
}
