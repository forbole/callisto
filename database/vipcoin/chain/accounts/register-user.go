/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package accounts

import (
	"context"
	"database/sql"

	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/forbole/bdjuno/v2/database/types"
)

func (r Repository) SaveRegisterUser(msg ...*accountstypes.MsgRegisterUser) error {
	if len(msg) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `INSERT INTO vipcoin_chain_accounts_register_user 
			(creator, address, hash, public_key, holder_wallet, ref_reward_wallet, 
			holder_wallet_extras, ref_reward_wallet_extras, referrer_hash) 
		VALUES 
			(:creator, :address, :hash, :public_key, :holder_wallet, :ref_reward_wallet, 
			:holder_wallet_extras, :ref_reward_wallet_extras, :referrer_hash)`

	for _, user := range msg {
		if _, err := tx.NamedExec(query, toRegisterUserDatabase(user)); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r Repository) GetRegisterUser(accfilter filter.Filter) ([]*accountstypes.MsgRegisterUser, error) {
	query, args := accfilter.Build("vipcoin_chain_accounts_register_user",
		`creator, address, hash, public_key, holder_wallet, ref_reward_wallet, 
		holder_wallet_extras, ref_reward_wallet_extras, referrer_hash`)

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
