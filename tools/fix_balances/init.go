/*
 * Copyright (c) 2023. Business Process Technologies. All rights reserved.
 */

package main

import (
	"embed"

	"git.ooo.ua/vipcoin/chain-client/pkg/api/v4/chain/assets"
	"git.ooo.ua/vipcoin/chain-client/pkg/client"
	"git.ooo.ua/vipcoin/lib/database"
	"git.ooo.ua/vipcoin/lib/log"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	_ "github.com/jackc/pgx/v5/stdlib"

	"fix_balances/config"
)

// initDependencies - initialize dependencies
func initDependencies() (logger log.Logger, db database.Executor, ovgChainClient client.VCG, err error) {
	cfg := config.Config

	// init logger
	logger = log.InitLogger("update_balance_in_database_from_blockchain", "", cfg.Logger, nil)

	// init database
	db = database.InitDatabase(cfg.Database, logger, embed.FS{})

	ovgChainClient, err = client.NewVCGClient(
		client.Options{
			Chain:         cfg.Chain.ChainID,
			BaseURL:       cfg.Chain.GRPCBaseURL,
			KeyringType:   keyring.BackendMemory,
			GasLimit:      cfg.Chain.GasLimit,
			AccountPrefix: assets.AssetOVG,
		},
	)
	if err != nil {
		logger.Fatal(err)
	}

	return logger, db, ovgChainClient, nil
}
