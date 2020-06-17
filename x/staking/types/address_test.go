package types

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestAddress(t *testing.T) {
	hexAddress := "E23B5B5C7DFE0B1AB3AAC9F4A7218315E0EE0810"

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(
		"desmos",
		"desmos"+sdk.PrefixPublic,
	)
	config.SetBech32PrefixForValidator(
		"desmos"+sdk.PrefixValidator+sdk.PrefixOperator,
		"desmos"+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic,
	)
	config.SetBech32PrefixForConsensusNode(
		"desmos"+sdk.PrefixValidator+sdk.PrefixConsensus,
		"desmos"+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic,
	)

	consAddr, err := sdk.ConsAddressFromHex(hexAddress)
	require.NoError(t, err)

	valAddr, err := sdk.ValAddressFromHex(hexAddress)
	require.NoError(t, err)

	fmt.Println(consAddr.String())
	fmt.Println(valAddr.String())
}
