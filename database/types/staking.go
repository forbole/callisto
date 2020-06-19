package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// ValidatorInfoRow represents a single row of the validator_info table.
// It implements types.Validator interface
type ValidatorInfoRow struct {
	ConsAddress string `db:"consensus_address"`
	ValAddress  string `db:"operator_address"`
	ConsPubKey  string `db:"consensus_pubkey"`
}

func (v ValidatorInfoRow) GetConsAddr() sdk.ConsAddress {
	addr, err := sdk.ConsAddressFromBech32(v.ConsAddress)
	if err != nil {
		panic(err)
	}

	return addr
}

func (v ValidatorInfoRow) GetConsPubKey() crypto.PubKey {
	return sdk.MustGetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, v.ValAddress)
}

func (v ValidatorInfoRow) GetOperator() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(v.ValAddress)
	if err != nil {
		panic(err)
	}

	return addr
}
