package remote

import (
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v5/node/remote"

	wormholesource "github.com/forbole/bdjuno/v4/modules/wormhole/source"
	wormholetypes "github.com/wormhole-foundation/wormchain/x/wormhole/types"
)

var (
	_ wormholesource.Source = &Source{}
)

// Source implements wormholesource.Source using a remote node
type Source struct {
	*remote.Source
	querier wormholetypes.QueryClient
}

// NewSource returns a new Source implementation
func NewSource(source *remote.Source, querier wormholetypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetGuardianSetAll implements wormholesource.Source
func (s Source) GetGuardianSetAll(height int64) ([]wormholetypes.GuardianSet, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var guardianSet []wormholetypes.GuardianSet
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.querier.GuardianSetAll(
			ctx,
			&wormholetypes.QueryAllGuardianSetRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 guardians set at once
				},
			},
		)
		if err != nil {
			return nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		guardianSet = append(guardianSet, res.GuardianSet...)
	}

	return guardianSet, nil
}

// GetGuardianValidatorAll implements wormholesource.Source
func (s Source) GetGuardianValidatorAll(height int64) ([]wormholetypes.GuardianValidator, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var guardianValidatorList []wormholetypes.GuardianValidator
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.querier.GuardianValidatorAll(
			ctx,
			&wormholetypes.QueryAllGuardianValidatorRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 guardians validators at once
				},
			},
		)
		if err != nil {
			return nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		guardianValidatorList = append(guardianValidatorList, res.GuardianValidator...)
	}

	return guardianValidatorList, nil
}

// GetAllowlistAll implements wormholesource.Source
func (s Source) GetAllowlistAll(height int64) ([]*wormholetypes.ValidatorAllowedAddress, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var validatorAllowedAddress []*wormholetypes.ValidatorAllowedAddress
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.querier.AllowlistAll(
			ctx,
			&wormholetypes.QueryAllValidatorAllowlist{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 at once
				},
			},
		)
		if err != nil {
			return nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		validatorAllowedAddress = append(validatorAllowedAddress, res.Allowlist...)
	}

	return validatorAllowedAddress, nil
}
