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

// GetAllMsgCreateAddress - method that get data from a db (overgold_feeexcluder_create_address).
func (r Repository) GetAllMsgCreateAddress(filter filter.Filter) ([]fe.MsgCreateAddress, error) {
	q, args := filter.Build(tableCreateAddress)

	var result []types.FeeExcluderCreateAddress
	if err := r.db.Select(&result, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableCreateAddress}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableCreateAddress}
	}

	return toMsgCreateAddressDomainList(result), nil
}

// InsertToMsgCreateAddress - insert new data in a database (overgold_feeexcluder_create_address).
func (r Repository) InsertToMsgCreateAddress(hash string, address fe.MsgCreateAddress) error {
	// 1) add address
	if _, err := r.InsertToAddress(nil, fe.Address{
		Address: address.Address,
		Creator: address.Creator,
	}); err != nil {
		return err
	}

	// 2) add create tariffs
	q := `
		INSERT INTO overgold_feeexcluder_create_address (
			tx_hash, creator, address
		) VALUES (
			$1, $2, $3
		) RETURNING
			id, tx_hash, creator, address
	`

	m := toMsgCreateAddressDatabase(hash, 0, address)
	if _, err := r.db.Exec(q, m.TxHash, m.Creator, m.Address); err != nil {
		if chain.IsAlreadyExists(err) {
			return nil
		}
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// UpdateMsgCreateAddress - method that updates in a database (overgold_feeexcluder_create_address).
func (r Repository) UpdateMsgCreateAddress(hash string, id uint64, address fe.MsgCreateAddress) error {
	q := `UPDATE overgold_feeexcluder_create_address SET
				 tx_hash = $1,
				 creator = $2,
				 address = $3
			 WHERE id = $4`

	m := toMsgCreateAddressDatabase(hash, id, address)
	if _, err := r.db.Exec(q, m.TxHash, m.Creator, m.Address, m.ID); err != nil {
		return err
	}

	return nil
}

// DeleteMsgCreateAddress - method that deletes data in a database (overgold_feeexcluder_create_address).
func (r Repository) DeleteMsgCreateAddress(id uint64) error {
	q := `DELETE FROM overgold_feeexcluder_create_address WHERE id IN ($1)`

	if _, err := r.db.Exec(q, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}
