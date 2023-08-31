/*
 * Copyright (c) 2023. Business Process Technologies. All rights reserved.
 */

package config

import (
	"fmt"
	"time"

	"git.ooo.ua/vipcoin/lib/database"
	"git.ooo.ua/vipcoin/lib/log"
)

// Custom Config variables.
const (
	// host address for DB
	host = "localhost"

	// credentials for DB
	credentials = "user=postgres dbname=postgres password=postgres" // local
)

type (
	// config defines configuration.
	config struct {
		Logger   log.Config
		Database database.Config
		Chain    Chain
	}

	// Chain defines cosmos api configuration.	SwapBalances string `yaml:"swap-balances" valid:"required"`.
	Chain struct {
		ChainID     string `yaml:"chain-id" valid:"required"`
		GRPCBaseURL string `yaml:"grpc-base-url" valid:"required"`
		GasLimit    uint64 `yaml:"gas-limit" valid:"required"`
	}
)

// Config exported variable
var Config = config{
	Logger: log.Config{
		Mode:                "dev",
		LogFormat:           "text",
		LogLevel:            "debug",
		DateTimeFormat:      "2006-01-02 15:04:05",
		UseTimestamp:        true,
		IncludeCallerMethod: true,
		OutputFilePath:      "",
	},
	Database: database.Config{
		ConnectionString: fmt.Sprintf("host=%s port=5432 %s sslmode=disable", host, credentials),
		Dialect:          "postgres",
		Driver:           "pgx",
		MaxRetries:       5,
		RetryDelay:       1 * time.Second,
		QueryTimeout:     60 * time.Second,
		AutoMigrate:      false,
	},
	Chain: Chain{
		ChainID:     "chain",
		GRPCBaseURL: "chain address",
		GasLimit:    100000000,
	},
}
