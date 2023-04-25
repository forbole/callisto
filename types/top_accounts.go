package types

// TopAccount represents a cheqd account from top accounts module
type TopAccount struct {
	Address string
	Type    string
}

// TopAccount builds a new TopAccount instance
func NewTopAccount(address, accountType string) TopAccount {
	return TopAccount{
		Address: address,
		Type:    accountType,
	}
}
