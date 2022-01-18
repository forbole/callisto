package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtype "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type GraphQLError struct {
	Message string `json:"message"`
}

// ========================= Address Payload =========================
type AddressPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            Address                `json:"input"`
}

type Address struct {
	Address string `json:"address"`
}

// ========================= Account Balance =========================
type AccountBalancePayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            AccountBalanceArgs     `json:"input"`
}

type AccountBalanceArgs struct {
	Address string `json:"address"`
	Height  int64  `json:"height"`
}

type Balance struct {
	Coins sdk.Coins `json:"coins"`
}

// ========================= Delegation =========================

type Delegation struct {
	DelAddress string   `json:"delegator_address"`
	ValAddress string   `json:"validator_address"`
	Coin       sdk.Coin `json:"coin"`
}

// ========================= Delegation Reward =========================

type DelegatorReward struct {
	DecCoins   sdk.DecCoins `json:"dec_coins"`
	ValAddress string       `json:"validator_address"`
}

// ========================= Validator Commission  =========================

type ValidatorCommission struct {
	DecCoin    sdk.DecCoin `json:"dec_coin"`
	ValAddress string      `json:"validator_address"`
}

// ========================= Unbonding Delegation  =========================

type UnbondingDelegation struct {
	DelegatorAddress string                                 `json:"delegator_address"`
	ValidatorAddress string                                 `json:"validator_address"`
	Entries          []stakingtype.UnbondingDelegationEntry `json:"entries"`
}

// ========================= Redelegation  =========================

type Redelegation struct {
	DelegatorAddress    string                          `json:"delegator_address"`
	ValidatorSrcAddress string                          `json:"validator_src_address"`
	ValidatorDstAddress string                          `json:"validator_dst_address"`
	Entries             []stakingtype.RedelegationEntry `json:"entries"`
}
