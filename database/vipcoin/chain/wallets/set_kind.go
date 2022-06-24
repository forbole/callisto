package wallets

import (
	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveKinds - inserts into the "vipcoin_chain_wallets_set_wallet_kind" table
func (r Repository) SaveKinds(messages *walletstypes.MsgSetWalletKind, transactionHash string) error {
	query := `INSERT INTO vipcoin_chain_wallets_set_wallet_kind 
			(transaction_hash, creator, address, kind) 
			VALUES 
			(:transaction_hash, :creator, :address, :kind)`

	if _, err := r.db.NamedExec(query, toSetKindDatabase(messages, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// DBSetWalletKind - get Kinds from the "vipcoin_chain_wallets_set_wallet_kind" table
func (r Repository) GetKinds(filter filter.Filter) ([]*walletstypes.MsgSetWalletKind, error) {
	query, args := filter.Build("vipcoin_chain_wallets_set_wallet_kind",
		types.FieldCreator, types.FieldAddress, types.FieldKind)

	var kindsDB []types.DBSetWalletKind
	if err := r.db.Select(&kindsDB, query, args...); err != nil {
		return []*walletstypes.MsgSetWalletKind{}, errs.Internal{Cause: err.Error()}
	}

	kinds := make([]*walletstypes.MsgSetWalletKind, 0, len(kindsDB))
	for _, kindDB := range kindsDB {
		kinds = append(kinds, toSetKindDomain(kindDB))
	}

	return kinds, nil
}
