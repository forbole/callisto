package local

import (
	staketypes "git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	"github.com/forbole/juno/v5/node/local"

	"github.com/forbole/bdjuno/v4/modules/overgold/chain/stake/source"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the QueryServer that works on a local node
type Source struct {
	*local.Source
	stakeServer staketypes.QueryServer
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, stakeServer staketypes.QueryServer) *Source {
	return &Source{
		Source:      source,
		stakeServer: stakeServer,
	}
}

// GetStakes implements Source
func (s Source) GetStakes(address []string, height int64) ([]*staketypes.Stake, error) {
	return nil, nil
}
