package remote

import (
	"github.com/forbole/juno/v5/node/remote"

	"github.com/forbole/bdjuno/v4/modules/overgold/chain/referral/source"

	referraltypes "git.ooo.ua/vipcoin/ovg-chain/x/referral/types"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the QueryClient that works on a remote node
type Source struct {
	*remote.Source
	client referraltypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, referralClient referraltypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: referralClient,
	}
}

// GetStats implements Source
func (s Source) GetStats(dates []string, height int64) ([]*referraltypes.Stats, error) {
	return []*referraltypes.Stats{}, nil
}
