package feeexcluder

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// BLOCK GenesisState
func toGenesisStateDomain(
	addresses []types.Address,
	dailyStats []types.DailyStats,
	stats []types.Stats,
	tariffs []types.Tariffs,
	g db.FeeExcluderGenesisState,
) types.GenesisState {
	return types.GenesisState{
		Params:          types.Params{}, // it is ok, params are always empty
		AddressList:     addresses,
		AddressCount:    g.AddressCount,
		DailyStatsList:  dailyStats,
		DailyStatsCount: g.DailyStatsCount,
		StatsList:       stats,
		TariffsList:     tariffs,
	}
}

func toGenesisStateDatabase(id uint64, g types.GenesisState) db.FeeExcluderGenesisState {
	return db.FeeExcluderGenesisState{
		ID:              id,
		AddressCount:    g.AddressCount,
		DailyStatsCount: g.DailyStatsCount,
	}
}
