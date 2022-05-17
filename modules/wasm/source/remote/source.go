package remote

import (
	"fmt"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v3/node/remote"

	wasmsource "github.com/forbole/bdjuno/v3/modules/wasm/source"
)

var (
	_ wasmsource.Source = &Source{}
)

// Source implements wasmsource.Source using a remote node
type Source struct {
	*remote.Source
	wasmClient wasmtypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, wasmClient wasmtypes.QueryClient) *Source {
	return &Source{
		Source:     source,
		wasmClient: wasmClient,
	}
}

// GetContractInfo implements wasmsource.Source
func (s Source) GetContractInfo(height int64, contractAddr string) (*wasmtypes.QueryContractInfoResponse, error) {
	res, err := s.wasmClient.ContractInfo(
		remote.GetHeightRequestContext(s.Ctx, height),
		&wasmtypes.QueryContractInfoRequest{
			Address: contractAddr,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting contract info: %s", err)
	}

	return res, nil
}

// GetContractStates implements wasmsource.Source
func (s Source) GetContractStates(height int64, contractAddr string) ([]wasmtypes.Model, error) {

	var models []wasmtypes.Model
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.wasmClient.AllContractState(
			remote.GetHeightRequestContext(s.Ctx, height),
			&wasmtypes.QueryAllContractStateRequest{
				Address: contractAddr,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 states at time
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error while getting contract state: %s", err)
		}

		stop = len(res.Pagination.NextKey) == 0
		models = append(models, res.Models...)
	}

	return models, nil
}
