package remote

import (
	"fmt"
	"strconv"

	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/forbole/bdjuno/v2/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// GetBalances implements bankkeeper.Source
func (s Source) GetBalances(addresses []string, height int64) ([]types.AccountBalance, error) {
	// header := remote.GetHeightRequestHeader(height)

	var header metadata.MD
	var balances []types.AccountBalance
	for _, address := range addresses {
		balRes, err := s.bankClient.AllBalances(
			metadata.AppendToOutgoingContext(s.Ctx, grpctypes.GRPCBlockHeightHeader, strconv.Itoa(int(height))),
			&banktypes.QueryAllBalancesRequest{Address: address},
			grpc.Header(&header),
		)
		if err != nil {
			return nil, fmt.Errorf("error while getting all balances: %s", err)
		}

		blockHeight := header.Get(grpctypes.GRPCBlockHeightHeader)
		fmt.Println("-- Blockheight: ", blockHeight)

		balances = append(balances, types.NewAccountBalance(
			address,
			balRes.Balances,
			height,
		))

	}

	return balances, nil
}
