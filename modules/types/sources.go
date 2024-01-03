package types

import (
	"fmt"
	"os"

	bundlestypes "github.com/KYVENetwork/chain/x/bundles/types"
	querytypes "github.com/KYVENetwork/chain/x/query/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cometbft/cometbft/libs/log"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v5/node/local"
	"github.com/forbole/juno/v5/node/remote"
	params "github.com/forbole/juno/v5/types/params"

	nodeconfig "github.com/forbole/juno/v5/node/config"

	banksource "github.com/forbole/bdjuno/v4/modules/bank/source"
	bundlessource "github.com/forbole/bdjuno/v4/modules/bundles/source"
	globalsource "github.com/forbole/bdjuno/v4/modules/global/source"
	stakerssource "github.com/forbole/bdjuno/v4/modules/stakers/source"

	kyveapp "github.com/KYVENetwork/chain/app"
	remotebanksource "github.com/forbole/bdjuno/v4/modules/bank/source/remote"
	remotebundlessource "github.com/forbole/bdjuno/v4/modules/bundles/source/remote"

	localbanksource "github.com/forbole/bdjuno/v4/modules/bank/source/local"

	distrsource "github.com/forbole/bdjuno/v4/modules/distribution/source"
	remotedistrsource "github.com/forbole/bdjuno/v4/modules/distribution/source/remote"
	govsource "github.com/forbole/bdjuno/v4/modules/gov/source"
	localgovsource "github.com/forbole/bdjuno/v4/modules/gov/source/local"
	remotegovsource "github.com/forbole/bdjuno/v4/modules/gov/source/remote"
	mintsource "github.com/forbole/bdjuno/v4/modules/mint/source"
	localmintsource "github.com/forbole/bdjuno/v4/modules/mint/source/local"
	remotemintsource "github.com/forbole/bdjuno/v4/modules/mint/source/remote"
	poolsource "github.com/forbole/bdjuno/v4/modules/pool/source"
	remotepoolsource "github.com/forbole/bdjuno/v4/modules/pool/source/remote"
	slashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source"
	localslashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source/local"
	remoteslashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source/remote"
	remotestakerssource "github.com/forbole/bdjuno/v4/modules/stakers/source/remote"
	stakingsource "github.com/forbole/bdjuno/v4/modules/staking/source"
	localstakingsource "github.com/forbole/bdjuno/v4/modules/staking/source/local"
	remotestakingsource "github.com/forbole/bdjuno/v4/modules/staking/source/remote"
)

type Sources struct {
	BankSource     banksource.Source
	BundlesSource  bundlessource.Source
	DistrSource    distrsource.Source
	GlobalSource   globalsource.Source
	GovSource      govsource.Source
	MintSource     mintsource.Source
	PoolSource     poolsource.Source
	SlashingSource slashingsource.Source
	StakingSource  stakingsource.Source
	StakersSource  stakerssource.Source
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

	app := kyveapp.NewKYVEApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, true, nil, nil,
	)

	sources := &Sources{
		BankSource: localbanksource.NewSource(source, banktypes.QueryServer(app.BankKeeper)),
		// DistrSource:    localdistrsource.NewSource(source, distrtypes.QueryServer(app.DistributionKeeper)),
		GovSource:      localgovsource.NewSource(source, govtypesv1.QueryServer(app.GovKeeper)),
		MintSource:     localmintsource.NewSource(source, minttypes.QueryServer(app.MintKeeper)),
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

	return &Sources{
		BankSource:     remotebanksource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		BundlesSource:  remotebundlessource.NewSource(source, bundlestypes.NewQueryClient(source.GrpcConn)),
		DistrSource:    remotedistrsource.NewSource(source, distrtypes.NewQueryClient(source.GrpcConn)),
		GovSource:      remotegovsource.NewSource(source, govtypesv1.NewQueryClient(source.GrpcConn)),
		MintSource:     remotemintsource.NewSource(source, minttypes.NewQueryClient(source.GrpcConn)),
		PoolSource:     remotepoolsource.NewSource(source, querytypes.NewQueryPoolClient(source.GrpcConn)),
		SlashingSource: remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		StakingSource:  remotestakingsource.NewSource(source, stakingtypes.NewQueryClient(source.GrpcConn)),
		StakersSource:  remotestakerssource.NewSource(source, stakerstypes.NewQueryClient(source.GrpcConn), querytypes.NewQueryStakersClient(source.GrpcConn)),
	}, nil
}
