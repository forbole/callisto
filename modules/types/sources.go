package types

import (
	"fmt"
	"os"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/forbole/juno/v5/node/remote"
	params "github.com/forbole/juno/v5/types/params"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	providertypes "github.com/cosmos/interchain-security/v3/x/ccv/provider/types"
	remotewasmsource "github.com/forbole/bdjuno/v4/modules/wasm/source/remote"
	nodeconfig "github.com/forbole/juno/v5/node/config"
	"github.com/forbole/juno/v5/node/local"

	neutronapp "github.com/MonikaCat/neutron/v2/app"
	banksource "github.com/forbole/bdjuno/v4/modules/bank/source"
	remotebanksource "github.com/forbole/bdjuno/v4/modules/bank/source/remote"
	providersource "github.com/forbole/bdjuno/v4/modules/ccv/provider/source"
	remoteprovidersource "github.com/forbole/bdjuno/v4/modules/ccv/provider/source/remote"
	slashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source"
	remoteslashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source/remote"
	wasmsource "github.com/forbole/bdjuno/v4/modules/wasm/source"

	localbanksource "github.com/forbole/bdjuno/v4/modules/bank/source/local"
	localslashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source/local"
)

type Sources struct {
	BankSource     banksource.Source
	ProviderSource providersource.Source
	SlashingSource slashingsource.Source
	WasmSource     wasmsource.Source
}

func BuildSources(nodeCfg nodeconfig.Config, encodingConfig params.EncodingConfig) (*Sources, error) {
	switch cfg := nodeCfg.Details.(type) {
	case *remote.Details:
		return buildRemoteSources(cfg)
	case *local.Details:
		return nil, fmt.Errorf("local source is not supported: %T", cfg)
	default:
		return nil, fmt.Errorf("invalid configuration type: %T", cfg)
	}
}

func buildLocalSources(cfg *local.Details, encodingConfig params.EncodingConfig) (*Sources, error) {
	source, err := local.NewSource(cfg.Home, encodingConfig)
	if err != nil {
		return nil, err
	}

	app := neutronapp.New(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, true, map[int64]bool{},
		cfg.Home, 0, neutronapp.MakeEncodingConfig(), nil, nil, nil,
	)

	sources := &Sources{
		BankSource:     localbanksource.NewSource(source, banktypes.QueryServer(app.BankKeeper)),
		SlashingSource: localslashingsource.NewSource(source, slashingtypes.QueryServer(app.SlashingKeeper)),
		// WasmSource:     localwasmsource.NewSource(source, wasmtypes.QueryServer(app.WasmKeeper)),
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

	providerSource, err := remote.NewSource(cfg.ProviderGRPC)
	if err != nil {
		return nil, fmt.Errorf("error while creating remote provider source: %s", err)
	}

	return &Sources{
		BankSource:     remotebanksource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		ProviderSource: remoteprovidersource.NewSource(providerSource, providertypes.NewQueryClient(providerSource.GrpcConn)),
		SlashingSource: remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		WasmSource:     remotewasmsource.NewSource(source, wasmtypes.NewQueryClient(source.GrpcConn)),
	}, nil
}
