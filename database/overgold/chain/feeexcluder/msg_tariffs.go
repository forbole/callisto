package feeexcluder

import (
	"context"
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v4/database/overgold/chain"
	"github.com/forbole/bdjuno/v4/database/types"
)

// GetAllTariffs - method that get data from a db (overgold_feeexcluder_tariffs). TODO: use JOIN and other db model
func (r Repository) GetAllTariffs(f filter.Filter) ([]fe.Tariffs, error) {
	q, args := f.Build(tableTariffs)

	// 1) get tariffs
	var tariffs []types.FeeExcluderTariffs
	if err := r.db.Select(&tariffs, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableTariffs}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(tariffs) == 0 {
		return nil, errs.NotFound{What: tableTariffs}
	}

	// TODO: refactor, use JOIN and other db model, e.g: overgold_feeexcluder_m2m_tariff_fees.fees_id...
	result := make([]fe.Tariffs, 0, len(tariffs))
	for _, ts := range tariffs {
		// 2) get m2m tariff tariff
		m2mTariff, err := r.GetAllM2MTariffTariffs(filter.NewFilter().SetArgument(types.FieldTariffsID, ts.ID))
		if err != nil {
			return nil, err
		}

		tariffIDs := make([]uint64, 0, len(m2mTariff))
		for _, m2m := range m2mTariff {
			tariffIDs = append(tariffIDs, m2m.TariffID)
		}

		// 3) get tariff
		tariff, err := r.GetAllTariff(filter.NewFilter().SetArgument(types.FieldID, tariffIDs))
		if err != nil {
			return nil, err
		}

		result = append(result, toTariffsDomain(ts, tariff))
	}

	return result, nil
}

// GetTariffsDB - method that get data from a db without domain model (overgold_feeexcluder_tariffs).
func (r Repository) GetTariffsDB(req filter.Filter) (types.FeeExcluderTariffs, error) {
	query, args := req.SetLimit(1).Build(tableTariffs)

	var result types.FeeExcluderTariffs
	if err := r.db.GetContext(context.Background(), &result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.FeeExcluderTariffs{}, errs.NotFound{What: tableTariffs}
		}

		return types.FeeExcluderTariffs{}, errs.Internal{Cause: err.Error()}
	}

	return result, nil
}

// GetAllTariffsDB - method that get data from a db without domain model (overgold_feeexcluder_tariffs).
func (r Repository) GetAllTariffsDB(f filter.Filter) ([]types.FeeExcluderTariffs, error) {
	q, args := f.Build(tableTariffs)

	var tariffs []types.FeeExcluderTariffs
	if err := r.db.Select(&tariffs, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableTariffs}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(tariffs) == 0 {
		return nil, errs.NotFound{What: tableTariffs}
	}

	return tariffs, nil
}

// InsertToTariffs - insert new data in a database (overgold_feeexcluder_tariffs).
func (r Repository) InsertToTariffs(_ *sqlx.Tx, tariffs fe.Tariffs) (lastID uint64, err error) {
	// 1) add tariffs
	q := `
		INSERT INTO overgold_feeexcluder_tariffs (
			denom, creator
		) VALUES (
			$1, $2
		) RETURNING id
	`

	m := toTariffsDatabase(0, tariffs)
	if err = r.db.QueryRowx(q, m.Denom, m.Creator).Scan(&lastID); err != nil {
		if chain.IsAlreadyExists(err) {
			return 0, nil
		}
		return 0, errs.Internal{Cause: err.Error()}
	}

	// 2) add tariff and save unique ids
	tariffIDs := make([]uint64, 0, len(tariffs.Tariffs))
	for _, t := range tariffs.Tariffs {
		id, err := r.InsertToTariff(nil, t)
		if err != nil {
			return 0, err
		}

		tariffIDs = append(tariffIDs, id)
	}

	// 3) add many-to-many tariff tariffs
	m2m := make([]types.FeeExcluderM2MTariffTariffs, 0, len(tariffs.Tariffs))
	for _, id := range tariffIDs {
		m2m = append(m2m, types.FeeExcluderM2MTariffTariffs{
			TariffsID: lastID,
			TariffID:  id,
		})
	}

	return lastID, r.InsertToM2MTariffTariffs(nil, m2m...)
}

// UpdateTariffs - method that updates in a database (overgold_feeexcluder_tariffs).
func (r Repository) UpdateTariffs(_ *sqlx.Tx, id uint64, tariffs fe.Tariffs) (err error) {
	// 1) update tariffs
	q := `UPDATE overgold_feeexcluder_tariffs SET
				 denom = $1,
				 creator = $2
			 WHERE id = $3`

	m := toTariffsDatabase(id, tariffs)
	if _, err = r.db.Exec(q, m.Denom, m.Creator, m.ID); err != nil {
		return err
	}

	// 2) get unique id from many-to-many tariff tariffs
	m2m, err := r.GetAllM2MTariffTariffs(filter.NewFilter().SetArgument(types.FieldTariffsID, id))
	if err != nil {
		return err
	}
	tariffIDs := make([]uint64, 0, len(m2m))
	for _, tariffs := range m2m {
		tariffIDs = append(tariffIDs, tariffs.TariffID)
	}

	tariffList, err := r.getAllTariffWithUniqueID(filter.NewFilter().SetArgument(types.FieldID, tariffIDs))
	if err != nil {
		return err
	}

	// 3) update tariff
	for _, t := range tariffList {
		for _, ts := range tariffs.Tariffs {
			if err = r.UpdateTariff(nil, t.ID, ts); err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteTariffs - method that deletes data in a database (overgold_feeexcluder_tariffs).
func (r Repository) DeleteTariffs(_ *sqlx.Tx, id uint64) (err error) {

	// 1) delete many-to-many tariff tariffs and get ids
	m2m, err := r.GetAllM2MTariffTariffs(filter.NewFilter().SetArgument(types.FieldTariffsID, id))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}

	if err = r.DeleteM2MTariffTariffsByTariffs(nil, id); err != nil {
		return err
	}

	// 2) delete tariff
	for _, m := range m2m {
		if err = r.DeleteTariff(nil, m.TariffID); err != nil {
			return err
		}
	}

	// 3) delete tariffs
	q := `DELETE FROM overgold_feeexcluder_tariffs WHERE id IN ($1)`

	if _, err = r.db.Exec(q, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// getTariffsWithUniqueID - method that get data from a db (overgold_feeexcluder_tariffs).
func (r Repository) getTariffsWithUniqueID(req filter.Filter) (types.FeeExcluderTariffs, error) {
	query, args := req.SetLimit(1).Build(tableTariffs)

	var result types.FeeExcluderTariffs
	if err := r.db.GetContext(context.Background(), &result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.FeeExcluderTariffs{}, errs.NotFound{What: tableTariffs}
		}

		return types.FeeExcluderTariffs{}, errs.Internal{Cause: err.Error()}
	}

	return result, nil
}

// getAllTariffsWithUniqueID - method that get data from a db (overgold_feeexcluder_tariffs).
func (r Repository) getAllTariffsWithUniqueID(f filter.Filter) ([]types.FeeExcluderTariffs, error) {
	q, args := f.Build(tableTariffs)

	var tariffs []types.FeeExcluderTariffs
	if err := r.db.Select(&tariffs, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableTariffs}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(tariffs) == 0 {
		return nil, errs.NotFound{What: tableTariffs}
	}

	return tariffs, nil
}
