package feeexcluder

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/jmoiron/sqlx"
)

type (
	// Repository - defines a repository for fee exclude repository
	Repository struct {
		cdc codec.Codec
		db  *sqlx.DB
	}
)

// NewRepository constructor.
func NewRepository(db *sqlx.DB, cdc codec.Codec) *Repository {
	return &Repository{
		cdc: cdc,
		db:  db,
	}
}
