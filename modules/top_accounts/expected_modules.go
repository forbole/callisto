package top_accounts

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type BankModule interface {
	UpdateBalances(addresses []string, height int64) error
}

type DistrModule interface {
	RefreshDelegatorRewards(addresses []string, height int64) error
}

type StakingModule interface {
	HandleMsgDelegate(height int64, msg *stakingtypes.MsgDelegate) error
	// HandleMsgBeginRedelegate(tx *juno.Tx, index int, msg *stakingtypes.MsgBeginRedelegate) error
	// HandleMsgUndelegate(tx *juno.Tx, index int, msg *stakingtypes.MsgUndelegate) error
}
