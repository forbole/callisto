package types

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/forbole/juno/v3/node/remote"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v3/node/local"
	markettypes "github.com/ovrclk/akash/x/market/types/v1beta2"

	nodeconfig "github.com/forbole/juno/v3/node/config"

	banksource "github.com/forbole/bdjuno/v3/modules/bank/source"
	localbanksource "github.com/forbole/bdjuno/v3/modules/bank/source/local"
	remotebanksource "github.com/forbole/bdjuno/v3/modules/bank/source/remote"
	distrsource "github.com/forbole/bdjuno/v3/modules/distribution/source"
	localdistrsource "github.com/forbole/bdjuno/v3/modules/distribution/source/local"
	remotedistrsource "github.com/forbole/bdjuno/v3/modules/distribution/source/remote"
	govsource "github.com/forbole/bdjuno/v3/modules/gov/source"
	localgovsource "github.com/forbole/bdjuno/v3/modules/gov/source/local"
	remotegovsource "github.com/forbole/bdjuno/v3/modules/gov/source/remote"
	marketsource "github.com/forbole/bdjuno/v3/modules/market/source"
	localmarketsource "github.com/forbole/bdjuno/v3/modules/market/source/local"
	remotemarketsource "github.com/forbole/bdjuno/v3/modules/market/source/remote"
	mintsource "github.com/forbole/bdjuno/v3/modules/mint/source"
	localmintsource "github.com/forbole/bdjuno/v3/modules/mint/source/local"
	remotemintsource "github.com/forbole/bdjuno/v3/modules/mint/source/remote"
	providersource "github.com/forbole/bdjuno/v3/modules/provider/source"
	localprovidersource "github.com/forbole/bdjuno/v3/modules/provider/source/local"
	remoteprovidersource "github.com/forbole/bdjuno/v3/modules/provider/source/remote"
	slashingsource "github.com/forbole/bdjuno/v3/modules/slashing/source"
	localslashingsource "github.com/forbole/bdjuno/v3/modules/slashing/source/local"
	remoteslashingsource "github.com/forbole/bdjuno/v3/modules/slashing/source/remote"
	stakingsource "github.com/forbole/bdjuno/v3/modules/staking/source"
	localstakingsource "github.com/forbole/bdjuno/v3/modules/staking/source/local"
	remotestakingsource "github.com/forbole/bdjuno/v3/modules/staking/source/remote"
	providertypes "github.com/ovrclk/akash/x/provider/types/v1beta2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	escrowKeeper "github.com/ovrclk/akash/x/escrow/keeper"
	akashmarket "github.com/ovrclk/akash/x/market"
	akashprovider "github.com/ovrclk/akash/x/provider"

	sdkclient "github.com/cosmos/cosmos-sdk/client"
	akashclient "github.com/ovrclk/akash/client"
	httpclient "github.com/tendermint/tendermint/rpc/client/http"

	akashapp "github.com/ovrclk/akash/app"
)

type Sources struct {
	BankSource     banksource.Source
	DistrSource    distrsource.Source
	GovSource      govsource.Source
	MarketSource   marketsource.Source
	MintSource     mintsource.Source
	ProviderSource providersource.Source
	SlashingSource slashingsource.Source
	StakingSource  stakingsource.Source
}

func BuildSources(nodeCfg nodeconfig.Config, encodingConfig *params.EncodingConfig) (*Sources, error) {
	switch cfg := nodeCfg.Details.(type) {
	case *remote.Details:
		return buildRemoteSources(cfg)
	case *local.Details:
		return buildLocalSources(cfg, encodingConfig)

	default:
		return nil, fmt.Errorf("invalid configuration type: %T", cfg)
	}
}

