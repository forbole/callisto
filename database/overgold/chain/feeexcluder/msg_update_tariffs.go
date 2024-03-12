package feeexcluder

import (
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"

	"github.com/forbole/bdjuno/v4/database/types"
)

// GetAllMsgUpdateTariffs - method that get data from a db (overgold_feeexcluder_update_tariffs).
// TODO: use JOIN and other db model
func (r Repository) GetAllMsgUpdateTariffs(f filter.Filter) ([]fe.MsgUpdateTariffs, error) {
	q, args := f.Build(tableUpdateTariffs)

	// 1) get update tariffs
	var updateTariffs []types.FeeExcluderUpdateTariffs
	if err := r.db.Select(&updateTariffs, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableUpdateTariffs}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(updateTariffs) == 0 {
		return nil, errs.NotFound{What: tableUpdateTariffs}
	}

	// 2) get tariff
	result := make([]fe.MsgUpdateTariffs, 0, len(updateTariffs))
	for _, ut := range updateTariffs {
		tariff, err := r.GetAllTariff(filter.NewFilter().SetArgument(types.FieldID, ut.TariffID))
		if err != nil {
			return nil, err
		}
		if len(tariff) == 0 {
			return nil, errs.NotFound{What: tableTariff}
		}

		result = append(result, toMsgUpdateTariffsDomain(ut, tariff[0]))
	}

	return result, nil
}

// InsertToMsgUpdateTariffs - insert new data in a database (overgold_feeexcluder_update_tariffs).
func (r Repository) InsertToMsgUpdateTariffs(hash string, ut fe.MsgUpdateTariffs) error {
	// 1) add tariff
	tariffID, err := r.InsertToTariff(nil, ut.Tariff)
	if err != nil {
		return err
	}
	if tariffID == 0 {
		return nil
	}

	// 2) add update tariffs
	q := `
		INSERT INTO overgold_feeexcluder_update_tariffs (
			tx_hash, creator, denom, tariff_id
		) VALUES (
			$1, $2, $3, $4
		) RETURNING
			id, tx_hash, creator, denom, tariff_id
	`

	m := toMsgUpdateTariffsDatabase(hash, 0, tariffID, ut)
	if _, err = r.db.Exec(q, m.TxHash, m.Creator, m.Denom, m.TariffID); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// UpdateMsgUpdateTariffs - method that updates in a database (overgold_feeexcluder_update_tariffs).
func (r Repository) UpdateMsgUpdateTariffs(hash string, id uint64, ut fe.MsgUpdateTariffs) error {
	// 1) get unique tariff id
	tariff, err := r.getTariffWithUniqueID(filter.NewFilter().SetArgument(types.FieldMsgID, ut.Tariff.Id))
	if err != nil {
		return err
	}

	// 2) update update tariffs
	q := `UPDATE overgold_feeexcluder_update_tariffs SET
				 tx_hash = $1,
				 creator = $2,
            	 tariff_id = $3,
            	 denom = $4
			 WHERE id = $5`

	m := toMsgUpdateTariffsDatabase(hash, id, tariff.ID, ut)
	if _, err = r.db.Exec(q, m.TxHash, m.Creator, m.TariffID, m.Denom, m.ID); err != nil {
		return err
	}

	// 3) update tariff
	if err = r.UpdateTariff(nil, tariff.ID, ut.Tariff); err != nil {
		return err
	}

	return nil
}

// DeleteMsgUpdateTariffs - method that deletes data in a database (overgold_feeexcluder_update_tariffs).
func (r Repository) DeleteMsgUpdateTariffs(id uint64) error {
	q := `DELETE FROM overgold_feeexcluder_update_tariffs WHERE id IN ($1)`

	if _, err := r.db.Exec(q, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}
