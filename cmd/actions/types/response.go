package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtype "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// ========================= Account Balance Response =========================

type Balance struct {
	Coins sdk.Coins `json:"coins"`
}

// ========================= Delegation Response =========================

type DelegationResponse struct {
	Delegations []Delegation        `json:"delegations"`
	Pagination  *query.PageResponse `json:"pagination"`
}

type Delegation struct {
	DelegatorAddress string   `json:"delegator_address"`
	ValidatorAddress string   `json:"validator_address"`
	Coins            sdk.Coin `json:"coins"`
}

// ========================= Delegation Reward Response =========================

type DelegationReward struct {
	Coins            sdk.DecCoins `json:"coins"`
	ValidatorAddress string       `json:"validator_address"`
}

// ========================= Validator Commission Response =========================

type ValidatorCommissionAmount struct {
	Coins sdk.DecCoins `json:"coins"`
}

// ========================= Unbonding Delegation Response =========================

type UnbondingDelegation struct {
	DelegatorAddress string                                 `json:"delegator_address"`
	ValidatorAddress string                                 `json:"validator_address"`
	Entries          []stakingtype.UnbondingDelegationEntry `json:"entries"`
}

// ========================= Redelegation Response =========================

type Redelegation struct {
	DelegatorAddress    string                          `json:"delegator_address"`
	ValidatorSrcAddress string                          `json:"validator_src_address"`
	ValidatorDstAddress string                          `json:"validator_dst_address"`
	Entries             []stakingtype.RedelegationEntry `json:"entries"`
}
