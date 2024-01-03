package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v5/node/local"
	"github.com/forbole/juno/v5/node/remote"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	markertypes "github.com/MonCatCat/provenance/x/marker/types"
	banksource "github.com/forbole/bdjuno/v4/modules/bank/source"
	remotebanksource "github.com/forbole/bdjuno/v4/modules/bank/source/remote"
	distrsource "github.com/forbole/bdjuno/v4/modules/distribution/source"
	remotedistrsource "github.com/forbole/bdjuno/v4/modules/distribution/source/remote"
	govsource "github.com/forbole/bdjuno/v4/modules/gov/source"
	remotegovsource "github.com/forbole/bdjuno/v4/modules/gov/source/remote"
	markersourcer "github.com/forbole/bdjuno/v4/modules/marker/source"
	remotemarkersource "github.com/forbole/bdjuno/v4/modules/marker/source/remote"
	mintsource "github.com/forbole/bdjuno/v4/modules/mint/source"
	remotemintsource "github.com/forbole/bdjuno/v4/modules/mint/source/remote"
	slashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source"
	remoteslashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source/remote"
	stakingsource "github.com/forbole/bdjuno/v4/modules/staking/source"
	remotestakingsource "github.com/forbole/bdjuno/v4/modules/staking/source/remote"
	wasmsource "github.com/forbole/bdjuno/v4/modules/wasm/source"
	remotewasmsource "github.com/forbole/bdjuno/v4/modules/wasm/source/remote"
	nodeconfig "github.com/forbole/juno/v5/node/config"
)

type Sources struct {
	BankSource     banksource.Source
	DistrSource    distrsource.Source
	GovSource      govsource.Source
	MarkerSource   markersourcer.Source
	MintSource     mintsource.Source
	SlashingSource slashingsource.Source
	StakingSource  stakingsource.Source
	WasmSource     wasmsource.Source
}

func BuildSources(nodeCfg nodeconfig.Config, encodingConfig *params.EncodingConfig) (*Sources, error) {
	switch cfg := nodeCfg.Details.(type) {
	case *remote.Details:
		return buildRemoteSources(cfg)
	case *local.Details:
		return nil, fmt.Errorf("local mode is currently not supported")

	default:
		return nil, fmt.Errorf("invalid configuration type: %T", cfg)
	}
}

// func buildLocalSources(cfg *local.Details, encodingConfig *params.EncodingConfig) (*Sources, error) {
// 	source, err := local.NewSource(cfg.Home, encodingConfig)
// 	if err != nil {
// 		return nil, err
// 	}

// 	app := provenanceapp.New(
// 		log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, true, map[int64]bool{},
// 		cfg.Home, 0, provenanceapp.MakeEncodingConfig(), simapp.EmptyAppOptions{},
// 	)

// 	sources := &Sources{
// 		BankSource:     localbanksource.NewSource(source, banktypes.QueryServer(app.BankKeeper)),
// 		DistrSource:    localdistrsource.NewSource(source, distrtypes.QueryServer(app.DistrKeeper)),
// 		GovSource:      localgovsource.NewSource(source, govtypesv1.QueryServer(app.GovKeeper), nil),
// 		MarkerSource:   localmarkersource.NewSource(source, markertypes.QueryServer(app.MarkerKeeper)),
// 		MintSource:     localmintsource.NewSource(source, minttypes.QueryServer(app.MintKeeper)),
// 		SlashingSource: localslashingsource.NewSource(source, slashingtypes.QueryServer(app.SlashingKeeper)),
// 		StakingSource:  localstakingsource.NewSource(source, stakingkeeper.Querier{Keeper: app.StakingKeeper}),
// 		WasmSource:     localwasmsource.NewSource(source, wasmkeeper.Querier(&app.WasmKeeper)),
// 	}

// 	// Mount and initialize the stores
// 	err = source.MountKVStores(app, "keys")
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = source.MountTransientStores(app, "tkeys")
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = source.MountMemoryStores(app, "memKeys")
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = source.InitStores()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return sources, nil
// }

func buildRemoteSources(cfg *remote.Details) (*Sources, error) {
	source, err := remote.NewSource(cfg.GRPC)
	if err != nil {
		return nil, fmt.Errorf("error while creating remote source: %s", err)
	}

	return &Sources{
		BankSource:     remotebanksource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		DistrSource:    remotedistrsource.NewSource(source, distrtypes.NewQueryClient(source.GrpcConn)),
		GovSource:      remotegovsource.NewSource(source, govtypesv1.NewQueryClient(source.GrpcConn), govtypesv1beta1.NewQueryClient(source.GrpcConn)),
		MarkerSource:   remotemarkersource.NewSource(source, markertypes.NewQueryClient(source.GrpcConn)),
		MintSource:     remotemintsource.NewSource(source, minttypes.NewQueryClient(source.GrpcConn)),
		SlashingSource: remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		StakingSource:  remotestakingsource.NewSource(source, stakingtypes.NewQueryClient(source.GrpcConn)),
		WasmSource:     remotewasmsource.NewSource(source, wasmtypes.NewQueryClient(source.GrpcConn)),
	}, nil
}
