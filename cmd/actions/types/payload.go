package types

// ========================= Address Payload =========================
type AddressPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            Address                `json:"input"`
}

type Address struct {
	Address string `json:"address"`
}

// ========================= Account Balance Payload =========================

type AccountBalancePayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            AccountBalanceArgs     `json:"input"`
}

type AccountBalanceArgs struct {
	Address string `json:"address"`
	Height  int64  `json:"height"`
}

// ========================= Staking Payload =========================

type StakingPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            StakingArgs            `json:"input"`
}

type StakingArgs struct {
	Address    string `json:"address"`
	Offset     uint64 `json:"offset"`
	Limit      uint64 `json:"limit"`
	CountTotal bool   `json:"count_total"`
}
