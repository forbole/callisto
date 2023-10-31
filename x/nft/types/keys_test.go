package types_test

import (
	"testing"

	"github.com/forbole/bdjuno/v4/x/nft/types"
	"github.com/stretchr/testify/require"
)

func TestSplitKeyDenomWithoutIBC(t *testing.T) {
	keyDenom := []byte("testdenomid/testtokenid")

	// nolint: govet
	denomID, tokenID, err := types.SplitKeyDenom(keyDenom)

	require.NoError(t, err)
	require.Equal(t, "testdenomid", denomID)
	require.Equal(t, "testtokenid", tokenID)
}

func TestSplitKeyDenomWithIBC(t *testing.T) {
	keyDenom := []byte("ibc/testdenomid/testtokenid")

	// nolint: govet
	denomID, tokenID, err := types.SplitKeyDenom(keyDenom)

	require.NoError(t, err)
	require.Equal(t, "ibc/testdenomid", denomID)
	require.Equal(t, "testtokenid", tokenID)
}
