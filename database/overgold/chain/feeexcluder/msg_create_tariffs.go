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

// GetAllMsgCreateTariffs - method that get data from a db (overgold_feeexcluder_create_tariffs).
// TODO: use JOIN and other db model
func (r Repository) GetAllMsgCreateTariffs(f filter.Filter) ([]fe.MsgCreateTariffs, error) {
	q, args := f.Build(tableCreateTariffs)

	// 1) get create tariffs
	var createTariffs []types.FeeExcluderCreateTariffs
	if err := r.db.Select(&createTariffs, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableCreateTariffs}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(createTariffs) == 0 {
		return nil, errs.NotFound{What: tableCreateTariffs}
	}

	result := make([]fe.MsgCreateTariffs, 0, len(createTariffs))
	for _, ct := range createTariffs {
		// 2) get tariff
		tariff, err := r.GetAllTariff(filter.NewFilter().SetArgument(types.FieldID, ct.TariffID))
		if err != nil {
			return nil, err
		}
		if len(tariff) == 0 {
			return nil, errs.NotFound{What: tableTariff}
		}

		result = append(result, toMsgCreateTariffsDomain(ct, tariff[0]))
	}

	return result, nil
}

// InsertToMsgCreateTariffs - insert new data in a database (overgold_feeexcluder_create_tariffs).
func (r Repository) InsertToMsgCreateTariffs(hash string, ct fe.MsgCreateTariffs) error {
	// 1) add tariff
	tariffID, err := r.InsertToTariff(nil, ct.Tariff)
	if err != nil {
		return err
	}

	// 2) add create tariffs
	q := `
		INSERT INTO overgold_feeexcluder_create_tariffs (
			tx_hash, creator, denom, tariff_id
		) VALUES (
			$1, $2, $3, $4
		) RETURNING
			id, tx_hash, creator, denom, tariff_id
	`

	m := toMsgCreateTariffsDatabase(hash, 0, tariffID, ct)
	if _, err = r.db.Exec(q, m.TxHash, m.Creator, m.Denom, m.TariffID); err != nil {
		if chain.IsAlreadyExists(err) {
			return nil
		}
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}
