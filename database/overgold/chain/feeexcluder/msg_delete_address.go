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

// GetAllMsgDeleteAddress - method that get data from a db (overgold_feeexcluder_delete_address).
func (r Repository) GetAllMsgDeleteAddress(filter filter.Filter) ([]fe.MsgDeleteAddress, error) {
	q, args := filter.Build(tableDeleteAddress)

	var result []types.FeeExcluderDeleteAddress
	if err := r.db.Select(&result, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableDeleteAddress}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableDeleteAddress}
	}

	return toMsgDeleteAddressDomainList(result), nil
}

// InsertToMsgDeleteAddress - insert new data in a database (overgold_feeexcluder_delete_address).
func (r Repository) InsertToMsgDeleteAddress(hash string, addresses ...fe.MsgDeleteAddress) error {
	if len(addresses) == 0 {
		return nil
	}

	q := `
		INSERT INTO overgold_feeexcluder_delete_address (
			id, tx_hash, creator
		) VALUES (
			$1, $2, $3
		) RETURNING
			id, tx_hash, creator
	`

	for _, a := range addresses {
		m := toMsgDeleteAddressDatabase(hash, a)
		if _, err := r.db.Exec(q, m.ID, m.TxHash, m.Creator); err != nil {
			if chain.IsAlreadyExists(err) {
				return nil
			}
			return errs.Internal{Cause: err.Error()}
		}
	}

	return nil
}

// UpdateMsgDeleteAddress - method that updates in a database (overgold_feeexcluder_delete_address).
func (r Repository) UpdateMsgDeleteAddress(hash string, addresses ...fe.MsgDeleteAddress) error {
	if len(addresses) == 0 {
		return nil
	}

	q := `UPDATE overgold_feeexcluder_delete_address SET
				 creator = $1
			 WHERE id = $2`

	for _, address := range addresses {
		m := toMsgDeleteAddressDatabase(hash, address)
		if _, err := r.db.Exec(q, m.Creator, m.ID); err != nil {
			return err
		}
	}

	return nil
}

// DeleteMsgDeleteAddress - method that deletes data in a database (overgold_feeexcluder_delete_address).
func (r Repository) DeleteMsgDeleteAddress(id uint64) error {
	q := `DELETE FROM overgold_feeexcluder_delete_address WHERE id IN ($1)`

	if _, err := r.db.Exec(q, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}
