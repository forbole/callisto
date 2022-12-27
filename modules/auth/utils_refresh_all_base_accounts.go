package auth

import (
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/forbole/bdjuno/v3/types"
)

func (m *Module) GetAllBaseAccounts(height int64) ([]types.Account, error) {
	anyAccounts, err := m.source.GetAllAnyAccounts(height)
	if err != nil {
		return nil, fmt.Errorf("error while getting any accounts: %s", err)
	}
	unpacked, err := m.unpackAnyAccounts(anyAccounts)
	if err != nil {
		return nil, err
	}

	return unpacked, nil

}

func (m *Module) unpackAnyAccounts(anyAccounts []*codectypes.Any) ([]types.Account, error) {
	accounts := []types.Account{}
	for _, account := range anyAccounts {
		var accountI authtypes.AccountI
		err := m.cdc.UnpackAny(account, &accountI)
		if err != nil {
			return nil, fmt.Errorf("error while unpacking any account: %s", err)
		}

		// Get only base accounts
		if baseAccount, ok := accountI.(*authtypes.BaseAccount); ok {
			accounts = append(accounts, types.NewAccount(baseAccount.Address))
		}
	}

	return accounts, nil

}
