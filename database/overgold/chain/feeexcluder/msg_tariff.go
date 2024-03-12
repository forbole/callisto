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

// GetAllTariff - method that get data from a db (overgold_feeexcluder_tariff).
// TODO: use JOIN and other db model, e.g.:
//
//	SELECT f.*, mt.*
//	FROM overgold_feeexcluder_fees AS f
//	JOIN overgold_feeexcluder_m2m_tariff_fees AS mt ON f.id = mt.fees_id
//	JOIN overgold_feeexcluder_tariff AS t ON mt.tariff_id = t.id;
func (r Repository) GetAllTariff(f filter.Filter) ([]*fe.Tariff, error) {
	q, args := f.Build(tableTariff)

	// 1) get tariff
	var tariffs []types.FeeExcluderTariff
	if err := r.db.Select(&tariffs, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableTariff}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(tariffs) == 0 {
		return nil, errs.NotFound{What: tableTariff}
	}

	// TODO: refactor, use JOIN and other db model, e.g: overgold_feeexcluder_m2m_tariff_fees.fees_id...
	result := make([]*fe.Tariff, 0, len(tariffs))
	for _, t := range tariffs {
		// 2) get m2m tariff fees
		m2mFees, err := r.GetAllM2MTariffFees(filter.NewFilter().SetArgument(types.FieldTariffID, t.ID))
		if err != nil {
			return nil, err
		}

		feeIDs := make([]uint64, 0, len(m2mFees))
		for _, m2m := range m2mFees {
			feeIDs = append(feeIDs, m2m.FeesID)
		}

		// 3) get fees
		fees, err := r.GetAllFees(filter.NewFilter().SetArgument(types.FieldID, feeIDs))
		if err != nil {
			return nil, err
		}

		result = append(result, toTariffDomain(t, fees))
	}

	return result, nil
}

// InsertToTariff - insert new data in a database (overgold_feeexcluder_tariff).
func (r Repository) InsertToTariff(_ *sqlx.Tx, tariff *fe.Tariff) (lastID uint64, err error) {
	// 1) add tariff
	q := `
		INSERT INTO overgold_feeexcluder_tariff (
			msg_id, amount, denom, min_ref_balance
		) VALUES (
			$1, $2, $3, $4
		) RETURNING id
	`

	m, err := toTariffDatabase(0, tariff)
	if err != nil {
		return 0, errs.Internal{Cause: err.Error()}
	}

	if err = r.db.QueryRowx(q, m.MsgID, m.Amount, m.Denom, m.MinRefBalance).Scan(&lastID); err != nil {
		if chain.IsAlreadyExists(err) {
			return 0, nil
		}
		return 0, errs.Internal{Cause: err.Error()}
	}

	// 2) add fees and save unique ids
	feesIDs := make([]uint64, 0, len(tariff.Fees))
	for _, f := range tariff.Fees {
		id, err := r.InsertToFees(nil, f)
		if err != nil {
			return 0, err
		}

		if id == 0 {
			continue
		}

		feesIDs = append(feesIDs, id)
	}

	// 3) add many-to-many tariff fees
	m2m := make([]types.FeeExcluderM2MTariffFees, 0, len(tariff.Fees))
	for _, id := range feesIDs {
		m2m = append(m2m, types.FeeExcluderM2MTariffFees{
			TariffID: lastID,
			FeesID:   id,
		})
	}

	return lastID, r.InsertToM2MTariffFees(nil, m2m...)
}

// UpdateTariff - method that updates in a database (overgold_feeexcluder_tariff).
func (r Repository) UpdateTariff(_ *sqlx.Tx, id uint64, tariff *fe.Tariff) (err error) {
	// 1) update tariff
	q := `UPDATE overgold_feeexcluder_tariff SET
                 msg_id = $1,
				 amount = $2,
				 denom = $3,
				 min_ref_balance = $4
			 WHERE id = $5`

	m, err := toTariffDatabase(id, tariff)
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	if _, err = r.db.Exec(q, m.MsgID, m.Amount, m.Denom, m.MinRefBalance, m.ID); err != nil {
		return err
	}

	// 2) get fees ids (custom unique ids and msg ids)
	m2mFees, err := r.GetAllM2MTariffFees(filter.NewFilter().SetArgument(types.FieldTariffID, id))
	if err != nil {
		return err
	}

	feesIDs := make([]uint64, 0, len(m2mFees))
	for _, m2m := range m2mFees {
		feesIDs = append(feesIDs, m2m.FeesID)
	}

	fees, err := r.getAllFeesWithUniqueID(filter.NewFilter().SetArgument(types.FieldID, feesIDs))
	if err != nil {
		return err
	}

	// 3) update fees
	for _, f := range fees {
		for _, msgFee := range tariff.Fees {
			if f.MsgID == msgFee.Id {
				if err = r.UpdateFees(nil, f.ID, msgFee); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// DeleteTariff - method that deletes data in a database (overgold_feeexcluder_tariff).
func (r Repository) DeleteTariff(_ *sqlx.Tx, id uint64) (err error) {
	// 1) delete many-to-many tariff fees and get ids
	m2m, err := r.GetAllM2MTariffFees(filter.NewFilter().SetArgument(types.FieldTariffID, id))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}

	if err = r.DeleteM2MTariffFeesByTariff(nil, id); err != nil {
		return err
	}

	// 2) delete fees
	for _, m := range m2m {
		if err = r.DeleteFees(nil, m.FeesID); err != nil {
			return err
		}
	}

	// 3) delete tariff
	q := `DELETE FROM overgold_feeexcluder_tariff WHERE id IN ($1)`

	if _, err = r.db.Exec(q, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// getTariffWithUniqueID - method that get data from a db (overgold_feeexcluder_tariffs).
func (r Repository) getTariffWithUniqueID(req filter.Filter) (types.FeeExcluderTariff, error) {
	query, args := req.SetLimit(1).Build(tableTariff)

	var result types.FeeExcluderTariff
	if err := r.db.GetContext(context.Background(), &result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.FeeExcluderTariff{}, errs.NotFound{What: tableTariff}
		}

		return types.FeeExcluderTariff{}, errs.Internal{Cause: err.Error()}
	}

	return result, nil
}

// getAllTariffWithUniqueID - method that get data from a db (overgold_feeexcluder_tariffs).
func (r Repository) getAllTariffWithUniqueID(f filter.Filter) ([]types.FeeExcluderTariff, error) {
	q, args := f.Build(tableTariff)

	var tariff []types.FeeExcluderTariff
	if err := r.db.Select(&tariff, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableTariff}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(tariff) == 0 {
		return nil, errs.NotFound{What: tableTariff}
	}

	return tariff, nil
}
