package types

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtype "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Coin struct {
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

func ConvertCoins(coins sdk.Coins) []Coin {
	amount := make([]Coin, 0)
	for _, coin := range coins {
		amount = append(amount, Coin{Amount: coin.Amount.String(), Denom: coin.Denom})
	}
	return amount
}

func ConvertDecCoins(coins sdk.DecCoins) []Coin {
	amount := make([]Coin, 0)
	for _, coin := range coins {
		amount = append(amount, Coin{Amount: coin.Amount.String(), Denom: coin.Denom})
	}
	return amount
}

// ========================= Withdraw Address Response =========================

type Address struct {
	Address string `json:"address"`
}

// ========================= Account Balance Response =========================

type Balance struct {
	Coins []Coin `json:"coins"`
}

// ========================= Delegation Response =========================

type DelegationResponse struct {
	Delegations []Delegation        `json:"delegations"`
	Pagination  *query.PageResponse `json:"pagination"`
}

type Delegation struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Coins            []Coin `json:"coins"`
}

// ========================= Delegation Reward Response =========================

type DelegationReward struct {
	Coins            []Coin `json:"coins"`
	ValidatorAddress string `json:"validator_address"`
}

// ========================= Validator Commission Response =========================

type ValidatorCommissionAmount struct {
	Coins []Coin `json:"coins"`
}

// ========================= Unbonding Delegation Response =========================

type UnbondingDelegationResponse struct {
	UnbondingDelegations []UnbondingDelegation `json:"unbonding_delegations"`
	Pagination           *query.PageResponse   `json:"pagination"`
}

type UnbondingDelegation struct {
	DelegatorAddress string                                 `json:"delegator_address"`
	ValidatorAddress string                                 `json:"validator_address"`
	Entries          []stakingtype.UnbondingDelegationEntry `json:"entries"`
}

// ========================= Redelegation Response =========================

type RedelegationResponse struct {
	Redelegations []Redelegation      `json:"redelegations"`
	Pagination    *query.PageResponse `json:"pagination"`
}

type Redelegation struct {
	DelegatorAddress    string              `json:"delegator_address"`
	ValidatorSrcAddress string              `json:"validator_src_address"`
	ValidatorDstAddress string              `json:"validator_dst_address"`
	RedelegationEntries []RedelegationEntry `json:"entries"`
}

type RedelegationEntry struct {
	CompletionTime time.Time   `json:"completion_time"`
	Balance        sdkmath.Int `json:"balance"`
}
