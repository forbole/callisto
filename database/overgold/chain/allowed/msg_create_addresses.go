package allowed

import (
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	allowed "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"

	"github.com/forbole/bdjuno/v4/database/overgold/chain"
	"github.com/forbole/bdjuno/v4/database/types"
)

// GetAllCreateAddresses - method that get data from a db (overgold_allowed_create_addresses).
func (r Repository) GetAllCreateAddresses(filter filter.Filter) ([]allowed.MsgCreateAddresses, error) {
	q, args := filter.Build(tableCreateAddresses)

	var result []types.AllowedCreateAddresses
	if err := r.db.Select(&result, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableCreateAddresses}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableCreateAddresses}
	}

	return toCreateAddressesDomainList(result), nil
}

// InsertToCreateAddresses - insert a new MsgCreateAddresses in a database (overgold_allowed_create_addresses).
func (r Repository) InsertToCreateAddresses(hash string, msgs ...*allowed.MsgCreateAddresses) error {
	if len(msgs) == 0 || hash == "" {
		return nil
	}

	q := `
		INSERT INTO overgold_allowed_create_addresses (
			tx_hash, creator, address
		) VALUES (
			$1, $2, $3
		) RETURNING
			id, tx_hash, creator, address
	`

	for _, msg := range msgs {
		m := toCreateAddressesDatabase(hash, msg)
		if _, err := r.db.Exec(q, m.TxHash, m.Creator, m.Address); err != nil {
			if chain.IsAlreadyExists(err) {
				continue
			}

			return errs.Internal{Cause: err.Error()}
		}
	}

	return nil
}
