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

// GetAllFees - method that get data from a db (overgold_feeexcluder_fees).
func (r Repository) GetAllFees(filter filter.Filter) ([]*fe.Fees, error) {
	q, args := filter.Build(tableFees)

	var result []types.FeeExcluderFees
	if err := r.db.Select(&result, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableFees}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableFees}
	}

	return toFeesDomainList(result), nil
}

// InsertToFees - insert new data in a database (overgold_feeexcluder_fees).
func (r Repository) InsertToFees(_ *sqlx.Tx, fees *fe.Fees) (lastID uint64, err error) {
	q := `
		INSERT INTO overgold_feeexcluder_fees (
			msg_id, creator, amount_from, fee, ref_reward, stake_reward, min_amount, no_ref_reward
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) RETURNING id
	`

	m, err := toFeesDatabase(0, fees)
	if err != nil {
		return 0, errs.Internal{Cause: err.Error()}
	}

	if err = r.db.QueryRowx(q,
		m.MsgID,
		m.Creator,
		m.AmountFrom,
		m.Fee,
		m.RefReward,
		m.StakeReward,
		m.MinAmount,
		m.NoRefReward,
	).Scan(&lastID); err != nil {
		if chain.IsAlreadyExists(err) {
			return 0, nil
		}
		return 0, errs.Internal{Cause: err.Error()}
	}

	return lastID, nil
}

// UpdateFees - method that updates in a database (overgold_feeexcluder_fees).
func (r Repository) UpdateFees(_ *sqlx.Tx, id uint64, fees *fe.Fees) (err error) {
	q := `UPDATE overgold_feeexcluder_fees SET
                 msg_id = $1,
				 creator = $2,
				 amount_from = $3, 
				 fee = $4, 
				 ref_reward = $5, 
				 stake_reward = $6, 
				 min_amount = $7, 
				 no_ref_reward = $8
			 WHERE id = $9`

	m, err := toFeesDatabase(id, fees)
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	if _, err = r.db.Exec(q,
		m.MsgID,
		m.Creator,
		m.AmountFrom,
		m.Fee,
		m.RefReward,
		m.StakeReward,
		m.MinAmount,
		m.NoRefReward,
		m.ID,
	); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// DeleteFees - method that deletes data in a database (overgold_feeexcluder_fees).
func (r Repository) DeleteFees(_ *sqlx.Tx, id uint64) (err error) {
	q := `DELETE FROM overgold_feeexcluder_fees WHERE id IN ($1)`

	if _, err = r.db.Exec(q, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// getFeesWithUniqueID - method that get data from a db (overgold_feeexcluder_fees).
func (r Repository) getFeesWithUniqueID(req filter.Filter) (types.FeeExcluderFees, error) {
	query, args := req.SetLimit(1).Build(tableFees)

	var result types.FeeExcluderFees
	if err := r.db.GetContext(context.Background(), &result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.FeeExcluderFees{}, errs.NotFound{What: tableFees}
		}

		return types.FeeExcluderFees{}, errs.Internal{Cause: err.Error()}
	}

	return result, nil
}

// getAllFeesListWithUniqueID - method that get data from a db (overgold_feeexcluder_fees).
func (r Repository) getAllFeesWithUniqueID(f filter.Filter) ([]types.FeeExcluderFees, error) {
	q, args := f.Build(tableFees)

	var fees []types.FeeExcluderFees
	if err := r.db.Select(&fees, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableFees}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(fees) == 0 {
		return nil, errs.NotFound{What: tableFees}
	}

	return fees, nil
}
