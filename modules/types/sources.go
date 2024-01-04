package types

import (
	"fmt"
	"os"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/forbole/juno/v5/types/params"

	minttypes "github.com/Stride-Labs/stride/v16/x/mint/types"
	providertypes "github.com/cosmos/interchain-security/v3/x/ccv/provider/types"
	providersource "github.com/forbole/bdjuno/v4/modules/ccv/provider/source"
	remoteprovidersource "github.com/forbole/bdjuno/v4/modules/ccv/provider/source/remote"
	"github.com/forbole/juno/v5/node/remote"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v5/node/local"

	strideapp "github.com/Stride-Labs/stride/v16/app"
	mintkeeper "github.com/Stride-Labs/stride/v16/x/mint/keeper"
	stakeibctypes "github.com/Stride-Labs/stride/v16/x/stakeibc/types"
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
	stakeibcsource "github.com/forbole/bdjuno/v4/modules/stakeibc/source"
	localstakeibcsource "github.com/forbole/bdjuno/v4/modules/stakeibc/source/local"
	remotestakeibcsource "github.com/forbole/bdjuno/v4/modules/stakeibc/source/remote"
	stakingsource "github.com/forbole/bdjuno/v4/modules/staking/source"
	localstakingsource "github.com/forbole/bdjuno/v4/modules/staking/source/local"
	remotestakingsource "github.com/forbole/bdjuno/v4/modules/staking/source/remote"
	nodeconfig "github.com/forbole/juno/v5/node/config"
)

type Sources struct {
	BankSource     banksource.Source
	DistrSource    distrsource.Source
	GovSource      govsource.Source
	MintSource     mintsource.Source
	SlashingSource slashingsource.Source
	StakeIBCSource stakeibcsource.Source
	StakingSource  stakingsource.Source
	ProviderSource providersource.Source
}

func BuildSources(nodeCfg nodeconfig.Config, encodingConfig params.EncodingConfig) (*Sources, error) {
	switch cfg := nodeCfg.Details.(type) {
	case *remote.Details:
		return buildRemoteSources(cfg)
	case *local.Details:
		return buildLocalSources(cfg, encodingConfig)

	default:
		return nil, fmt.Errorf("invalid configuration type: %T", cfg)
	}
}

func buildLocalSources(cfg *local.Details, encodingConfig params.EncodingConfig) (*Sources, error) {
	source, err := local.NewSource(cfg.Home, encodingConfig)
	if err != nil {
		return nil, err
	}

	app := strideapp.NewStrideApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, true, map[int64]bool{},
		cfg.Home, 0, strideapp.MakeEncodingConfig(), nil,
	)

	sources := &Sources{
		BankSource:     localbanksource.NewSource(source, banktypes.QueryServer(app.BankKeeper)),
		DistrSource:    localdistrsource.NewSource(source, distributionkeeper.Querier{Keeper: app.DistrKeeper}),
		GovSource:      localgovsource.NewSource(source, govtypesv1.QueryServer(app.GovKeeper)),
		MintSource:     localmintsource.NewSource(source, mintkeeper.Querier{Keeper: app.MintKeeper}),
		SlashingSource: localslashingsource.NewSource(source, slashingtypes.QueryServer(app.SlashingKeeper)),
		StakeIBCSource: localstakeibcsource.NewSource(source, stakeibctypes.QueryServer(app.StakeibcKeeper)),
		StakingSource:  localstakingsource.NewSource(source, stakingkeeper.Querier{Keeper: &app.StakingKeeper}),
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
		DistrSource:    remotedistrsource.NewSource(source, distrtypes.NewQueryClient(source.GrpcConn)),
		GovSource:      remotegovsource.NewSource(source, govtypesv1.NewQueryClient(source.GrpcConn)),
		MintSource:     remotemintsource.NewSource(source, minttypes.NewQueryClient(source.GrpcConn)),
		SlashingSource: remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		StakeIBCSource: remotestakeibcsource.NewSource(source, stakeibctypes.NewQueryClient(source.GrpcConn)),
		StakingSource:  remotestakingsource.NewSource(source, stakingtypes.NewQueryClient(source.GrpcConn)),
		ProviderSource: remoteprovidersource.NewSource(providerSource, providertypes.NewQueryClient(providerSource.GrpcConn)),
	}, nil
}
