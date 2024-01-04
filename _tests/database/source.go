package database

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v4/database/overgold/chain/allowed"
	"github.com/forbole/bdjuno/v4/database/overgold/chain/bank"
	"github.com/forbole/bdjuno/v4/database/overgold/chain/core"
	"github.com/forbole/bdjuno/v4/database/overgold/chain/feeexcluder"
	"github.com/forbole/bdjuno/v4/database/overgold/chain/last_block"
	"github.com/forbole/bdjuno/v4/database/overgold/chain/referral"
	"github.com/forbole/bdjuno/v4/database/overgold/chain/stake"
)

const (
	host       = "localhost"
	dbName     = "juno"
	dbPassword = "postgres"
)

const (
	TestAddressCreator = "ovg18p9heyy3m4dsq7fj86p6v9yzzx8a64f86eec7u"
)

var (
	DB    *sqlx.DB
	Codec codec.Codec

	Datastore struct {
		Allowed     *allowed.Repository
		Bank        *bank.Repository
		Core        *core.Repository
		FeeExcluder *feeexcluder.Repository
		LastBlock   *last_block.Repository
		Referral    *referral.Repository
		Stake       *stake.Repository
	}
)

func init() {
	var err error

	DB, err = sqlx.Connect("pgx", fmt.Sprintf("host=%s port=5432 user=postgres dbname=%s password=%s sslmode=disable",
		host,
		dbName,
		dbPassword,
	))
	if err != nil {
		log.Fatal().Err(err)
	}

	// Create the codec.
	// TODO: rework it: Codec = registrar.Context{}.EncodingConfig.Codec

	// OverGold modules
	Datastore.Allowed = allowed.NewRepository(DB, Codec)
	Datastore.Core = core.NewRepository(DB, Codec)
	Datastore.FeeExcluder = feeexcluder.NewRepository(DB, Codec)
	Datastore.Referral = referral.NewRepository(DB, Codec)
	Datastore.Stake = stake.NewRepository(DB, Codec)

	// Cosmos modules
	Datastore.Bank = bank.NewRepository(DB)
	Datastore.LastBlock = last_block.NewRepository(DB)
}
