package source

import (
	superfluidtypes "github.com/osmosis-labs/osmosis/v15/x/superfluid/types"
)

type Source interface {
	GetSuperfluidDelegationsByDelegator(address string, height int64) ([]superfluidtypes.SuperfluidDelegationRecord, error)
}
