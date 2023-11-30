package last_block

import (
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v4/database/overgold/chain"
)

var _ chain.LastBlock = &Repository{}

type (
	// Repository - defines a repository for last block repository
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

// create - define repository method for save last block.
func (r Repository) create(id uint64) error {
	query := `INSERT INTO last_block (block) VALUES ($1)`

	if _, err := r.db.Exec(query, id); err != nil {
		return err
	}

	return nil
}

// Get - define repository method which gets last block from db.
func (r Repository) Get() (uint64, error) {
	query := `SELECT block FROM last_block`

	var blockNum uint64
	if err := r.db.Get(&blockNum, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, r.create(0)
		}

		return 0, errs.Internal{Cause: err.Error()}
	}

	return blockNum, nil
}

// Update - define repository method for update last block.
func (r Repository) Update(id uint64) error {
	query := `UPDATE last_block SET block = $1`

	if _, err := r.db.Exec(query, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}
