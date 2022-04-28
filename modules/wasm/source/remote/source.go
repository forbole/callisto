package remote

import (
	"fmt"

	wasmdtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/forbole/juno/v3/node/remote"

	wasmsource "github.com/forbole/bdjuno/v3/modules/wasm/source"
)

var (
	_ wasmsource.Source = &Source{}
)

// Source implements stakingsource.Source using a remote node
type Source struct {
	*remote.Source
	wasmClient wasmdtypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, wasmClient wasmdtypes.QueryClient) *Source {
	return &Source{
		Source:     source,
		wasmClient: wasmClient,
	}
}

// GetContractHisotry implements wasmsource.Source
func (s Source) GetContractHisotry(height int64) (*wasmdtypes.QueryContractHistoryResponse, error) {
	res, err := s.wasmClient.ContractHistory(s.Ctx, &wasmdtypes.QueryContractHistoryRequest{})
	if err != nil {
		return nil, fmt.Errorf("error while getting contract history: %s", err)
	}

	return res, nil
}
