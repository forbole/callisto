package source

import (
	ccvprovidertypes "github.com/cosmos/interchain-security/x/ccv/provider/types"
)

type Source interface {
	GetAllConsumerChains(height int64) ([]*ccvprovidertypes.Chain, error)
}
