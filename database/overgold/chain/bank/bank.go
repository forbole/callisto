package bank

import (
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v4/database/overgold/chain"
)

var _ chain.Bank = &Repository{}

type (
	// Repository - defines a repository for bank repository
	Repository struct {
		db *sqlx.DB
	}
)

// NewRepository constructor.
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}
