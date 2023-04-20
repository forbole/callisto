package source

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	superfluidtypes "github.com/osmosis-labs/osmosis/v15/x/superfluid/types"
)

type Source interface {
	GetSuperfluidDelegationsByDelegator(address string, height int64) ([]superfluidtypes.SuperfluidDelegationRecord, error)
	GetTotalSuperfluidDelegationsByDelegator(address string, height int64) (sdk.Coins, error)
}
