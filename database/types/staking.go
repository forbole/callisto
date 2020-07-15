package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/tendermint/tendermint/crypto"
)

// StakingPoolRow represents a single row inside the staking_pool table
type StakingPoolRow struct {
	BondedTokens    int64     `db:"bonded_tokens"`
	NotBondedTokens int64     `db:"not_bonded_tokens"`
	Height          int64     `db:"height"`
	Timestamp       time.Time `db:"timestamp"`
}

// NewStakingPoolRow allows to easily create a new StakingPoolRow
func NewStakingPoolRow(bondedTokens, notBondedTokens int64, height int64, timestamp time.Time) StakingPoolRow {
	return StakingPoolRow{
		BondedTokens:    bondedTokens,
		NotBondedTokens: notBondedTokens,
		Height:          height,
		Timestamp:       timestamp,
	}
}

// Equal allows to tells whether r and as represent the same rows
func (r StakingPoolRow) Equal(s StakingPoolRow) bool {
	return r.BondedTokens == s.BondedTokens &&
		r.NotBondedTokens == s.NotBondedTokens &&
		r.Height == s.Height &&
		r.Timestamp.Equal(s.Timestamp)
}

// ________________________________________________

// ValidatorRow represents a single row of the validator table
type ValidatorRow struct {
	ConsAddress string `db:"consensus_address"`
	ConsPubKey  string `db:"consensus_pubkey"`
}

// NewValidatorRow returns a new ValidatorRow
func NewValidatorRow(consAddress, consPubKey string) ValidatorRow {
	return ValidatorRow{
		ConsAddress: consAddress,
		ConsPubKey:  consPubKey,
	}
}

// Equal tells whether v and w contain the same data
func (v ValidatorRow) Equal(w ValidatorRow) bool {
	return v.ConsAddress == w.ConsAddress &&
		v.ConsPubKey == w.ConsPubKey
}

// ________________________________________________

// ValidatorInfoRow represents a single row of the validator_info table
type ValidatorInfoRow struct {
	ConsAddress         string `db:"consensus_address"`
	ValAddress          string `db:"operator_address"`
	Moniker             string `db:"moniker"`
	Identity            string `db:"identity"`
	Website             string `db:"website"`
	SecurityContact     string `db:"security_contact"`
	Details             string `db:"details"`
	SelfDelegateAddress string `db:"self_delegate_address"`
}

// NewValidatorInfoRow allows to build a new ValidatorInfoRow
func NewValidatorInfoRow(
	consAddress string, valAddress string, selfDelegateAddress string, moniker string, identity string,
	website string, securityContact string, details string,
) ValidatorInfoRow {
	return ValidatorInfoRow{
		ConsAddress:         consAddress,
		ValAddress:          valAddress,
		Moniker:             moniker,
		Identity:            identity,
		Website:             website,
		SecurityContact:     securityContact,
		Details:             details,
		SelfDelegateAddress: selfDelegateAddress,
	}
}

// Equal tells whether v and w represent the same rows
func (v ValidatorInfoRow) Equal(w ValidatorInfoRow) bool {
	return v.ConsAddress == w.ConsAddress &&
		v.ValAddress == w.ValAddress &&
		v.Moniker == w.Moniker &&
		v.Identity == w.Identity &&
		v.Website == w.Website &&
		v.SecurityContact == w.SecurityContact &&
		v.Details == w.Details &&
		v.SelfDelegateAddress == w.SelfDelegateAddress
}

// ________________________________________________

// ValidatorData contains all the data of a single validator.
// It implements types.Validator interface
type ValidatorData struct {
	ConsAddress         string `db:"consensus_address"`
	ValAddress          string `db:"operator_address"`
	ConsPubKey          string `db:"consensus_pubkey"`
	SelfDelegateAddress string `db:"self_delegate_address"`
	Moniker             string `db:"moniker"`
	Identity            string `db:"identity"`
	Website             string `db:"website"`
	SecurityContact     string `db:"security_contact"`
	Details             string `db:"details"`
}

// NewValidatorData allows to build a new ValidatorData
func NewValidatorData(
	consAddress, valAddress, consPubKey string, selfDelegateAddress string, moniker string, identity string,
	website string, securityContact string, details string,
) ValidatorData {
	return ValidatorData{
		ConsAddress:         consAddress,
		ValAddress:          valAddress,
		ConsPubKey:          consPubKey,
		SelfDelegateAddress: selfDelegateAddress,
		Moniker:             moniker,
		Identity:            identity,
		Website:             website,
		SecurityContact:     securityContact,
		Details:             details,
	}
}

