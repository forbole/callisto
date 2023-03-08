package source

import (
	ccvconsumertypes "github.com/cosmos/interchain-security/x/ccv/consumer/types"
)

type Source interface {
	GetNextFeeDistribution(height int64) (*ccvconsumertypes.NextFeeDistributionEstimate, error)
}
