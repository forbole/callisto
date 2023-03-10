package local

import (
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/juno/v4/node/local"

	sdk "github.com/cosmos/cosmos-sdk/types"
	source "github.com/forbole/bdjuno/v4/modules/auth/source"
)

var (
	_ source.Source = &Source{}
)

// Source implements authsource.Source by using a local node
type Source struct {
	*local.Source
	q authtypes.QueryServer
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, q authtypes.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      q,
	}
}

func (s Source) GetAllAnyAccounts(height int64) ([]*codectypes.Any, error) {
	log.Debug().Msg("getting all accounts")
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var accounts []*codectypes.Any
	var nextKey []byte
	var stop = false
	var counter uint64
	var totalCounts uint64

	// Get 1000 accounts per query
	var pageLimit uint64 = 1000

	for !stop {
		// Get accounts
		res, err := s.q.Accounts(
			sdk.WrapSDKContext(ctx),
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
