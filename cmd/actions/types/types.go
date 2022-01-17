package types

import sdk "github.com/cosmos/cosmos-sdk/types"

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

// ========================= Account Balances Payload =========================
type AccountBalancesPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            AccountBalancesArgs    `json:"input"`
}

type AccountBalancesArgs struct {
	Address string `json:"address"`
	Height  int64  `json:"height"`
}

// ========================= Coins =========================
type Coins struct {
	Coins []sdk.Coin `json:"coins"`
}

type DecCoins struct {
	DecCoins []sdk.DecCoin `json:"dec_coins"`
}

// ========================= Delegator Rewards =========================

type DelegatorReward struct {
	DecCoins   sdk.DecCoins `json:"dec_coins"`
	ValAddress string       `json:"validator_address"`
}
