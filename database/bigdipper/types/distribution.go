package types

import "github.com/forbole/bdjuno/database/types"

// CommunityPoolRow represents a single row inside the total_supply table
type CommunityPoolRow struct {
	OneRowID bool              `db:"one_row_id"`
	Coins    *types.DbDecCoins `db:"coins"`
	Height   int64             `db:"height"`
}

// NewCommunityPoolRow allows to easily create a new CommunityPoolRow
func NewCommunityPoolRow(coins types.DbDecCoins, height int64) CommunityPoolRow {
	return CommunityPoolRow{
		OneRowID: true,
		Coins:    &coins,
		Height:   height,
	}
}

// Equals return true if one CommunityPoolRow representing the same row as the original one
func (v CommunityPoolRow) Equals(w CommunityPoolRow) bool {
	return v.Coins.Equal(w.Coins) &&
		v.Height == w.Height
}

// -------------------------------------------------------------------------------------------------------------------

// ValidatorCommissionAmountRow represents a single row of the "validator_commission_amount" table
type ValidatorCommissionAmountRow struct {
	ValidatorAddr string           `db:"validator_address"`
	Amount        types.DbDecCoins `db:"amount"`
	Height        int64            `db:"height"`
}

// NewValidatorCommissionAmountRow returns a new ValidatorCommissionAmountRow instance
func NewValidatorCommissionAmountRow(valAddr string, amount types.DbDecCoins, height int64) ValidatorCommissionAmountRow {
	return ValidatorCommissionAmountRow{
		ValidatorAddr: valAddr,
		Amount:        amount,
		Height:        height,
	}
}

// Equals returns true iff v and w contain the same data
func (v ValidatorCommissionAmountRow) Equals(w ValidatorCommissionAmountRow) bool {
	return v.ValidatorAddr == w.ValidatorAddr &&
		v.Amount.Equal(&w.Amount) &&
		v.Height == w.Height
}

// -------------------------------------------------------------------------------------------------------------------

// DelegationRewardRow represents a single row inside the "delegation_reward" table
type DelegationRewardRow struct {
	ValidatorConsAddress string           `db:"validator_address"`
	DelegatorAddress     string           `db:"delegator_address"`
	WithdrawAddress      string           `db:"withdraw_address"`
	Amount               types.DbDecCoins `db:"amount"`
	Height               int64            `db:"height"`
}

// NewDelegationRewardRow returns a new DelegationRewardRow instance
func NewDelegationRewardRow(delegatorAddr, valConsAddr, withdrawAddr string, amount types.DbDecCoins, height int64) DelegationRewardRow {
	return DelegationRewardRow{
		ValidatorConsAddress: valConsAddr,
		DelegatorAddress:     delegatorAddr,
		WithdrawAddress:      withdrawAddr,
		Amount:               amount,
		Height:               height,
	}
}

// Equals returns true iff v and w contain the same data
func (v DelegationRewardRow) Equals(w DelegationRewardRow) bool {
	return v.ValidatorConsAddress == w.ValidatorConsAddress &&
		v.DelegatorAddress == w.DelegatorAddress &&
		v.WithdrawAddress == w.WithdrawAddress &&
		v.Amount.Equal(&w.Amount) &&
		v.Height == w.Height
}
