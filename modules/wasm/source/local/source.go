package local

import (
	"fmt"
	"strings"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v5/node/local"

	wasmsource "github.com/forbole/bdjuno/v4/modules/wasm/source"
)

var (
	_ wasmsource.Source = &Source{}
)

// Source implements wasmsource.Source using a local node
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

// GetContractStates implements wasmsource.Source
func (s Source) GetContractStates(height int64, contractAddr string) ([]wasmtypes.Model, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var models []wasmtypes.Model
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.q.AllContractState(
			sdk.WrapSDKContext(ctx),
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
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var codes []wasmtypes.CodeInfoResponse
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.q.Codes(
			sdk.WrapSDKContext(ctx),
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
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var contracts []string
	res, err := s.q.ContractsByCode(
		sdk.WrapSDKContext(ctx),
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
