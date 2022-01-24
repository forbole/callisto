package remote

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/forbole/juno/v2/node/remote"

	bankkeeper "github.com/forbole/bdjuno/v2/modules/bank/source"
)

var (
	_ bankkeeper.Source = &Source{}
)

type Source struct {
	*remote.Source
	bankClient banktypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, bankClient banktypes.QueryClient) *Source {
	return &Source{
		Source:     source,
		bankClient: bankClient,
	}
}

// GetSupply implements bankkeeper.Source
func (s Source) GetSupply(height int64) (sdk.Coins, error) {
	header := remote.GetHeightRequestHeader(height)
	res, err := s.bankClient.TotalSupply(s.Ctx, &banktypes.QueryTotalSupplyRequest{}, header)
	if err != nil {
		return nil, fmt.Errorf("error while getting total supply: %s", err)
	}

	return res.Supply, nil
}
