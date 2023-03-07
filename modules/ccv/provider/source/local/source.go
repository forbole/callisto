package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	ccvprovidertypes "github.com/cosmos/interchain-security/x/ccv/provider/types"
	"github.com/forbole/juno/v4/node/local"

	ccvprovidersource "github.com/forbole/bdjuno/v4/modules/ccv/provider/source"
)

var (
	_ ccvprovidersource.Source = &Source{}
)

// Source implements ccvprovidersource.Source using a local node
type Source struct {
	*local.Source
	querier ccvprovidertypes.QueryServer
}

// NewSource implements a new Source instance
func NewSource(source *local.Source, querier ccvprovidertypes.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

