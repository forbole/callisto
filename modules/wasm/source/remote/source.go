package remote

import (
	"fmt"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v5/node/remote"

	wasmsource "github.com/forbole/callisto/v4/modules/wasm/source"
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

// GetCodesInfos implements wasmsource.Source
func (s Source) GetCodesInfos(height int64) ([]wasmtypes.CodeInfoResponse, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var codeInfosRes []wasmtypes.CodeInfoResponse
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.wasmClient.Codes(
			ctx,
			&wasmtypes.QueryCodesRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 codes at time
				},
			},
		)
		if err != nil {
			return nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		codeInfosRes = append(codeInfosRes, res.CodeInfos...)
	}

	return codeInfosRes, nil
}

// GetCodeBinary implements wasmsource.Source
func (s Source) GetCodeBinary(codeID uint64, height int64) ([]byte, error) {
	res, err := s.wasmClient.Code(
		remote.GetHeightRequestContext(s.Ctx, height),
		&wasmtypes.QueryCodeRequest{
			CodeId: codeID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting contract code binary: %s", err)
	}

	return res.Data, nil
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