func buildLocalSources(cfg *local.Details, encodingConfig *params.EncodingConfig) (*Sources, error) {
	source, err := local.NewSource(cfg.Home, encodingConfig)
	if err != nil {
		return nil, err
	}

	app := simapp.NewSimApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, true, map[int64]bool{},
		cfg.Home, 0, simapp.MakeTestEncodingConfig(), simapp.EmptyAppOptions{},
	)

	// For MarketSource & ProviderSource
	akashapp := akashapp.NewApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, true, 0, map[int64]bool{},
		cfg.Home, simapp.EmptyAppOptions{},
	)
	escrowKeeper := escrowKeeper.NewKeeper(encodingConfig.Marshaler, sdk.NewKVStoreKey(akashprovider.StoreKey), app.BankKeeper)

	sources := &Sources{
		BankSource:     localbanksource.NewSource(source, banktypes.QueryServer(app.BankKeeper)),
		DistrSource:    localdistrsource.NewSource(source, distrtypes.QueryServer(app.DistrKeeper)),
		GovSource:      localgovsource.NewSource(source, govtypes.QueryServer(app.GovKeeper)),
		MarketSource:   localmarketsource.NewSource(source, akashmarket.NewKeeper(encodingConfig.Marshaler, sdk.NewKVStoreKey(akashprovider.StoreKey), akashapp.GetSubspace(markettypes.ModuleName), escrowKeeper).NewQuerier()),
		MintSource:     localmintsource.NewSource(source, minttypes.QueryServer(app.MintKeeper)),
		ProviderSource: localprovidersource.NewSource(source, akashprovider.NewKeeper(encodingConfig.Marshaler, sdk.NewKVStoreKey(akashprovider.StoreKey)).NewQuerier()),
		SlashingSource: localslashingsource.NewSource(source, slashingtypes.QueryServer(app.SlashingKeeper)),
		StakingSource:  localstakingsource.NewSource(source, stakingkeeper.Querier{Keeper: app.StakingKeeper}),
	}

	// Mount and initialize the stores
	err = source.MountKVStores(app, "keys")
	if err != nil {
		return nil, err
	}

	err = source.MountTransientStores(app, "tkeys")
	if err != nil {
		return nil, err
	}

	err = source.MountMemoryStores(app, "memKeys")
	if err != nil {
		return nil, err
	}

	err = source.InitStores()
	if err != nil {
		return nil, err
	}

	return sources, nil
}

func buildRemoteSources(cfg *remote.Details) (*Sources, error) {
	source, err := remote.NewSource(cfg.GRPC)
	if err != nil {
		return nil, fmt.Errorf("error while creating remote source: %s", err)
	}

	akashclient, err := buildAkashClient(cfg.RPC.Address)
	if err != nil {
		return nil, fmt.Errorf("error while building akash client: %s", err)
	}

	return &Sources{
		BankSource:     remotebanksource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		DistrSource:    remotedistrsource.NewSource(source, distrtypes.NewQueryClient(source.GrpcConn)),
		GovSource:      remotegovsource.NewSource(source, govtypes.NewQueryClient(source.GrpcConn)),
		MarketSource:   remotemarketsource.NewSource(source, markettypes.NewQueryClient(source.GrpcConn)),
		MintSource:     remotemintsource.NewSource(source, minttypes.NewQueryClient(source.GrpcConn)),
		ProviderSource: remoteprovidersource.NewSource(source, providertypes.NewQueryClient(source.GrpcConn), akashclient),
		SlashingSource: remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		StakingSource:  remotestakingsource.NewSource(source, stakingtypes.NewQueryClient(source.GrpcConn)),
	}, nil
}

func buildAkashClient(rpcAddress string) (akashclient.QueryClient, error) {
	// Build rpc client
	rpcClient, err := httpclient.New(rpcAddress, "/websocket")
	if err != nil {
		return nil, fmt.Errorf("error while building rpc client: %s", err)
	}

	err = rpcClient.Start()
	if err != nil {
		return nil, fmt.Errorf("error while starting rpc client: %s", err)
	}

	// Build akashClient
	akashclient := akashclient.NewQueryClientFromCtx(sdkclient.Context{
		Client:  rpcClient,
		NodeURI: rpcAddress,
	})

	return akashclient, nil
}
