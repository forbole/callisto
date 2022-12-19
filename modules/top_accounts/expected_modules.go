package top_accounts

import (
	juno "github.com/forbole/juno/v3/types"
)

type BankModule interface {
	UpdateBalances(addresses []string, height int64) error
}

type DistrModule interface {
	RefreshDelegatorRewards(addresses []string, height int64) error
}

type StakingModule interface {
	RefreshDelegations(height int64, delegatorAddr string) error
	RefreshRedelegations(tx *juno.Tx, index int, delegatorAddr string) error
	// HandleMsgUndelegate(tx *juno.Tx, index int, msg *stakingtypes.MsgUndelegate) error
}
