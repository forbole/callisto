package assets

import (
	"context"
	"database/sql"

	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/bdjuno/v2/database/types"
	"github.com/jmoiron/sqlx"
)

type (
	// Repository - defines a repository for assets repository
	Repository struct {
		db  *sqlx.DB
		cdc codec.Marshaler
	}
)

// NewRepository constructor
func NewRepository(db *sqlx.DB, cdc codec.Marshaler) *Repository {
	return &Repository{
		db:  db,
		cdc: cdc,
	}
}

// SaveAssets - method that save assets to the "vipcoin_chain_assets_assets" table
func (r Repository) SaveAssets(assets ...*assetstypes.Asset) error {
	if len(assets) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_assets_assets 
       (
        "issuer", "name", "policies", "state", "issued", "burned",
        "withdrawn", "in_circulation", "precision", "fee_percent", "extras"
        ) 
     VALUES 
       (
        :issuer, :name, :policies, :state, :issued, :burned,
        :withdrawn, :in_circulation, :precision, :fee_percent, :extras
        )`

	if _, err := r.db.NamedExec(query, toAssetsArrDatabase(assets...)); err != nil {
		return err
	}

	return nil
}

// GetAssets - method that get assets from the "vipcoin_chain_assets_assets" table
func (r Repository) GetAssets(filter filter.Filter) ([]*assetstypes.Asset, error) {
	query, args := filter.Build("vipcoin_chain_assets_assets",
		`issuer, name, policies, state, issued, burned, withdrawn, in_circulation, precision, fee_percent, extras`)

	var result []types.DBAssets
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*assetstypes.Asset{}, err
	}

	assets := make([]*assetstypes.Asset, 0, len(result))
	for _, w := range result {
		assets = append(assets, toAssetDomain(w))
	}

	return assets, nil
}

// UpdateAssets - method that updates the assets in the "vipcoin_chain_assets_assets" table
func (r Repository) UpdateAssets(assets ...*assetstypes.Asset) error {
	if len(assets) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := `UPDATE vipcoin_chain_assets_assets SET
				 issuer = :issuer, name = :name, policies = :policies,
				 state = :state, issued = :issued, burned = :burned, withdrawn = :withdrawn, 
			     in_circulation = :in_circulation, precision = :precision, fee_percent = :fee_percent, extras = :extras
			 WHERE issuer = :issuer`

	for _, asset := range assets {
		assetDB := toAssetDatabase(asset)

		if _, err := tx.NamedExec(query, assetDB); err != nil {
			return err
		}
	}

	return tx.Commit()
}
