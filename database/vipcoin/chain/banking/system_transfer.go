package banking

import (
	"context"
	"database/sql"

	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveSystemTransfers - method that create system transfers to the "vipcoin_chain_banking_system_transfer" table
func (r Repository) SaveSystemTransfers(transfers ...*bankingtypes.SystemTransfer) error {
	if len(transfers) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	queryBaseTransfer := `INSERT INTO vipcoin_chain_banking_base_transfers 
       ("asset", "amount", "kind", "extras", "timestamp", "tx_hash") 
     VALUES 
       (:asset, :amount, :kind, :extras, :timestamp, :tx_hash)
     RETURNING id`

	querySystemTransfer := `INSERT INTO vipcoin_chain_banking_system_transfer
			("id", "wallet_from", "wallet_to")
			VALUES
			(:id, :wallet_from, :wallet_to)`

	for _, transfer := range transfers {
		transferDB := toSystemTransferDatabase(transfer)

		resp, err := tx.NamedQuery(queryBaseTransfer, transferDB)
		if err != nil {
			return err
		}

		for resp.Next() {
			if err := resp.Scan(&transferDB.ID); err != nil {
				return err
			}
		}

		if err := resp.Err(); err != nil {
			return err
		}

		if _, err := tx.NamedExec(querySystemTransfer, transferDB); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetSystemTransfers - method that get system transfers from the "vipcoin_chain_banking_system_transfer" table
func (r Repository) GetSystemTransfers(filter filter.Filter) ([]*bankingtypes.SystemTransfer, error) {
	query, args := filter.ToJoiner().
		PrepareTable(tableTransfers, types.FieldID, types.FieldAsset, types.FieldAmount, types.FieldKind, types.FieldExtras, types.FieldTimestamp, types.FieldTxHash).
		PrepareTable(tableSystemTransfer, types.FieldID, types.FieldWalletFrom, types.FieldWalletTo).
		PrepareJoinStatement("INNER JOIN vipcoin_chain_banking_base_transfers on vipcoin_chain_banking_base_transfers.id = vipcoin_chain_banking_system_transfer.id").
		Build(tableSystemTransfer)

	var transfersDB []types.DBSystemTransfer
	if err := r.db.Select(&transfersDB, query, args...); err != nil {
		return []*bankingtypes.SystemTransfer{}, err
	}

	result := make([]*bankingtypes.SystemTransfer, 0, len(transfersDB))
	for _, transfer := range transfersDB {
		result = append(result, toSystemTransferDomain(transfer))
	}

	return result, nil
}

// UpdateSystemTransfers - method that update the transfer in the "vipcoin_chain_banking_system_transfer" table
func (r Repository) UpdateSystemTransfers(transfers ...*bankingtypes.SystemTransfer) error {
	if len(transfers) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	queryBaseTransfer := `UPDATE vipcoin_chain_banking_base_transfers SET
	asset =:asset, amount =:amount, kind =:kind, extras =:extras, timestamp =:timestamp, tx_hash =:tx_hash
	WHERE id =:id;
	`

	queryTransfer := `UPDATE vipcoin_chain_banking_system_transfer SET
	wallet_from =:wallet_from, wallet_to =:wallet_to
	WHERE id =:id;
	`

	for _, transfer := range transfers {
		transferDB := toSystemTransfersDatabase(transfer)

		if _, err := tx.NamedExec(queryBaseTransfer, transferDB); err != nil {
			return err
		}

		if _, err := tx.NamedExec(queryTransfer, transferDB); err != nil {
			return err
		}
	}

	return tx.Commit()
}
