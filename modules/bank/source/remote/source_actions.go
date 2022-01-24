package remote

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"google.golang.org/grpc/metadata"
)

// GetAccountBalances implements bankkeeper.Source
func (s Source) GetAccountBalance(address string, height int64) ([]sdk.Coin, error) {

	// Get account balance at certain height
	balRes, err := s.bankClient.AllBalances(
		metadata.AppendToOutgoingContext(
			s.Ctx,
			grpctypes.GRPCBlockHeightHeader,
			strconv.Itoa(int(height)),
		),
		&banktypes.QueryAllBalancesRequest{Address: address},
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting all balances: %s", err)
	}

	fmt.Println("Balance response: ", balRes)

	return balRes.Balances, nil
}
