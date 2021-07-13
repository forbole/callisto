package types

import( 
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Account represents a chain account
type Account struct {
	Address string
	Details authtypes.AccountI
}

// NewAccount builds a new Account instance
func NewAccount(address string, accountI authtypes.AccountI) Account {
	return Account{
		Address: address,
		Details: accountI,
	}
}
