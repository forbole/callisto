package staking

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestAddress(t *testing.T) {
	prefix := "desmos"

	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(
		prefix,
		prefix+sdk.PrefixPublic,
	)
	cfg.SetBech32PrefixForValidator(
		prefix+sdk.PrefixValidator+sdk.PrefixOperator,
		prefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic,
	)
	cfg.SetBech32PrefixForConsensusNode(
		prefix+sdk.PrefixValidator+sdk.PrefixConsensus,
		prefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic,
	)

	consAddr, err := sdk.ConsAddressFromHex("0549E17DCCE2AA4BAB0640C70A692389FE081351")
	require.NoError(t, err)

	fmt.Printf(consAddr.String())
}
