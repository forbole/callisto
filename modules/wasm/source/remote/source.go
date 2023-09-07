package remote

import (
	"fmt"
	"strings"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v5/node/remote"

	wasmsource "github.com/forbole/bdjuno/v4/modules/wasm/source"
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

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		models = append(models, res.Models...)
	}

	return models, nil
}

// GetCodes implements wasmsource.Source
func (s Source) GetCodes(height int64) ([]wasmtypes.CodeInfoResponse, error) {

	var codes []wasmtypes.CodeInfoResponse
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.wasmClient.Codes(
			remote.GetHeightRequestContext(s.Ctx, height),
			&wasmtypes.QueryCodesRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 states at time
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error while getting contract codes: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		codes = append(codes, res.CodeInfos...)
	}

	return codes, nil
}

// GetContractsByCode implements wasmsource.Source
func (s Source) GetContractsByCode(height int64, codeID uint64) ([]string, error) {
	var contracts []string
	res, err := s.wasmClient.ContractsByCode(
		remote.GetHeightRequestContext(s.Ctx, height),
		&wasmtypes.QueryContractsByCodeRequest{
			CodeId: codeID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting contracts by code info: %s", err)
	}

	for _, c := range res.Contracts {
		v := strings.Split(c, ",") // Split the values
		contracts = append(contracts, v...)
	}
	return contracts, nil
}
