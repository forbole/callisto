package accounts

import (
	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveRegisterUser - saves the given user inside the database
func (r Repository) SaveRegisterUser(msg *accountstypes.MsgRegisterUser, transactionHash string) error {
	query := `INSERT INTO overgold_chain_accounts_register_user 
			(transaction_hash, creator, address, hash, public_key, holder_wallet, ref_reward_wallet, 
			holder_wallet_extras, ref_reward_wallet_extras, referrer_hash) 
		VALUES 
			(:transaction_hash, :creator, :address, :hash, :public_key, :holder_wallet, :ref_reward_wallet, 
			:holder_wallet_extras, :ref_reward_wallet_extras, :referrer_hash)`

	if _, err := r.db.NamedExec(query, toRegisterUserDatabase(msg, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
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
		return []*accountstypes.MsgRegisterUser{}, errs.Internal{Cause: err.Error()}
	}

	users := make([]*accountstypes.MsgRegisterUser, 0, len(result))
	for _, user := range result {
		users = append(users, toRegisterUserDomain(user))
	}

	return users, nil
}
