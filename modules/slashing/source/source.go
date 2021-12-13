package source

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

type Source interface {
	GetSigningInfo(height int64, consAddr sdk.ConsAddress) (slashingtypes.ValidatorSigningInfo, error)
	GetSigningInfos(height int64) ([]slashingtypes.ValidatorSigningInfo, error)
	GetParams(height int64) (slashingtypes.Params, error)
}
