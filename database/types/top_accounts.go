package types

type TopAccountsRow struct {
	Address    string `db:"address"`
	Type       string `db:"type"`
	Available  int64  `db:"available"`
	Delegation int64  `db:"delegation"`
	Unbonding  int64  `db:"unbonding"`
	Reward     int64  `db:"reward"`
	Sum        int64  `db:"sum"`
	Height     int64  `db:"height"`
}

func NewTopAccountsRow(
	address, accountType string, available, delegation, unbonding, reward, sum, height int64,
) TopAccountsRow {
	return TopAccountsRow{
		Address:    address,
		Type:       accountType,
		Available:  available,
		Delegation: delegation,
		Unbonding:  unbonding,
		Reward:     reward,
		Sum:        sum,
		Height:     height,
	}
}

// Equals return true if one TopAccountsRow representing the same row as the original one
func (a TopAccountsRow) Equals(b TopAccountsRow) bool {
	return a.Address == b.Address &&
		a.Type == b.Type &&
		a.Available == b.Available &&
		a.Delegation == b.Delegation &&
		a.Unbonding == b.Unbonding &&
		a.Reward == b.Reward &&
		a.Sum == b.Sum &&
		a.Height == b.Height
}
