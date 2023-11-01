package types

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/forbole/juno/v4/node/remote"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v4/node/local"
	mintkeeper "github.com/osmosis-labs/osmosis/v20/x/mint/keeper"
	minttypes "github.com/osmosis-labs/osmosis/v20/x/mint/types"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banksource "github.com/forbole/bdjuno/v4/modules/bank/source"
	localbanksource "github.com/forbole/bdjuno/v4/modules/bank/source/local"
	remotebanksource "github.com/forbole/bdjuno/v4/modules/bank/source/remote"
	distrsource "github.com/forbole/bdjuno/v4/modules/distribution/source"
	localdistrsource "github.com/forbole/bdjuno/v4/modules/distribution/source/local"
	remotedistrsource "github.com/forbole/bdjuno/v4/modules/distribution/source/remote"
	govsource "github.com/forbole/bdjuno/v4/modules/gov/source"
	localgovsource "github.com/forbole/bdjuno/v4/modules/gov/source/local"
	remotegovsource "github.com/forbole/bdjuno/v4/modules/gov/source/remote"
	mintsource "github.com/forbole/bdjuno/v4/modules/mint/source"
	localmintsource "github.com/forbole/bdjuno/v4/modules/mint/source/local"
	remotemintsource "github.com/forbole/bdjuno/v4/modules/mint/source/remote"
	slashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source"
	localslashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source/local"
	remoteslashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source/remote"
	stakingsource "github.com/forbole/bdjuno/v4/modules/staking/source"
	localstakingsource "github.com/forbole/bdjuno/v4/modules/staking/source/local"
	remotestakingsource "github.com/forbole/bdjuno/v4/modules/staking/source/remote"
	superfluidsource "github.com/forbole/bdjuno/v4/modules/superfluid/source"
	localsuperfluidsource "github.com/forbole/bdjuno/v4/modules/superfluid/source/local"
	remotesuperfluidsource "github.com/forbole/bdjuno/v4/modules/superfluid/source/remote"
	wasmsource "github.com/forbole/bdjuno/v4/modules/wasm/source"
	localwasmsource "github.com/forbole/bdjuno/v4/modules/wasm/source/local"
	remotewasmsource "github.com/forbole/bdjuno/v4/modules/wasm/source/remote"
	nodeconfig "github.com/forbole/juno/v4/node/config"
	osmosisapp "github.com/osmosis-labs/osmosis/v20/app"
	superfluidkeeper "github.com/osmosis-labs/osmosis/v20/x/superfluid/keeper"
	superfluidtypes "github.com/osmosis-labs/osmosis/v20/x/superfluid/types"
)

type Sources struct {
	BankSource       banksource.Source
	DistrSource      distrsource.Source
	GovSource        govsource.Source
	MintSource       mintsource.Source
	SlashingSource   slashingsource.Source
	StakingSource    stakingsource.Source
	SuperfluidSource superfluidsource.Source
	WasmSource       wasmsource.Source
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

	app := osmosisapp.NewOsmosisApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, true, map[int64]bool{},
		cfg.Home, 0, simapp.EmptyAppOptions{}, osmosisapp.EmptyWasmOpts,
	)

	sources := &Sources{
		BankSource:       localbanksource.NewSource(source, bankkeeper.Querier{BaseKeeper: *app.BankKeeper}),
		DistrSource:      localdistrsource.NewSource(source, distrtypes.QueryServer(app.DistrKeeper)),
		GovSource:        localgovsource.NewSource(source, govtypes.QueryServer(app.GovKeeper)),
		MintSource:       localmintsource.NewSource(source, mintkeeper.Querier{Keeper: *app.MintKeeper}),
		SlashingSource:   localslashingsource.NewSource(source, slashingtypes.QueryServer(app.SlashingKeeper)),
		StakingSource:    localstakingsource.NewSource(source, stakingkeeper.Querier{Keeper: *app.StakingKeeper}),
		SuperfluidSource: localsuperfluidsource.NewSource(source, superfluidkeeper.Querier{Keeper: *app.SuperfluidKeeper}),
		WasmSource:       localwasmsource.NewSource(source, wasmkeeper.Querier(app.WasmKeeper)),
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

	return &Sources{
		BankSource:       remotebanksource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		DistrSource:      remotedistrsource.NewSource(source, distrtypes.NewQueryClient(source.GrpcConn)),
		GovSource:        remotegovsource.NewSource(source, govtypes.NewQueryClient(source.GrpcConn)),
		MintSource:       remotemintsource.NewSource(source, minttypes.NewQueryClient(source.GrpcConn)),
		SlashingSource:   remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		StakingSource:    remotestakingsource.NewSource(source, stakingtypes.NewQueryClient(source.GrpcConn)),
		SuperfluidSource: remotesuperfluidsource.NewSource(source, superfluidtypes.NewQueryClient(source.GrpcConn)),
		WasmSource:       remotewasmsource.NewSource(source, wasmtypes.NewQueryClient(source.GrpcConn)),
	}, nil
}
