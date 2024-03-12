package feeexcluder

import (
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"

	"github.com/forbole/bdjuno/v4/database/overgold/chain"
	"github.com/forbole/bdjuno/v4/database/types"
)

// GetAllGenesisState - method that get data from a db (overgold_feeexcluder_genesis_state).
// TODO: use JOIN and other db model
func (r Repository) GetAllGenesisState(f filter.Filter) ([]fe.GenesisState, error) {
	q, args := f.Build(tableGenesisState)

	// 1) get genesis state
	var gsList []types.FeeExcluderGenesisState
	if err := r.db.Select(&gsList, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableGenesisState}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(gsList) == 0 {
		return nil, errs.NotFound{What: tableGenesisState}
	}

	result := make([]fe.GenesisState, 0, len(gsList))
	for _, gs := range gsList {
		// 2) get m2m genesis state address and ids
		m2mAddress, err := r.GetAllM2MGenesisStateAddress(filter.NewFilter().
			SetArgument(types.FieldGenesisStateID, gs.ID))
		if err != nil {
			return nil, err
		}

		addressIDs := make([]uint64, 0, len(m2mAddress))
		for _, m2m := range m2mAddress {
			addressIDs = append(addressIDs, m2m.AddressID)
		}

		// 3) get m2m genesis state daily stats and ids
		m2mDailyStats, err := r.GetAllM2MGenesisStateDailyStats(filter.NewFilter().
			SetArgument(types.FieldGenesisStateID, gs.ID))
		if err != nil {
			return nil, err
		}

		dailyStatsIDs := make([]uint64, 0, len(m2mDailyStats))
		for _, m2m := range m2mDailyStats {
			dailyStatsIDs = append(dailyStatsIDs, m2m.DailyStatsID)
		}

		// 4) get m2m genesis state daily stats and ids
		m2mStats, err := r.GetAllM2MGenesisStateStats(filter.NewFilter().
			SetArgument(types.FieldGenesisStateID, gs.ID))
		if err != nil {
			return nil, err
		}

		statsIDs := make([]string, 0, len(m2mStats))
		for _, m2m := range m2mStats {
			statsIDs = append(statsIDs, m2m.StatsID)
		}

		// 5) get m2m genesis state daily stats and ids
		m2mTariffs, err := r.GetAllM2MGenesisStateTariffs(filter.NewFilter().
			SetArgument(types.FieldGenesisStateID, gs.ID))
		if err != nil {
			return nil, err
		}

		tariffsIDs := make([]uint64, 0, len(m2mTariffs))
		for _, m2m := range m2mTariffs {
			tariffsIDs = append(tariffsIDs, m2m.TariffsID)
		}

		// 6) get address
		addressList, err := r.GetAllAddress(filter.NewFilter().SetArgument(types.FieldID, addressIDs))
		if err != nil {
			return nil, err
		}

		// 7) get daily stats
		dailyStatsList, err := r.GetAllDailyStats(filter.NewFilter().SetArgument(types.FieldID, dailyStatsIDs))
		if err != nil {
			return nil, err
		}

		// 8) get daily stats
		statsList, err := r.GetAllStats(filter.NewFilter().SetArgument(types.FieldID, statsIDs))
		if err != nil {
			return nil, err
		}

		// 9) get tariff
		tariffsList, err := r.GetAllTariffs(filter.NewFilter().SetArgument(types.FieldID, tariffsIDs))
		if err != nil {
			return nil, err
		}

		result = append(result, toGenesisStateDomain(addressList, dailyStatsList, statsList, tariffsList, gs))
	}

	return result, nil
}

