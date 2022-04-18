package utils_test

import (
	"testing"

	"github.com/forbole/bdjuno/v3/modules/utils"

	"github.com/stretchr/testify/require"
)

func TestFilterNonAccountAddresses(t *testing.T) {
	addresses := []string{
		"cosmos1hafptm4zxy5nw8rd2pxyg83c5ls2v62tstzuv2",
		"cosmosvaloper1hafptm4zxy5nw8rd2pxyg83c5ls2v62t4lkfqe",
	}

	filtered := utils.FilterNonAccountAddresses(addresses)
	require.Equal(t, []string{
		"cosmos1hafptm4zxy5nw8rd2pxyg83c5ls2v62tstzuv2",
	}, filtered)
}
