package types

import (
	"fmt"
	"os"

	"cosmossdk.io/simapp"
	"cosmossdk.io/simapp/params"
	allowedtypes "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
	coretypes "git.ooo.ua/vipcoin/ovg-chain/x/core/types"
	feeexcludertypes "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	referraltypes "git.ooo.ua/vipcoin/ovg-chain/x/referral/types"
	staketypes "git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	"github.com/cometbft/cometbft/libs/log"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	nodeconfig "github.com/forbole/juno/v5/node/config"
	"github.com/forbole/juno/v5/node/local"
	"github.com/forbole/juno/v5/node/remote"

	banksource "github.com/forbole/bdjuno/v4/modules/bank/source"
	localbanksource "github.com/forbole/bdjuno/v4/modules/bank/source/local"
	remotebanksource "github.com/forbole/bdjuno/v4/modules/bank/source/remote"
	distrsource "github.com/forbole/bdjuno/v4/modules/distribution/source"
	remotedistrsource "github.com/forbole/bdjuno/v4/modules/distribution/source/remote"
	govsource "github.com/forbole/bdjuno/v4/modules/gov/source"
	localgovsource "github.com/forbole/bdjuno/v4/modules/gov/source/local"
	remotegovsource "github.com/forbole/bdjuno/v4/modules/gov/source/remote"
	mintsource "github.com/forbole/bdjuno/v4/modules/mint/source"
	localmintsource "github.com/forbole/bdjuno/v4/modules/mint/source/local"
	remotemintsource "github.com/forbole/bdjuno/v4/modules/mint/source/remote"
	overgoldAllowedSource "github.com/forbole/bdjuno/v4/modules/overgold/chain/allowed/source"
	remoteOvergoldAllowedSource "github.com/forbole/bdjuno/v4/modules/overgold/chain/allowed/source/remote"
	overgoldBankSource "github.com/forbole/bdjuno/v4/modules/overgold/chain/bank/source"
	remoteOvergoldBankSource "github.com/forbole/bdjuno/v4/modules/overgold/chain/bank/source/remote"
	overgoldCoreSource "github.com/forbole/bdjuno/v4/modules/overgold/chain/core/source"
	remoteOvergoldCoreSource "github.com/forbole/bdjuno/v4/modules/overgold/chain/core/source/remote"
	overgoldFeeExcluderSource "github.com/forbole/bdjuno/v4/modules/overgold/chain/feeexcluder/source"
	remoteOvergoldFeeExcluderSource "github.com/forbole/bdjuno/v4/modules/overgold/chain/feeexcluder/source/remote"
	overgoldReferralSource "github.com/forbole/bdjuno/v4/modules/overgold/chain/referral/source"
	remoteOvergoldReferralSource "github.com/forbole/bdjuno/v4/modules/overgold/chain/referral/source/remote"
	overgoldStakeSource "github.com/forbole/bdjuno/v4/modules/overgold/chain/stake/source"
	remoteOvergoldStakeSource "github.com/forbole/bdjuno/v4/modules/overgold/chain/stake/source/remote"
	slashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source"
	localslashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source/local"
	remoteslashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source/remote"
	stakingsource "github.com/forbole/bdjuno/v4/modules/staking/source"
	localstakingsource "github.com/forbole/bdjuno/v4/modules/staking/source/local"
	remotestakingsource "github.com/forbole/bdjuno/v4/modules/staking/source/remote"
)

type Sources struct {
	BankSource     banksource.Source
	DistrSource    distrsource.Source
	GovSource      govsource.Source
	MintSource     mintsource.Source
	SlashingSource slashingsource.Source
	StakingSource  stakingsource.Source

	// Custom OVG sources
	OverGoldAllowedSource     overgoldAllowedSource.Source
	OverGoldCoreSource        overgoldCoreSource.Source
	OverGoldFeeExcluderSource overgoldFeeExcluderSource.Source
	OverGoldReferralSource    overgoldReferralSource.Source
	OverGoldStakeSource       overgoldStakeSource.Source

	// Custom SDK sources
	OverGoldBankSource overgoldBankSource.Source
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
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, true, nil, nil,
	)

	sources := &Sources{
		BankSource: localbanksource.NewSource(source, banktypes.QueryServer(app.BankKeeper)),
		// DistrSource:    localdistrsource.NewSource(source, distrtypes.QueryServer(app.DistrKeeper)),
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
		DistrSource:    remotedistrsource.NewSource(source, distrtypes.NewQueryClient(source.GrpcConn)),
		GovSource:      remotegovsource.NewSource(source, govtypesv1.NewQueryClient(source.GrpcConn)),
		MintSource:     remotemintsource.NewSource(source, minttypes.NewQueryClient(source.GrpcConn)),
		SlashingSource: remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		StakingSource:  remotestakingsource.NewSource(source, stakingtypes.NewQueryClient(source.GrpcConn)),

		// Custom OVG sources
		OverGoldAllowedSource:     remoteOvergoldAllowedSource.NewSource(source, allowedtypes.NewQueryClient(source.GrpcConn)),
		OverGoldBankSource:        remoteOvergoldBankSource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		OverGoldCoreSource:        remoteOvergoldCoreSource.NewSource(source, coretypes.NewQueryClient(source.GrpcConn)),
		OverGoldFeeExcluderSource: remoteOvergoldFeeExcluderSource.NewSource(source, feeexcludertypes.NewQueryClient(source.GrpcConn)),
		OverGoldReferralSource:    remoteOvergoldReferralSource.NewSource(source, referraltypes.NewQueryClient(source.GrpcConn)),
		OverGoldStakeSource:       remoteOvergoldStakeSource.NewSource(source, staketypes.NewQueryClient(source.GrpcConn)),
	}, nil
}
