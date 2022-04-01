package accounts

import (
	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveRegisterUser - saves the given user inside the database
func (r Repository) SaveRegisterUser(msg ...*accountstypes.MsgRegisterUser) error {
	if len(msg) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_accounts_register_user 
			(creator, address, hash, public_key, holder_wallet, ref_reward_wallet, 
			holder_wallet_extras, ref_reward_wallet_extras, referrer_hash) 
		VALUES 
			(:creator, :address, :hash, :public_key, :holder_wallet, :ref_reward_wallet, 
			:holder_wallet_extras, :ref_reward_wallet_extras, :referrer_hash)`

	if _, err := r.db.NamedExec(query, toRegisterUsersDatabase(msg...)); err != nil {
		return err
	}

	return nil
}

// GetRegisterUser - get the given user from database
func (r Repository) GetRegisterUser(accountFilter filter.Filter) ([]*accountstypes.MsgRegisterUser, error) {
	query, args := accountFilter.Build(
		tableRegisterUser,
		types.FieldCreator, types.FieldAddress, types.FieldHash, types.FieldPublicKey,
		types.FieldHolderWallet, types.FieldRefRewardWallet, types.FieldHolderWalletExtras,
		types.FieldRefRewardWalletExtras, types.FieldReferrerHash,
	)

	var result []types.DBRegisterUser
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*accountstypes.MsgRegisterUser{}, err
	}

	users := make([]*accountstypes.MsgRegisterUser, 0, len(result))
	for _, user := range result {
		users = append(users, toRegisterUserDomain(user))
	}

	return users, nil
}
