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

// GetAllM2MGenesisStateTariffs - method that get data from a db (overgold_feeexcluder_m2m_genesis_state_tariffs).
func (r Repository) GetAllM2MGenesisStateTariffs(filter filter.Filter) ([]types.FeeExcluderM2MGenesisStateTariffs, error) {
	q, args := filter.Build(tableM2MGenesisStateTariffs)

	var result []types.FeeExcluderM2MGenesisStateTariffs
	if err := r.db.Select(&result, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableM2MGenesisStateTariffs}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableM2MGenesisStateTariffs}
	}

	return result, nil
}

// InsertToM2MGenesisStateTariffs - insert new data in a database (overgold_feeexcluder_m2m_genesis_state_tariffs).
func (r Repository) InsertToM2MGenesisStateTariffs(_ *sqlx.Tx, t ...types.FeeExcluderM2MGenesisStateTariffs) (err error) {

	q := `
		INSERT INTO overgold_feeexcluder_m2m_genesis_state_tariffs (
			genesis_state_id, tariffs_id
		) VALUES (
			$1, $2
		) RETURNING
			genesis_state_id, tariffs_id
	`

	for _, m := range t {
		if _, err = r.db.Exec(q, m.GenesisStateID, m.TariffsID); err != nil {
			if chain.IsAlreadyExists(err) {
				continue
			}
			return errs.Internal{Cause: err.Error()}
		}
	}

	return nil
}

// DeleteM2MGenesisStateTariffsByGenesisState - method that deletes data in a database (overgold_feeexcluder_m2m_genesis_state_tariffs).
func (r Repository) DeleteM2MGenesisStateTariffsByGenesisState(tx *sqlx.Tx, id uint64) (err error) {
	q := `DELETE FROM overgold_feeexcluder_m2m_genesis_state_tariffs WHERE genesis_state_id IN ($1)`

	if _, err = tx.Exec(q, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}
