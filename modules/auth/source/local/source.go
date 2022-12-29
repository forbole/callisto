package local

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/forbole/juno/v3/node/local"

	source "github.com/forbole/bdjuno/v3/modules/auth/source"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the bank keeper that works on a local node
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

	return nil, nil
}