// InsertToGenesisState - insert new data in a database (overgold_feeexcluder_genesis_state).
func (r Repository) InsertToGenesisState(gs fe.GenesisState) error {
	var genesisStateID uint64
	m2mAddreses := make([]types.FeeExcluderM2MGenesisStateAddress, 0, len(gs.AddressList))
	m2mDailyStats := make([]types.FeeExcluderM2MGenesisStateDailyStats, 0, len(gs.DailyStatsList))
	m2mStats := make([]types.FeeExcluderM2MGenesisStateStats, 0, len(gs.StatsList))
	m2mTariffs := make([]types.FeeExcluderM2MGenesisStateTariffs, 0, len(gs.TariffsList))

	// 1) insert genesis state
	q := `
		INSERT INTO overgold_feeexcluder_genesis_state (
			address_count, daily_stats_count
		) VALUES (
			$1, $2
		) RETURNING id
	`

	m := toGenesisStateDatabase(genesisStateID, gs)
	if err := r.db.QueryRowx(q, m.AddressCount, m.DailyStatsCount).Scan(&genesisStateID); err != nil {
		if !chain.IsAlreadyExists(err) {
			return errs.Internal{Cause: err.Error()}
		}
	}

	// 2) insert address
	for _, a := range gs.AddressList {
		id, err := r.InsertToAddress(nil, a)
		if err != nil {
			return err
		}
		if id == 0 { // skip if already exists
			continue
		}

		m2mAddreses = append(m2mAddreses, types.FeeExcluderM2MGenesisStateAddress{
			GenesisStateID: genesisStateID,
			AddressID:      id,
		})
	}

	// 3) insert daily stats
	for _, d := range gs.DailyStatsList {
		id, err := r.InsertToDailyStats(nil, d)
		if err != nil {
			return err
		}
		if id == 0 { // skip if already exists
			continue
		}

		m2mDailyStats = append(m2mDailyStats, types.FeeExcluderM2MGenesisStateDailyStats{
			GenesisStateID: genesisStateID,
			DailyStatsID:   id,
		})
	}

	// 4) insert stats
	for _, s := range gs.StatsList {
		id, err := r.InsertToStats(nil, s)
		if err != nil {
			return err
		}
		if id == "" { // skip if already exists
			continue
		}

		m2mStats = append(m2mStats, types.FeeExcluderM2MGenesisStateStats{
			GenesisStateID: genesisStateID,
			StatsID:        id,
		})
	}

	// 5) insert tariffs
	for _, t := range gs.TariffsList {
		id, err := r.InsertToTariffs(nil, t)
		if err != nil {
			return err
		}
		if id == 0 { // skip if already exists
			continue
		}

		m2mTariffs = append(m2mTariffs, types.FeeExcluderM2MGenesisStateTariffs{
			GenesisStateID: genesisStateID,
			TariffsID:      id,
		})
	}

	// 6) insert m2m genesis state address
	if err := r.InsertToM2MGenesisStateAddress(nil, m2mAddreses...); err != nil {
		return err
	}

	// 7) insert m2m genesis state daily stats
	if err := r.InsertToM2MGenesisStateDailyStats(nil, m2mDailyStats...); err != nil {
		return err
	}

	// 8) insert m2m genesis state stats
	if err := r.InsertToM2MGenesisStateStats(nil, m2mStats...); err != nil {
		return err
	}

	// 9) insert m2m genesis state tariffs
	if err := r.InsertToM2MGenesisStateTariffs(nil, m2mTariffs...); err != nil {
		return err
	}

	return nil
}

// DeleteGenesisState - method that deletes data in a database (overgold_feeexcluder_genesis_state).
func (r Repository) DeleteGenesisState(id uint64) error {
	gsFilter := filter.NewFilter().SetArgument(types.FieldGenesisStateID, id)

	// 1) delete m2m genesis state address
	m2mAddress, err := r.GetAllM2MGenesisStateAddress(gsFilter)
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}

	if err = r.DeleteM2MGenesisStateAddressByGenesisState(nil, id); err != nil {
		return err
	}

	// 2) delete m2m genesis state daily stats
	m2mDailyStats, err := r.GetAllM2MGenesisStateDailyStats(gsFilter)
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}

	if err = r.DeleteM2MGenesisStateDailyStatsByGenesisState(nil, id); err != nil {
		return err
	}

	// 3) delete m2m genesis state stats
	m2mStats, err := r.GetAllM2MGenesisStateStats(gsFilter)
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}

	if err = r.DeleteM2MGenesisStateStatsByGenesisState(nil, id); err != nil {
		return err
	}

	// 4) delete m2m genesis state tariffs
	m2mTariffs, err := r.GetAllM2MGenesisStateTariffs(gsFilter)
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}

	if err = r.DeleteM2MGenesisStateTariffsByGenesisState(nil, id); err != nil {
		return err
	}

	// 5) delete address
	for _, m := range m2mAddress {
		if err = r.DeleteAddress(nil, m.AddressID); err != nil {
			return err
		}
	}

	// 6) delete daily stats
	for _, m := range m2mDailyStats {
		if err = r.DeleteDailyStats(nil, m.DailyStatsID); err != nil {
			return err
		}
	}

	// 7) delete daily stats
	for _, m := range m2mStats {
		if err = r.DeleteStats(nil, m.StatsID); err != nil {
			return err
		}
	}

	// 8) delete tariffs
	for _, m := range m2mTariffs {
		if err = r.DeleteTariffs(nil, m.TariffsID); err != nil {
			return err
		}
	}

	// 9) delete genesis state
	q := `DELETE FROM overgold_feeexcluder_genesis_state WHERE id IN ($1)`

	if _, err = r.db.Exec(q, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}
