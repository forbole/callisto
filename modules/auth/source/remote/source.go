package remote

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v4/node/remote"
	"github.com/rs/zerolog/log"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	source "github.com/forbole/bdjuno/v3/modules/auth/source"
)

var (
	_ source.Source = &Source{}
)

type Source struct {
	*remote.Source
	authClient authtypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, authClient authtypes.QueryClient) *Source {
	return &Source{
		Source:     source,
		authClient: authClient,
	}
}

func (s Source) GetAllAnyAccounts(height int64) ([]*codectypes.Any, error) {
	log.Debug().Msg("getting all accounts")
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var accounts []*codectypes.Any
	var nextKey []byte
	var stop = false
	var counter uint64
	var totalCounts uint64

	// Get 1000 accounts per query
	var pageLimit uint64 = 1000

	for !stop {
		// Get accounts
		res, err := s.authClient.Accounts(
			ctx,
			&authtypes.QueryAccountsRequest{
				Pagination: &query.PageRequest{
					Key:        nextKey,
					Limit:      pageLimit,
					CountTotal: true,
				},
			})
		if err != nil {
			return nil, fmt.Errorf("error while getting any accounts from source: %s", err)
		}
		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		accounts = append(accounts, res.Accounts...)

		// Log getting accounts progress
		if res.Pagination.GetTotal() != 0 {
			totalCounts = res.Pagination.GetTotal()
		}
		counter += uint64(len(res.Accounts))
		log.Debug().Uint64("total accounts", totalCounts).Uint64("current counter", counter).Msg("getting accounts...")
	}

	return accounts, nil
}
