package remote

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	// "github.com/tendermint/tendermint/proto/tendermint"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"

	// bytes "github.com/cometbft/cometbft/libs/bytes"

	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	// basetypes "github.com/cosmos/cosmos-sdk/x/types"
	tmservice "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/forbole/juno/v4/node/remote"

	stakingsource "github.com/forbole/bdjuno/v4/modules/staking/source"
)

var (
	_ stakingsource.Source = &Source{}
)

// Source implements stakingsource.Source using a remote node
type Source struct {
	*remote.Source
	stakingClient stakingtypes.QueryClient
	baseClient    tmservice.ServiceClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, stakingClient stakingtypes.QueryClient, baseClient tmservice.ServiceClient) *Source {
	return &Source{
		Source:        source,
		stakingClient: stakingClient,
		baseClient:    baseClient,
	}
}

// GetValidator implements stakingsource.Source
func (s Source) GetValidator(height int64, valOper string) (stakingtypes.Validator, error) {
	res, err := s.stakingClient.Validator(
		remote.GetHeightRequestContext(s.Ctx, height),
		&stakingtypes.QueryValidatorRequest{ValidatorAddr: valOper},
	)
	if err != nil {
		return stakingtypes.Validator{}, fmt.Errorf("error while getting validator: %s", err)
	}

	return res.Validator, nil
}

// GetValidatorsWithStatus implements stakingsource.Source
func (s Source) GetValidatorsWithStatus(height int64, status string) ([]stakingtypes.Validator, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var validators []stakingtypes.Validator
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.stakingClient.Validators(
			ctx,
			&stakingtypes.QueryValidatorsRequest{
				Status: status,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 validators at time
				},
			},
		)
		if err != nil {
			return nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		validators = append(validators, res.Validators...)
	}

	return validators, nil
}

// GetPool implements stakingsource.Source
func (s Source) GetPool(height int64) (stakingtypes.Pool, error) {
	res, err := s.stakingClient.Pool(remote.GetHeightRequestContext(s.Ctx, height), &stakingtypes.QueryPoolRequest{})
	if err != nil {
		return stakingtypes.Pool{}, err
	}

	return res.Pool, nil
}

// GetParams implements stakingsource.Source
func (s Source) GetParams(height int64) (stakingtypes.Params, error) {
	res, err := s.stakingClient.Params(remote.GetHeightRequestContext(s.Ctx, height), &stakingtypes.QueryParamsRequest{})
	if err != nil {
		return stakingtypes.Params{}, err
	}

	return res.Params, nil
}

// GetTmValidator implements stakingsource.Source
func (s Source) GetTmValidator(height int64, valOper string) (stakingtypes.Validator, error) {
	var ptr ValidatorQueryResult
	resp, err := http.Get(fmt.Sprintf("https://rpc-testnet.neutron.forbole.com/validators"))
	if err != nil {
		return stakingtypes.Validator{}, err
	}

	defer resp.Body.Close()

	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return stakingtypes.Validator{}, fmt.Errorf("error while reading response body: %s", err)
	}

	err = json.Unmarshal(bz, &ptr)
	if err != nil {
		return stakingtypes.Validator{}, fmt.Errorf("error while unmarshaling response body: %s", err)
	}
	// pp := crypto.PublicKey{}

	pp = crypto.PublicKey{}
	pubkey, err := cryptocodec.FromTmProtoPublicKey(ptr.Result.Validators[0].PubKey.Value)
	if err != nil {
		// An error here would indicate that the validator updates
		// received from the provider are invalid.
		panic(err)
	}
	addr := pubkey.Address()
	fmt.Printf("\n\n addr %v \n\n", addr)

	return stakingtypes.Validator{}, nil
}
