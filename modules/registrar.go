package modules

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/forbole/juno/v2/modules/messages"
	"github.com/forbole/juno/v2/modules/pruning"
	"github.com/forbole/juno/v2/modules/registrar"
	"github.com/forbole/juno/v2/modules/telemetry"
	"github.com/forbole/juno/v2/node/local"
	"github.com/forbole/juno/v2/node/remote"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/modules/auth"
	"github.com/forbole/bdjuno/v2/modules/bank"
	"github.com/forbole/bdjuno/v2/modules/consensus"
	"github.com/forbole/bdjuno/v2/modules/distribution"
	"github.com/forbole/bdjuno/v2/modules/gov"
	"github.com/forbole/bdjuno/v2/modules/mint"
	"github.com/forbole/bdjuno/v2/modules/modules"
	"github.com/forbole/bdjuno/v2/modules/pricefeed"
	"github.com/forbole/bdjuno/v2/modules/slashing"
	"github.com/forbole/bdjuno/v2/modules/staking"
	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/accounts"
	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/assets"
	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/banking"
	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/wallets"
	"github.com/forbole/bdjuno/v2/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	banksource "github.com/forbole/bdjuno/v2/modules/bank/source"
	localbanksource "github.com/forbole/bdjuno/v2/modules/bank/source/local"
	remotebanksource "github.com/forbole/bdjuno/v2/modules/bank/source/remote"
	distrsource "github.com/forbole/bdjuno/v2/modules/distribution/source"
	localdistrsource "github.com/forbole/bdjuno/v2/modules/distribution/source/local"
	remotedistrsource "github.com/forbole/bdjuno/v2/modules/distribution/source/remote"
	govsource "github.com/forbole/bdjuno/v2/modules/gov/source"
	localgovsource "github.com/forbole/bdjuno/v2/modules/gov/source/local"
	remotegovsource "github.com/forbole/bdjuno/v2/modules/gov/source/remote"
	mintsource "github.com/forbole/bdjuno/v2/modules/mint/source"
	localmintsource "github.com/forbole/bdjuno/v2/modules/mint/source/local"
	remotemintsource "github.com/forbole/bdjuno/v2/modules/mint/source/remote"
	slashingsource "github.com/forbole/bdjuno/v2/modules/slashing/source"
	localslashingsource "github.com/forbole/bdjuno/v2/modules/slashing/source/local"
	remoteslashingsource "github.com/forbole/bdjuno/v2/modules/slashing/source/remote"
	stakingsource "github.com/forbole/bdjuno/v2/modules/staking/source"
	localstakingsource "github.com/forbole/bdjuno/v2/modules/staking/source/local"
	remotestakingsource "github.com/forbole/bdjuno/v2/modules/staking/source/remote"

	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"

	jmodules "github.com/forbole/juno/v2/modules"
	nodeconfig "github.com/forbole/juno/v2/node/config"

	vipcoinaccountssource "github.com/forbole/bdjuno/v2/modules/vipcoin/chain/accounts/source"
	remotevipcoinaccountssource "github.com/forbole/bdjuno/v2/modules/vipcoin/chain/accounts/source/remote"

	vipcoinwalletssource "github.com/forbole/bdjuno/v2/modules/vipcoin/chain/wallets/source"
	remotevipcoinwalletsssource "github.com/forbole/bdjuno/v2/modules/vipcoin/chain/wallets/source/remote"

	vipcoinbankingsource "github.com/forbole/bdjuno/v2/modules/vipcoin/chain/banking/source"
	remotevipcoinbankingsource "github.com/forbole/bdjuno/v2/modules/vipcoin/chain/banking/source/remote"

	vipcoinassetssource "github.com/forbole/bdjuno/v2/modules/vipcoin/chain/assets/source"
	remotevipcoinassetssource "github.com/forbole/bdjuno/v2/modules/vipcoin/chain/assets/source/remote"
)

// UniqueAddressesParser returns a wrapper around the given parser that removes all duplicated addresses
func UniqueAddressesParser(parser messages.MessageAddressesParser) messages.MessageAddressesParser {
	return func(cdc codec.Marshaler, msg sdk.Msg) ([]string, error) {
		addresses, err := parser(cdc, msg)
		if err != nil {
			return nil, err
		}

		return utils.RemoveDuplicateValues(addresses), nil
	}
}

// --------------------------------------------------------------------------------------------------------------------

var (
	_ registrar.Registrar = &Registrar{}
)

