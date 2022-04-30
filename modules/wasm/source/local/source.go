package local

import (
	"fmt"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v3/node/local"

	wasmsource "github.com/forbole/bdjuno/v3/modules/wasm/source"
)

var (
	_ wasmsource.Source = &Source{}
)

// Source implements stakingsource.Source using a local node
type Source struct {
	*local.Source
	q wasmtypes.QueryServer
}

// NewSource returns a new Source instance
func NewSource(source *local.Source, querier wasmtypes.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      querier,
	}
}

// GetContractInfo implements wasmsource.Source
func (s Source) GetContractInfo(height int64, contractAddr string) (*wasmtypes.QueryContractInfoResponse, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.q.ContractInfo(
		sdk.WrapSDKContext(ctx),
		&wasmtypes.QueryContractInfoRequest{
			Address: contractAddr,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting contract info: %s", err)
	}

	return res, nil
}
