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

// GetAllUpdateAddresses - method that get data from a db (overgold_allowed_update_addresses).
func (r Repository) GetAllUpdateAddresses(filter filter.Filter) ([]allowed.MsgUpdateAddresses, error) {
	q, args := filter.Build(tableUpdateAddresses)

	var result []types.AllowedUpdateAddresses
	if err := r.db.Select(&result, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableUpdateAddresses}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableUpdateAddresses}
	}

	return toUpdateAddressesDomainList(result), nil
}

// InsertToUpdateAddresses - insert a new MsgUpdateAddresses in a database (overgold_allowed_update_addresses).
func (r Repository) InsertToUpdateAddresses(hash string, msgs ...*allowed.MsgUpdateAddresses) error {
	if len(msgs) == 0 || hash == "" {
		return nil
	}

	q := `
		INSERT INTO overgold_allowed_update_addresses (
			id, tx_hash, creator, address
		) VALUES (
			$1, $2, $3, $4
		) RETURNING
			id, tx_hash, creator, address
	`

	for _, msg := range msgs {
		m := toUpdateAddressesDatabase(hash, msg)
		if _, err := r.db.Exec(q, m.ID, m.TxHash, m.Creator, m.Address); err != nil {
			if chain.IsAlreadyExists(err) {
				continue
			}

			return errs.Internal{Cause: err.Error()}
		}
	}

	return nil
}