func (v ValidatorData) GetConsAddr() sdk.ConsAddress {
	addr, err := sdk.ConsAddressFromBech32(v.ConsAddress)
	if err != nil {
		panic(err)
	}

	return addr
}

func (v ValidatorData) GetConsPubKey() crypto.PubKey {
	return sdk.MustGetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, v.ConsPubKey)
}

func (v ValidatorData) GetOperator() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(v.ValAddress)
	if err != nil {
		panic(err)
	}

	return addr
}

func (v ValidatorData) GetDescription() staking.Description {
	return staking.NewDescription(
		v.Moniker,
		v.Identity,
		v.Website,
		v.SecurityContact,
		v.Details,
	)
}

func (v ValidatorData) GetSelfDelegateAddress() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(v.SelfDelegateAddress)
	if err != nil {
		panic(err)
	}

	return addr
}

// ________________________________________________

// ValidatorUptimeRow represents a single row of the validator_uptime table
type ValidatorUptimeRow struct {
	ConsAddr           string `db:"validator_address"`
	Height             int64  `db:"height"`
	SignedBlockWindow  int64  `db:"signed_blocks_window"`
	MissedBlockCounter int64  `db:"missed_blocks_counter"`
}

// NewValidatorUptimeRow allows to build a new ValidatorUptimeRow
func NewValidatorUptimeRow(consAddr string, signedBlocWindow, missedBlocksCounter, height int64) ValidatorUptimeRow {
	return ValidatorUptimeRow{
		ConsAddr:           consAddr,
		SignedBlockWindow:  signedBlocWindow,
		MissedBlockCounter: missedBlocksCounter,
		Height:             height,
	}
}

// Equal tells whether v and w contain the same data
func (v ValidatorUptimeRow) Equal(w ValidatorUptimeRow) bool {
	return v.ConsAddr == w.ConsAddr &&
		v.Height == w.Height &&
		v.SignedBlockWindow == w.SignedBlockWindow &&
		v.MissedBlockCounter == w.MissedBlockCounter
}

// ________________________________________________

// ValidatorDelegationRow represents a single validator_delegation table row
type ValidatorDelegationRow struct {
	ConsensusAddress string    `db:"consensus_address"`
	DelegatorAddress string    `db:"delegator_address"`
	Amount           DbCoin    `db:"amount"`
	Height           int64     `db:"height"`
	Timestamp        time.Time `db:"timestamp"`
}

// NewValidatorDelegationRow allows to build a new ValidatorDelegationRow
func NewValidatorDelegationRow(
	consAddr, delegator string, amount DbCoin,
	height int64, timestamp time.Time,
) ValidatorDelegationRow {
	return ValidatorDelegationRow{
		ConsensusAddress: consAddr,
		DelegatorAddress: delegator,
		Amount:           amount,
		Height:           height,
		Timestamp:        timestamp,
	}
}

// Equals tells whether v and w represent the same row
func (v ValidatorDelegationRow) Equal(w ValidatorDelegationRow) bool {
	return v.ConsensusAddress == w.ConsensusAddress &&
		v.DelegatorAddress == w.DelegatorAddress &&
		v.Amount.Equal(w.Amount) &&
		v.Height == w.Height &&
		v.Timestamp.Equal(w.Timestamp)
}

// ________________________________________________

// ValidatorUnbondingDelegationRow represents a single row inside the
// validator_unbonding_delegation table
type ValidatorUnbondingDelegationRow struct {
	ConsensusAddress    string    `db:"consensus_address"`
	DelegatorAddress    string    `db:"delegator_address"`
	Amount              DbCoin    `db:"amount"`
	CompletionTimestamp time.Time `db:"completion_timestamp"`
	Height              int64     `db:"height"`
	Timestamp           time.Time `db:"timestamp"`
}

// NewValidatorUnbondingDelegationRow allows to build a new
// ValidatorUnbondingDelegationRow instance
func NewValidatorUnbondingDelegationRow(
	consAddr, delegator string, amount DbCoin, completionTimestamp time.Time,
	height int64, timestamp time.Time,
) ValidatorUnbondingDelegationRow {
	return ValidatorUnbondingDelegationRow{
		ConsensusAddress:    consAddr,
		DelegatorAddress:    delegator,
		Amount:              amount,
		CompletionTimestamp: completionTimestamp,
		Height:              height,
		Timestamp:           timestamp,
	}
}

