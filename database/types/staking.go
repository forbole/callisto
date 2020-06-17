package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// ValidatorInfoRow represents a single row of the validator_info table
type ValidatorInfoRow struct {
	ConsAddress string `db:"consensus_address"`
	ValAddress  string `db:"operator_address"`
	ConsPubKey  string `db:"consensus_pubkey"`
}

// Validator data implements bstaking.Validator interface
type ValidatorData struct {
	ConsAddress sdk.ConsAddress
	ValAddress  sdk.ValAddress
	ConsPubKey  crypto.PubKey
}

func (v ValidatorData) GetConsAddr() sdk.ConsAddress {
	return v.ConsAddress
}

func (v ValidatorData) GetConsPubKey() crypto.PubKey {
	return v.ConsPubKey
}

func (v ValidatorData) GetOperator() sdk.ValAddress {
	return v.ValAddress
}
