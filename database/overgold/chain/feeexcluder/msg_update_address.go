package feeexcluder

import (
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"

	"github.com/forbole/bdjuno/v4/database/types"
)

// GetAllMsgUpdateAddress - method that get data from a db (overgold_feeexcluder_update_address).
func (r Repository) GetAllMsgUpdateAddress(filter filter.Filter) ([]fe.MsgUpdateAddress, error) {
	q, args := filter.Build(tableUpdateAddress)

	var result []types.FeeExcluderUpdateAddress
	if err := r.db.Select(&result, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableUpdateAddress}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableUpdateAddress}
	}

	return toMsgUpdateAddressDomainList(result), nil
}

// InsertToMsgUpdateAddress - insert new data in a database (overgold_feeexcluder_update_address).
func (r Repository) InsertToMsgUpdateAddress(hash string, address fe.MsgUpdateAddress) error {
	q := `
		INSERT INTO overgold_feeexcluder_update_address (
			id, tx_hash, creator, address
		) VALUES (
			$1, $2, $3, $4
		) RETURNING
			id, tx_hash, creator, address
	`

	m := toMsgUpdateAddressDatabase(hash, address)
	if _, err := r.db.Exec(q, m.ID, m.TxHash, m.Creator, m.Address); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// UpdateMsgUpdateAddress - method that updates in a database (overgold_feeexcluder_update_address).
func (r Repository) UpdateMsgUpdateAddress(hash string, addresses ...fe.MsgUpdateAddress) error {
	if len(addresses) == 0 {
		return nil
	}

	q := `UPDATE overgold_feeexcluder_update_address SET
				 creator = $1,
				 address = $2
			 WHERE id = $3`

	for _, address := range addresses {
		m := toMsgUpdateAddressDatabase(hash, address)
		if _, err := r.db.Exec(q, m.Creator, m.Address, m.ID); err != nil {
			return err
		}
	}

	return nil
}

// DeleteMsgUpdateAddress - method that deletes data in a database (overgold_feeexcluder_update_address).
func (r Repository) DeleteMsgUpdateAddress(id uint64) error {
	q := `DELETE FROM overgold_feeexcluder_update_address WHERE id IN ($1)`

	if _, err := r.db.Exec(q, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}
