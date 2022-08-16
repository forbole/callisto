package remote

import (
	"fmt"

	providersource "github.com/forbole/bdjuno/v3/modules/provider/source"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v3/node/remote"

	sdkclient "github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"

	akashclient "github.com/ovrclk/akash/client"
	"github.com/ovrclk/akash/provider"
	gwrest "github.com/ovrclk/akash/provider/gateway/rest"
	providertypes "github.com/ovrclk/akash/x/provider/types/v1beta2"
)

var (
	_ providersource.Source = &Source{}
)

// Source implements providersource.Source using a remote node
type Source struct {
	*remote.Source
	providerClient            providertypes.QueryClient
	providerGatewatRestClient gwrest.Client
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, providerClient providertypes.QueryClient) *Source {
	return &Source{
		Source:         source,
		providerClient: providerClient,
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

func (s Source) ProviderStatus(address string, height int64) (*provider.Status, error) {
	bech32Addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, fmt.Errorf("error while converting address to AccAddress: %s", err)
	}

	cctx, err := sdkclient.GetClientTxContext(nil)
	if err != nil {
		return nil, fmt.Errorf("error while getting client tx context: %s", err)
	}

	gclient, err := gwrest.NewClient(akashclient.NewQueryClientFromCtx(cctx), bech32Addr, nil)
	if err != nil {
		return nil, fmt.Errorf("error while building new akash gateway rest client: %s", err)
	}

	res, err := gclient.Status(s.Ctx)
	if err != nil {
		return nil, fmt.Errorf("error while getting provider status with gateway rest client: %s", err)
	}

	return res, nil
}

// akash provider status akash1wxr49evm8hddnx9ujsdtd86gk46s7ejnccqfmy
// {
//   "cluster": {
//     "leases": 3,
//     "inventory": {
//       "active": [
//         {
//           "cpu": 8000,
//           "memory": 8589934592,
//           "storage_ephemeral": 5384815247360
//         },
//         {
//           "cpu": 100000,
//           "memory": 450971566080,
//           "storage_ephemeral": 982473768960
//         },
//         {
//           "cpu": 8000,
//           "memory": 8589934592,
//           "storage_ephemeral": 2000000000000
//         }
//       ],
//       "available": {
//         "nodes": [
//           {
//             "cpu": 111495,
//             "memory": 466163988480,
//             "storage_ephemeral": 2375935850345
//           },
//           {
//             "cpu": 118780,
//             "memory": 474497601536,
//             "storage_ephemeral": 7760751097705
//           },
//           {
//             "cpu": 110800,
//             "memory": 465918152704,
//             "storage_ephemeral": 5760751097705
//           },
//           {
//             "cpu": 19525,
//             "memory": 23846356992,
//             "storage_ephemeral": 6778277328745
//           }
//         ]
//       }
//     }
//   },
//   "bidengine": {
//     "orders": 0
//   },
//   "manifest": {
//     "deployments": 0
//   },
//   "cluster_public_hostname": "provider.bigtractorplotting.com"
// }