// Registrar represents the modules.Registrar that allows to register all modules that are supported by BigDipper
type Registrar struct {
	parser messages.MessageAddressesParser
}

// NewRegistrar allows to build a new Registrar instance
func NewRegistrar(parser messages.MessageAddressesParser) *Registrar {
	return &Registrar{
		parser: UniqueAddressesParser(parser),
	}
}

// BuildModules implements modules.Registrar
func (r *Registrar) BuildModules(ctx registrar.Context) jmodules.Modules {
	cdc := ctx.EncodingConfig.Marshaler
	db := database.Cast(ctx.Database)

	sources, err := BuildSources(ctx.JunoConfig.Node, ctx.EncodingConfig)
	if err != nil {
		panic(err)
	}

	authModule := auth.NewModule(r.parser, cdc, db)
	bankModule := bank.NewModule(r.parser, sources.BankSource, cdc, db)
	consensusModule := consensus.NewModule(db)
	distrModule := distribution.NewModule(sources.DistrSource, cdc, db)
	mintModule := mint.NewModule(sources.MintSource, cdc, db)
	slashingModule := slashing.NewModule(sources.SlashingSource, cdc, db)
	stakingModule := staking.NewModule(sources.StakingSource, slashingModule, cdc, db)
	govModule := gov.NewModule(sources.GovSource, authModule, distrModule, mintModule, slashingModule, stakingModule, cdc, db)

	vipcoinAccountsModule := accounts.NewModule(sources.VipcoinAccountsSource, cdc, db)
	vipcoinWalletsModule := wallets.NewModule(r.parser, sources.VipcoinWalletsSource, cdc, db)
	vipcoinBankingModule := banking.NewModule(sources.VipcoinBankingSource, cdc, db)
	vipcoinAssetsModule := assets.NewModule(sources.VipcoinAssetsSource, cdc, db)

	return []jmodules.Module{
		messages.NewModule(r.parser, cdc, ctx.Database),
		telemetry.NewModule(ctx.JunoConfig),
		pruning.NewModule(ctx.JunoConfig, db, ctx.Logger),

		authModule,
		bankModule,
		consensusModule,
		distrModule,
		govModule,
		mintModule,
		modules.NewModule(ctx.JunoConfig.Chain, db),
		pricefeed.NewModule(ctx.JunoConfig, cdc, db),
		slashingModule,
		stakingModule,

		vipcoinAccountsModule,
		vipcoinWalletsModule,
		vipcoinBankingModule,
		vipcoinAssetsModule,
	}
}

type Sources struct {
	BankSource            banksource.Source
	DistrSource           distrsource.Source
	GovSource             govsource.Source
	MintSource            mintsource.Source
	SlashingSource        slashingsource.Source
	StakingSource         stakingsource.Source
	VipcoinAccountsSource vipcoinaccountssource.Source
	VipcoinWalletsSource  vipcoinwalletssource.Source
	VipcoinBankingSource  vipcoinbankingsource.Source
	VipcoinAssetsSource   vipcoinassetssource.Source
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

	sources := &Sources{
		BankSource:     localbanksource.NewSource(source, banktypes.QueryServer(app.BankKeeper)),
		DistrSource:    localdistrsource.NewSource(source, distrtypes.QueryServer(app.DistrKeeper)),
		GovSource:      localgovsource.NewSource(source, govtypes.QueryServer(app.GovKeeper)),
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
		BankSource:            remotebanksource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		DistrSource:           remotedistrsource.NewSource(source, distrtypes.NewQueryClient(source.GrpcConn)),
		GovSource:             remotegovsource.NewSource(source, govtypes.NewQueryClient(source.GrpcConn)),
		MintSource:            remotemintsource.NewSource(source, minttypes.NewQueryClient(source.GrpcConn)),
		SlashingSource:        remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		StakingSource:         remotestakingsource.NewSource(source, stakingtypes.NewQueryClient(source.GrpcConn)),
		VipcoinAccountsSource: remotevipcoinaccountssource.NewSource(source, accountstypes.NewQueryClient(source.GrpcConn)),
		VipcoinWalletsSource:  remotevipcoinwalletsssource.NewSource(source, walletstypes.NewQueryClient(source.GrpcConn)),
		VipcoinBankingSource:  remotevipcoinbankingsource.NewSource(source, bankingtypes.NewQueryClient(source.GrpcConn)),
		VipcoinAssetsSource:   remotevipcoinassetssource.NewSource(source, assetstypes.NewQueryClient(source.GrpcConn)),
	}, nil
}
