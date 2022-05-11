package utils

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

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

// ConvertAddressPrefix converts the bech32 address to the desired prefix
func ConvertAddressPrefix(prefix string, bech32Add string) (string, error) {
	_, bz, err := bech32.DecodeAndConvert(bech32Add)
	if err != nil {
		return "", fmt.Errorf("error while decoding bech32 address(%s): %s", bech32Add, err)
	}

	newBech32Add, err := bech32.ConvertAndEncode(prefix, bz)
	if err != nil {
		return "", fmt.Errorf("error while encoding bech32 address(%s) with %s prefix: %s", bech32Add, prefix, err)
	}

	return newBech32Add, nil
}
