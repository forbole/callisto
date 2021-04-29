package types

// CommunityPoolRow represents a single row inside the total_supply table
type CommunityPoolRow struct {
	Coins *DbDecCoins `db:"coins"`
}

// NewCommunityPoolRow allows to easily create a new NewCommunityPoolRow
func NewCommunityPoolRow(coins DbDecCoins) CommunityPoolRow {
	return CommunityPoolRow{
		Coins: &coins,
	}
}

// Equals return true if one CommunityPoolRow representing the same row as the original one
func (v CommunityPoolRow) Equals(w CommunityPoolRow) bool {
	return v.Coins.Equal(w.Coins)
}

// --------------------------------------------------------------------------------------------------------------------

type ValidatorCommissionAmountRow struct {
	ValidatorAddress string     `db:"validator_address"`
	Amount           DbDecCoins `db:"amount"`
}

func NewValidatorCommissionAmountRow(validatorAddress string, amount DbDecCoins) ValidatorCommissionAmountRow {
	return ValidatorCommissionAmountRow{
		ValidatorAddress: validatorAddress,
		Amount:           amount,
	}
}

// Equals return true iff v and w contain the same data
func (v ValidatorCommissionAmountRow) Equals(w ValidatorCommissionAmountRow) bool {
	return v.ValidatorAddress == w.ValidatorAddress &&
		v.Amount.Equal(&w.Amount)
}

// --------------------------------------------------------------------------------------------------------------------

// DelegatorRewardRow represents a single row inside the "delegation_reward" table
type DelegatorRewardRow struct {
	ValidatorConsAddr string     `db:"validator_address"`
	DelegatorAddr     string     `db:"delegator_address"`
	WithdrawAddr      string     `db:"withdraw_address"`
	Amount            DbDecCoins `db:"amount"`
}

func NewDelegatorRewardRow(validatorConsAdd, delegatorAddr, withdrawAddr string, amount DbDecCoins) DelegatorRewardRow {
	return DelegatorRewardRow{
		ValidatorConsAddr: validatorConsAdd,
		DelegatorAddr:     delegatorAddr,
		WithdrawAddr:      withdrawAddr,
		Amount:            amount,
	}
}

// Equals return true iff v and w contain the same data
func (v DelegatorRewardRow) Equals(w DelegatorRewardRow) bool {
	return v.ValidatorConsAddr == w.ValidatorConsAddr &&
		v.DelegatorAddr == w.DelegatorAddr &&
		v.WithdrawAddr == w.WithdrawAddr &&
		v.Amount.Equal(&w.Amount)
}
