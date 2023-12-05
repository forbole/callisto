package database

import (
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v4/database/types"
)

// GetBlock - get block from database
func (db *Db) GetBlock(filter filter.Filter) (types.BlockRow, error) {
	query, args := filter.SetLimit(1).Build("block")

	var result types.BlockRow
	if err := db.Sqlx.Get(&result, query, args...); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return types.BlockRow{}, errs.Internal{Cause: err.Error()}
		}

		return types.BlockRow{}, errs.NotFound{What: "blocks"}
	}

	return result, nil
}

// GetBlocks - get blocks from database
func (db *Db) GetBlocks(filter filter.Filter) ([]types.BlockRow, error) {
	query, args := filter.Build("block")

	var val []types.BlockRow
	if err := db.Sqlx.Select(&val, query, args...); err != nil {
		return []types.BlockRow{}, errs.Internal{Cause: err.Error()}
	}

	if len(val) == 0 {
		return []types.BlockRow{}, errs.NotFound{What: "blocks"}
	}

	return val, nil
}
