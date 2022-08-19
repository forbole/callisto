package remote

import (
	"fmt"

	providersource "github.com/forbole/bdjuno/v3/modules/provider/source"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v3/node/remote"

	"github.com/ovrclk/akash/provider"
	gwrest "github.com/ovrclk/akash/provider/gateway/rest"
	providertypes "github.com/ovrclk/akash/x/provider/types/v1beta2"

	sdkclient "github.com/cosmos/cosmos-sdk/client"
	akashclient "github.com/ovrclk/akash/client"
	httpclient "github.com/tendermint/tendermint/rpc/client/http"
)

var (
	_ providersource.Source = &Source{}
)

// Source implements providersource.Source using a remote node
type Source struct {
	*remote.Source
	providerClient providertypes.QueryClient
	rpcClient      *httpclient.HTTP
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, providerClient providertypes.QueryClient, rpcClient *httpclient.HTTP) *Source {
	return &Source{
		Source:         source,
		providerClient: providerClient,
		rpcClient:      rpcClient,
	}
}

// GetProvider implements providersource.Source
func (s Source) GetProvider(height int64, ownerAddress string) (providertypes.Provider, error) {
	res, err := s.providerClient.Provider(remote.GetHeightRequestContext(s.Ctx, height),
		&providertypes.QueryProviderRequest{
			Owner: ownerAddress,
		},
	)
	if err != nil {
		return providertypes.Provider{}, fmt.Errorf("error while getting provider: %s", err)
	}

	return res.Provider, nil
}

// GetProviders implements providersource.Source
func (s Source) GetProviders(height int64) ([]providertypes.Provider, error) {
	var providers []providertypes.Provider
	var nextKey []byte
	var stop bool
	for !stop {
		res, err := s.providerClient.Providers(remote.GetHeightRequestContext(s.Ctx, height),
			&providertypes.QueryProvidersRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 1000, // Query 1000 providers at a time
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error while getting providers: %s", err)
		}
		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		providers = append(providers, res.Providers...)
	}

	return providers, nil
}

func (s Source) GetProviderInventoryStatus(address string) (*provider.Status, error) {
	s.rpcClient.Start()

	// Get sdk address
	bech32Addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, fmt.Errorf("error while converting address to AccAddress: %s", err)
	}

	// Build akashClient
	akashclient := akashclient.NewQueryClientFromCtx(sdkclient.Context{
		Client:  s.rpcClient,
		NodeURI: "https://rpc.akash.forbole.com:443",
	})

	// Builde gateway rest client
	gclient, err := gwrest.NewClient(akashclient, bech32Addr, nil)
	if err != nil {
		return nil, fmt.Errorf("error while building rest client: %s", err)
	}

	res, err := gclient.Status(s.Ctx)
	if err != nil {
		return nil, fmt.Errorf("error while getting provider inventory status with rest client: %s", err)
	}

	return res, nil
}
