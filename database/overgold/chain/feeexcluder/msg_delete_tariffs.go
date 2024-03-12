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

// GetAllMsgDeleteTariffs - method that get data from a db (overgold_feeexcluder_delete_tariffs).
func (r Repository) GetAllMsgDeleteTariffs(f filter.Filter) ([]fe.MsgDeleteTariffs, error) {
	q, args := f.Build(tableDeleteTariffs)

	var deleteTariffs []types.FeeExcluderDeleteTariffs
	if err := r.db.Select(&deleteTariffs, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableDeleteTariffs}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(deleteTariffs) == 0 {
		return nil, errs.NotFound{What: tableDeleteTariffs}
	}

	return toMsgDeleteTariffsDomainList(deleteTariffs), nil
}

// InsertToMsgDeleteTariffs - insert new data in a database (overgold_feeexcluder_delete_tariffs).
func (r Repository) InsertToMsgDeleteTariffs(hash string, dt fe.MsgDeleteTariffs) error {
	// 1) get unique tariff id
	tariff, err := r.getTariffWithUniqueID(filter.NewFilter().SetArgument(types.FieldMsgID, dt.TariffID))
	if err != nil {
		return err
	}

	// 2) get unique fees id
	fees, err := r.getFeesWithUniqueID(filter.NewFilter().SetArgument(types.FieldMsgID, dt.FeeID))
	if err != nil {
		return err
	}

	// 3) insert delete tariffs
	q := `
		INSERT INTO overgold_feeexcluder_delete_tariffs (
			tx_hash, creator, denom, tariff_id, fees_id
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING
			id, tx_hash, creator, denom, tariff_id, fees_id
	`

	m, err := toMsgDeleteTariffsDatabase(hash, 0, dt)
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	if _, err := r.db.Exec(q, m.TxHash, m.Creator, m.Denom, tariff.ID, fees.ID); err != nil {
		if chain.IsAlreadyExists(err) {
			return nil
		}
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// UpdateMsgDeleteTariffs - method that deletes in a database (overgold_feeexcluder_delete_tariffs).
func (r Repository) UpdateMsgDeleteTariffs(hash string, id uint64, ut fe.MsgDeleteTariffs) error {
	q := `UPDATE overgold_feeexcluder_delete_tariffs SET
				 tx_hash = $1,
				 creator = $2,
            	 tariff_id = $3,
            	 denom = $4,
                 fees_id = $5
			 WHERE id = $6`

	m, err := toMsgDeleteTariffsDatabase(hash, id, ut)
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	if _, err := r.db.Exec(q, m.TxHash, m.Creator, m.TariffID, m.Denom, m.FeesID, m.ID); err != nil {
		return err
	}

	return nil
}

// DeleteMsgDeleteTariffs - method that deletes data in a database (overgold_feeexcluder_delete_tariffs).
func (r Repository) DeleteMsgDeleteTariffs(id uint64) error {
	q := `DELETE FROM overgold_feeexcluder_delete_tariffs WHERE id IN ($1)`

	if _, err := r.db.Exec(q, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}
