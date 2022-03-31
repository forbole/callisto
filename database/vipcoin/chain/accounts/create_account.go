package accounts

import (
	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveCreateAccount - saves the given create account message inside the database
func (r Repository) SaveCreateAccount(msg ...*accountstypes.MsgCreateAccount) error {
	if len(msg) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_accounts_create_account 
			(creator, hash, address, public_key, kinds, state, extras) 
		VALUES 
			(:creator, :hash, :address, :public_key, :kinds, :state, :extras)`

	if _, err := r.db.NamedExec(query, toCreateAccountsDatabase(msg...)); err != nil {
		return err
	}

	return nil
}

// GetCreateAccount - get the given create account message from database
func (r Repository) GetCreateAccount(accfilter filter.Filter) ([]*accountstypes.MsgCreateAccount, error) {
	query, args := accfilter.Build(
		tableCreateAccount,
		types.FieldCreator, types.FieldHash, types.FieldAddress, types.FieldPublicKey,
		types.FieldKinds, types.FieldState, types.FieldExtras,
	)

	var result []types.DBCreateAccount
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*accountstypes.MsgCreateAccount{}, err
	}

	accounts := make([]*accountstypes.MsgCreateAccount, 0, len(result))
	for _, account := range result {
		accounts = append(accounts, toCreateAccountDomain(account))
	}

	return accounts, nil
}
