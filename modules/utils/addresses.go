package utils

import sdk "github.com/cosmos/cosmos-sdk/types"

// FilterNonAccountAddresses filters all the non-account addresses from the given slice of addresses, returning a new
// slice containing only account addresses.
func FilterNonAccountAddresses(addresses []string) []string {
	// Filter using only the account addresses as the MessageAddressesParser might return also validator addresses
	var accountAddresses []string
	for _, address := range addresses {
		_, err := sdk.AccAddressFromBech32(address)
		if err == nil {
			accountAddresses = append(accountAddresses, address)
		}
	}
	return accountAddresses
}
