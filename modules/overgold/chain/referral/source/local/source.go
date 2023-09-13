package local

import (
	referraltypes "git.ooo.ua/vipcoin/ovg-chain/x/referral/types"
	"github.com/forbole/juno/v5/node/local"

	"github.com/forbole/bdjuno/v4/modules/overgold/chain/referral/source"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the QueryServer that works on a local node
type Source struct {
	*local.Source
	referralServer referraltypes.QueryServer
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, referralServer referraltypes.QueryServer) *Source {
	return &Source{
		Source:         source,
		referralServer: referralServer,
	}
}

// GetStats implements Source
func (s Source) GetStats(dates []string, height int64) ([]*referraltypes.Stats, error) {
	return nil, nil
}
