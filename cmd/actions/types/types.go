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

type Coins struct {
	Coins sdk.Coins `json:"coins"`
}

// ========================= Delegation Reward =========================

type DelegatorReward struct {
	DecCoins   sdk.DecCoins `json:"dec_coins"`
	ValAddress string       `json:"validator_address"`
}

type DecCoins struct {
	DecCoins []sdk.DecCoin `json:"dec_coins"`
}

// ========================= Validator Commission  =========================

type ValidatorCommission struct {
	DecCoin     sdk.DecCoin  `json:"dec_coins"`
	ValAddress  string       `json:"validator_address"`
}

// ========================= Unbonding Delegation  =========================

type UnbondingDelegation struct {
	DelegatorAddress string                                  `json:"delegator_address"`
	ValidatorAddress string                                  `json:"validator_address"`
	Entries          []stakingtype.UnbondingDelegationEntry  `json:"entries"`
}