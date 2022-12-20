package types

type TopAccountsRow struct {
	Address      string `db:"address"`
	Available    int64  `db:"available"`
	Delegation   int64  `db:"delegation"`
	Redelegation int64  `db:"redelegation"`
	Unbonding    int64  `db:"unbonding"`
	Reward       int64  `db:"reward"`
	Sum          int64  `db:"sum"`
}

func NewTopAccountsRow(
	address string, available, delegation, redelegation, unbonding, reward, sum int64,
) TopAccountsRow {
	return TopAccountsRow{
		Address:      address,
		Available:    available,
		Delegation:   delegation,
		Redelegation: redelegation,
		Unbonding:    unbonding,
		Reward:       reward,
		Sum:          sum,
	}
}

// Equals return true if one TopAccountsRow representing the same row as the original one
func (a TopAccountsRow) Equals(b TopAccountsRow) bool {
	return a.Address == b.Address &&
		a.Available == b.Available &&
		a.Redelegation == b.Redelegation &&
		a.Unbonding == b.Unbonding &&
		a.Reward == b.Reward
}
