package source

import (
	ccvprovidertypes "github.com/cosmos/interchain-security/x/ccv/provider/types"
)

type Source interface {
	// GetParams(height int64) (ccvprovidertypes.Params, error)
	GetAllConsumerChains(height int64) ([]*ccvprovidertypes.Chain, error)
	GetConsumerChainStarts(height int64) (*ccvprovidertypes.ConsumerAdditionProposals, error)
	GetConsumerChainStops(height int64) (*ccvprovidertypes.ConsumerRemovalProposals, error)
}