// Equal tells whether v and w represent the same rows
func (v ValidatorUnbondingDelegationRow) Equal(w ValidatorUnbondingDelegationRow) bool {
	return v.ConsensusAddress == w.ConsensusAddress &&
		v.DelegatorAddress == w.DelegatorAddress &&
		v.Amount.Equal(w.Amount) &&
		v.CompletionTimestamp.Equal(w.CompletionTimestamp) &&
		v.Height == w.Height &&
		v.Timestamp.Equal(w.Timestamp)
}

// ________________________________________________

// ValidatorReDelegationRow represents a single row of the
// validator_redelegation database table
type ValidatorReDelegationRow struct {
	DelegatorAddress    string    `db:"delegator_address"`
	SrcValidatorAddress string    `db:"src_validator_address"`
	DstValidatorAddress string    `db:"dst_validator_address"`
	Amount              DbCoin    `db:"amount"`
	Height              int64     `db:"height"`
	CompletionTime      time.Time `db:"completion_time"`
}

// NewValidatorReDelegationRow allows to easily build a new
// ValidatorReDelegationRow instance
func NewValidatorReDelegationRow(
	delegator, srcConsAddr, dstConsAddr string, amount DbCoin,
	height int64, completionTime time.Time,
) ValidatorReDelegationRow {
	return ValidatorReDelegationRow{
		DelegatorAddress:    delegator,
		SrcValidatorAddress: srcConsAddr,
		DstValidatorAddress: dstConsAddr,
		Amount:              amount,
		Height:              height,
		CompletionTime:      completionTime,
	}
}

// Equal tells whether v and w represent the same database rows
func (v ValidatorReDelegationRow) Equal(w ValidatorReDelegationRow) bool {
	return v.DelegatorAddress == w.DelegatorAddress &&
		v.SrcValidatorAddress == w.SrcValidatorAddress &&
		v.DstValidatorAddress == w.DstValidatorAddress &&
		v.Amount.Equal(w.Amount) &&
		v.Height == w.Height &&
		v.CompletionTime.Equal(w.CompletionTime)
}

// ValidatorCommission represents a single row of the
// validator_commission database table
type ValidatorCommission struct {
	ValidatorAddress  string    `db:"validator_address"`
	Timestamp         time.Time `db:"timestamp"`
	Commission        int64     `db:"commission"`
	MinSelfDelegation int64     `db:"min_self_delegation"`
	Height            int64     `db:"height"`
}

// NewValidatorCommission allows to easily build a new
// ValidatorCommission instance
func NewValidatorCommission(
	validatorAddress string, commission int64, minSelfDelegation int64, height int64, timestamp time.Time,
) ValidatorCommission {
	return ValidatorCommission{
		ValidatorAddress:  validatorAddress,
		Timestamp:         timestamp,
		Commission:        commission,
		MinSelfDelegation: minSelfDelegation,
		Height:            height,
	}
}

// Equal tells whether v and w represent the same rows
func (v ValidatorCommission) Equal(w ValidatorCommission) bool {
	return v.ValidatorAddress == w.ValidatorAddress &&
		v.Timestamp.Equal(w.Timestamp) &&
		v.Commission == w.Commission &&
		v.MinSelfDelegation == w.MinSelfDelegation &&
		v.Height == w.Height
}

//ValidatorDelegation store the return of validator_delegation_shares
type ValidatorDelegationSharesRow struct {
	OperatorAddress  string    `db:"operator_address"`
	DelegatorAddress string    `db:"delegator_address"`
	Shares           int64     `db:"shares"`
	Timestamp        time.Time `db:"timestamp"`
	Height           int64     `db:"height"`
}

//Equal determain two validatorDelegation refer as same row
func (v ValidatorDelegationSharesRow) Equal(w ValidatorDelegationSharesRow) bool {
	return v.OperatorAddress == w.OperatorAddress &&
		v.DelegatorAddress == w.DelegatorAddress &&
		v.Shares == w.Shares &&
		v.Timestamp.Equal(w.Timestamp) &&
		v.Height == w.Height
}

// NewValidatorDelegationSharesRow make a new instance of ValidatorDelegationSharesRow
func NewValidatorDelegationSharesRow(
	operatorAddress string, delegatorAddress string, shares int64,
	timestamp time.Time, height int64,
) ValidatorDelegationSharesRow {
	return ValidatorDelegationSharesRow{
		OperatorAddress:  operatorAddress,
		DelegatorAddress: delegatorAddress,
		Shares:           shares,
		Timestamp:        timestamp,
		Height:           height,
	}
}
