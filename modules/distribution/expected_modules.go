package distribution

// BankModule represents the bank module we expect
type BankModule interface {
	RefreshBalances(height int64, addresses []string) error
}
