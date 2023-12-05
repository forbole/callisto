package remote

import (
	"github.com/forbole/juno/v5/node/remote"

	"github.com/forbole/bdjuno/v4/modules/overgold/chain/stake/source"

	staketypes "git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the QueryClient that works on a remote node
type Source struct {
	*remote.Source
	client staketypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, stakeClient staketypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: stakeClient,
	}
}

// GetStakes implements Source
func (s Source) GetStakes(address []string, height int64) ([]*staketypes.Stake, error) {
	return []*staketypes.Stake{}, nil
}
