package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/tendermint/tendermint/crypto"
)

// Validator represents a single validator.
// This is defined as an interface so that we can use the SDK types
// as well as database types properly.
type Validator interface {
	GetConsAddr() sdk.ConsAddress
	GetConsPubKey() crypto.PubKey
	GetOperator() sdk.ValAddress
}

// ValidatorUptime contains the uptime information of a single
// validator for a specific height and point in time
type ValidatorUptime struct {
	ValidatorAddress    sdk.ConsAddress
	Height              int64
	SignedBlocksWindow  int64
	MissedBlocksCounter int64
}

// ValidatorDelegations contains both a validator delegations as
// well as its unbonding delegations
type ValidatorDelegations struct {
	ConsAddress          sdk.ConsAddress
	Delegations          staking.DelegationResponses
	UnbondingDelegations staking.UnbondingDelegations
	Height               int64
	Timestamp            time.Time
}
