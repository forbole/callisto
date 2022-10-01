package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	margintypes "github.com/Sifchain/sifnode/x/margin/types"
)

// MarginParamsRow represents a single row inside the margin_params table
type MarginParamsRow struct {
	OneRowID bool   `db:"one_row_id"`
	Params   string `db:"params"`
	Height   int64  `db:"height"`
}

// NewMarginParamsRow builds a new MarginParamsRow instance
func NewMarginParamsRow(
	params string, height int64,
) MarginParamsRow {
	return MarginParamsRow{
		OneRowID: true,
		Params:   params,
		Height:   height,
	}
}

// Equal reports whether m and n represent the same table rows.
func (m MarginParamsRow) Equal(n MarginParamsRow) bool {
	return m.Params == n.Params &&
		m.Height == n.Height
}

// NewMarginParams builds a new sifchain x/margin Params instance
func NewMarginParams(leverageMax, interestRateMax, interestRateMin, interestRateIncrease, interestRateDecrease, healthGainFactor sdk.Dec,
	epochLength int64, pools []string, removalQueueThreshold sdk.Dec, maxOpenPositions uint64, poolOpenThreshold, forceCloseFundPercentage sdk.Dec,
	forceCloseFundAddress string, incrementalInterestPaymentFundPercentage sdk.Dec, incrementalInterestPaymentFundAddress string,
	sqModifier, safetyFactor sdk.Dec, closedPools []string, incrementalInterestPaymentEnabled, whitelistingEnabled bool) *margintypes.Params {
	return &margintypes.Params{
		LeverageMax:                              leverageMax,
		InterestRateMax:                          interestRateMax,
		InterestRateMin:                          interestRateMin,
		InterestRateIncrease:                     interestRateIncrease,
		InterestRateDecrease:                     interestRateDecrease,
		HealthGainFactor:                         healthGainFactor,
		EpochLength:                              epochLength,
		Pools:                                    pools,
		RemovalQueueThreshold:                    removalQueueThreshold,
		MaxOpenPositions:                         maxOpenPositions,
		PoolOpenThreshold:                        poolOpenThreshold,
		ForceCloseFundPercentage:                 forceCloseFundPercentage,
		ForceCloseFundAddress:                    forceCloseFundAddress,
		IncrementalInterestPaymentFundPercentage: incrementalInterestPaymentFundPercentage,
		IncrementalInterestPaymentFundAddress:    incrementalInterestPaymentFundAddress,
		SqModifier:                               sqModifier,
		SafetyFactor:                             safetyFactor,
		ClosedPools:                              closedPools,
		IncrementalInterestPaymentEnabled:        incrementalInterestPaymentEnabled,
		WhitelistingEnabled:                      whitelistingEnabled,
	}
}
