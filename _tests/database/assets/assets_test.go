package assets

import (
	"reflect"
	"testing"

	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
	extratypes "git.ooo.ua/vipcoin/chain/x/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/cosmos/cosmos-sdk/simapp"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	assetsdb "github.com/forbole/bdjuno/v3/database/overgold/chain/assets"
)

func TestRepository_SaveAsset(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	codec := simapp.MakeTestEncodingConfig()

	if err != nil {
		t.Error(err)
		return
	}

	type args struct {
		assets []*assetstypes.Asset
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] save assets into db",
			args: args{
				assets: []*assetstypes.Asset{
					{
						Issuer: "ovg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdz",
						Name:   "This is Name for assets name",
						Policies: []assetstypes.AssetPolicy{
							assetstypes.ASSET_POLICY_ISSUABLE,
							assetstypes.ASSET_POLICY_EXCHANGEABLE,
						},
						State:         assetstypes.ASSET_STATE_ACTIVE,
						Issued:        785,
						Burned:        4556,
						Withdrawn:     658655656,
						InCirculation: 454562,
						Properties: assetstypes.Properties{
							Precision:  8,
							FeePercent: 100,
						},
						Extras: []*extratypes.Extra{},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := assetsdb.NewRepository(db, codec.Marshaler)

			if err := r.SaveAssets(tt.args.assets...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveAssets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetAsset(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		accfilter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		assets  []*assetstypes.Asset
		wantErr bool
	}{
		{
			name: "[success] get assets from db",
			args: args{
				accfilter: filter.NewFilter().SetArgument("issuer", "this is isseuer"),
			},
			assets: []*assetstypes.Asset{
				{
					Issuer: "ovg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdy",
					Name:   assetstypes.AssetOVG,
					Policies: []assetstypes.AssetPolicy{
						assetstypes.ASSET_POLICY_ISSUABLE,
					},
					State:         assetstypes.ASSET_STATE_INACTIVE,
					Issued:        6785,
					Burned:        64556,
					Withdrawn:     7658655656,
					InCirculation: 8454562,
					Properties: assetstypes.Properties{
						Precision:  9,
						FeePercent: 1000,
					},
					Extras: []*extratypes.Extra{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := assetsdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetAssets(tt.args.accfilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetAsset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.assets) {
				t.Errorf("Repository.GetAsset() = %v, msg %v", got, tt.assets)
			}
		})
	}
}

func TestRepository_UpdateAsset(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	codec := simapp.MakeTestEncodingConfig()

	if err != nil {
		t.Error(err)
		return
	}

	type args struct {
		assets []*assetstypes.Asset
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] update assets in db",
			args: args{
				assets: []*assetstypes.Asset{
					{
						Issuer: "ovg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdz",
						Name:   assetstypes.AssetOVG,
						Policies: []assetstypes.AssetPolicy{
							assetstypes.ASSET_POLICY_ISSUABLE,
						},
						State:         assetstypes.ASSET_STATE_INACTIVE,
						Issued:        6785,
						Burned:        64556,
						Withdrawn:     7658655656,
						InCirculation: 8454562,
						Properties: assetstypes.Properties{
							Precision:  9,
							FeePercent: 1000,
						},
						Extras: []*extratypes.Extra{},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := assetsdb.NewRepository(db, codec.Marshaler)

			if err := r.UpdateAssets(tt.args.assets...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateAssets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
