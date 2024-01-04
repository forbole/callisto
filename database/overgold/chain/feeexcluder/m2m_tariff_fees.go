package feeexcluder

import (
	"context"
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v4/database/types"
)

// GetAllM2MTariffFees - method that get data from a db (overgold_feeexcluder_m2m_tariff_fees).
func (r Repository) GetAllM2MTariffFees(filter filter.Filter) ([]types.FeeExcluderM2MTariffFees, error) {
	q, args := filter.Build(tableM2MTariffFees)

	var result []types.FeeExcluderM2MTariffFees
	if err := r.db.Select(&result, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableM2MTariffFees}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableM2MTariffFees}
	}

	return result, nil
}

// InsertToM2MTariffFees - insert new data in a database (overgold_feeexcluder_m2m_tariff_fees).
func (r Repository) InsertToM2MTariffFees(tx *sqlx.Tx, ids ...types.FeeExcluderM2MTariffFees) (err error) {
	if tx == nil {
		tx, err = r.db.BeginTxx(context.Background(), &sql.TxOptions{})
		if err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		defer commit(tx, err)
	}

	q := `
		INSERT INTO overgold_feeexcluder_m2m_tariff_fees (
			tariff_id, fees_id
		) VALUES (
			$1, $2
		) RETURNING
			tariff_id, fees_id
	`

	for _, m := range ids {
		if _, err = tx.Exec(q, m.TariffID, m.FeesID); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return nil
}

// DeleteM2MTariffFeesByTariff - method that deletes data in a database (overgold_feeexcluder_m2m_tariff_fees).
func (r Repository) DeleteM2MTariffFeesByTariff(tx *sqlx.Tx, tariffID uint64) (err error) {
	if tx == nil {
		tx, err = r.db.BeginTxx(context.Background(), &sql.TxOptions{})
		if err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		defer commit(tx, err)
	}

	q := `DELETE FROM overgold_feeexcluder_m2m_tariff_fees WHERE tariff_id IN ($1)`

	if _, err = tx.Exec(q, tariffID); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}
