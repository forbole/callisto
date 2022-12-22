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

	assetsdb "github.com/forbole/bdjuno/v2/database/vipcoin/chain/assets"
)

func TestRepository_SaveCreateAsset(t *testing.T) {
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
		msgAssetCreate []*assetstypes.MsgAssetCreate
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] save create assets into db",
			args: args{
				msgAssetCreate: []*assetstypes.MsgAssetCreate{
					{
						Creator: "ovg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdz",
						Name:    "This is Name for assets name",
						Issuer:  "ovg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdz",
						Policies: []assetstypes.AssetPolicy{
							assetstypes.ASSET_POLICY_ISSUABLE,
							assetstypes.ASSET_POLICY_EXCHANGEABLE,
						},
						State: assetstypes.ASSET_STATE_ACTIVE,
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

			if err := r.SaveCreateAsset(tt.args.msgAssetCreate[0], ""); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveCreateAssets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetCreateAsset(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		accfilter filter.Filter
	}
	tests := []struct {
		name           string
		args           args
		msgAssetCreate []*assetstypes.MsgAssetCreate
		wantErr        bool
	}{
		{
			name: "[success] get create assets from db",
			args: args{
				accfilter: filter.NewFilter().SetArgument("creator", "ovg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdz"),
			},
			msgAssetCreate: []*assetstypes.MsgAssetCreate{
				{
					Creator: "ovg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdz",
					Name:    "This is Name for assets name",
					Issuer:  "ovg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdz",
					Policies: []assetstypes.AssetPolicy{
						assetstypes.ASSET_POLICY_ISSUABLE,
						assetstypes.ASSET_POLICY_EXCHANGEABLE,
					},
					State: assetstypes.ASSET_STATE_ACTIVE,
					Properties: assetstypes.Properties{
						Precision:  8,
						FeePercent: 100,
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

			got, err := r.GetCreateAsset(tt.args.accfilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetCreateAsset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.msgAssetCreate) {
				t.Errorf("Repository.GetCreateAsset() = %v, msg %v", got, tt.msgAssetCreate)
			}
		})
	}
}
