package utils_test

import (
	"testing"

	utils2 "github.com/forbole/bdjuno/modules/common/utils"

	"github.com/stretchr/testify/require"
)

func TestFilterNonAccountAddresses(t *testing.T) {
	addresses := []string{
		"cosmos1hafptm4zxy5nw8rd2pxyg83c5ls2v62tstzuv2",
		"cosmosvaloper1hafptm4zxy5nw8rd2pxyg83c5ls2v62t4lkfqe",
	}

	filtered := utils2.FilterNonAccountAddresses(addresses)
	require.Equal(t, []string{
		"cosmos1hafptm4zxy5nw8rd2pxyg83c5ls2v62tstzuv2",
	}, filtered)
}
