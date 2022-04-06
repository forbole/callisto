package banking

import (
	"context"
	"database/sql"

	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v2/database/types"
)

type (
	// Repository - defines a repository for banking repository
	Repository struct {
		db  *sqlx.DB
		cdc codec.Marshaler
	}
)

// NewRepository constructor
func NewRepository(db *sqlx.DB, cdc codec.Marshaler) *Repository {
	return &Repository{
		db:  db,
		cdc: cdc,
	}
}

// SaveBaseTransfers - method that create transfers to the "vipcoin_chain_banking_transfers" table
func (r Repository) SaveBaseTransfers(transfers ...*bankingtypes.BaseTransfer) error {
	if len(transfers) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_banking_base_transfers 
       ("id", "asset", "amount", "kind", "extras", "timestamp", "tx_hash") 
     VALUES 
       (:id, :asset, :amount, :kind, :extras, :timestamp, :tx_hash)`

	if _, err := r.db.NamedExec(query, toTransfersDatabase(transfers...)); err != nil {
		return err
	}

	return nil
}

// UpdateBaseTransfers - method that updates the transfers in the "vipcoin_chain_banking_transfers" table
func (r Repository) UpdateBaseTransfers(transfers ...*bankingtypes.BaseTransfer) error {
	if len(transfers) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := `UPDATE vipcoin_chain_banking_base_transfers SET
				asset = :asset, amount = :amount,
				kind = :kind, extras = :extras, timestamp = :timestamp,
				tx_hash = :tx_hash
			 WHERE id = :id`

	for _, transfer := range transfers {
		if _, err := tx.NamedExec(query, toTransferDatabase(transfer)); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetBaseTransfers - method that get transfers from the "vipcoin_chain_banking_transfers" table
func (r Repository) GetBaseTransfers(filter filter.Filter) ([]*bankingtypes.BaseTransfer, error) {
	query, args := filter.Build(
		tableTransfers,
		types.FieldID, types.FieldAsset, types.FieldAmount, types.FieldKind,
		types.FieldExtras, types.FieldTimestamp, types.FieldTxHash,
	)

	var result []types.DBTransfer
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*bankingtypes.BaseTransfer{}, err
	}

	transfers := make([]*bankingtypes.BaseTransfer, 0, len(result))
	for _, transfer := range result {
		transfers = append(transfers, toTransferDomain(transfer))
	}

	return transfers, nil
}
