package feeexcluder

import (
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v4/database/overgold/chain"
	"github.com/forbole/bdjuno/v4/database/types"
)

// GetAllM2MTariffTariffs - method that get data from a db (overgold_feeexcluder_m2m_tariff_tariffs).
func (r Repository) GetAllM2MTariffTariffs(filter filter.Filter) ([]types.FeeExcluderM2MTariffTariffs, error) {
	q, args := filter.Build(tableM2MTariffTariffs)

	var result []types.FeeExcluderM2MTariffTariffs
	if err := r.db.Select(&result, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableM2MTariffTariffs}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableM2MTariffTariffs}
	}

	return result, nil
}

// InsertToM2MTariffTariffs - insert new data in a database (overgold_feeexcluder_m2m_tariff_tariffs).
func (r Repository) InsertToM2MTariffTariffs(_ *sqlx.Tx, ids ...types.FeeExcluderM2MTariffTariffs) (err error) {
	if len(ids) == 0 {
		return nil
	}

	q := `
		INSERT INTO overgold_feeexcluder_m2m_tariff_tariffs (
			tariff_id, tariffs_id
		) VALUES (
			$1, $2
		) RETURNING
			tariff_id, tariffs_id
	`

	for _, m := range ids {
		if _, err = r.db.Exec(q, m.TariffID, m.TariffsID); err != nil {
			if chain.IsAlreadyExists(err) {
				continue
			}
			return errs.Internal{Cause: err.Error()}
		}
	}

	return nil
}

// DeleteM2MTariffTariffsByTariffs - method that deletes data in a database (overgold_feeexcluder_m2m_tariff_tariffs).
func (r Repository) DeleteM2MTariffTariffsByTariffs(tx *sqlx.Tx, id uint64) (err error) {
	q := `DELETE FROM overgold_feeexcluder_m2m_tariff_tariffs WHERE tariffs_id IN ($1)`

	if _, err = r.db.Exec(q, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}
