package types

// Account represents a chain account
type Account struct {
	Address string
}

// NewAccount builds a new Account instance
func NewAccount(address string) Account {
	return Account{
		Address: address,
	}
}

// TopAccount represents a cheqd top account
type TopAccount struct {
	Address     string
	AccountType string
}

// TopAccount builds a new TopAccount instance
func NewTopAccount(address, accountType string) TopAccount {
	return TopAccount{
		Address:     address,
		AccountType: accountType,
	}
}
