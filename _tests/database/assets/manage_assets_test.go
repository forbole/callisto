package assets

import (
	"reflect"
	"testing"

	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/cosmos/cosmos-sdk/simapp"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	assetsdb "github.com/forbole/bdjuno/v2/database/vipcoin/chain/assets"
)

func TestRepository_SaveManageAsset(t *testing.T) {
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
		msgAssetManage []*assetstypes.MsgAssetManage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] save message assets into db",
			args: args{
				msgAssetManage: []*assetstypes.MsgAssetManage{
					{
						Creator: "vcg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdz",
						Name:    "This is Name for assets name",
						Policies: []assetstypes.AssetPolicy{
							assetstypes.ASSET_POLICY_ISSUABLE,
							assetstypes.ASSET_POLICY_EXCHANGEABLE,
						},
						State: assetstypes.ASSET_STATE_ACTIVE,
						Properties: assetstypes.Properties{
							Precision:  8,
							FeePercent: 100,
						},
						Issued:        6785,
						Burned:        64556,
						Withdrawn:     7658655656,
						InCirculation: 8454562,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := assetsdb.NewRepository(db, codec.Marshaler)

			if err := r.SaveManageAsset(tt.args.msgAssetManage, ""); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveManageAssets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetManageAsset(t *testing.T) {
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
		msgAssetManage []*assetstypes.MsgAssetManage
		wantErr        bool
	}{
		{
			name: "[success] get manage assets from db",
			args: args{
				accfilter: filter.NewFilter().SetArgument("creator", "vcg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdz"),
			},
			msgAssetManage: []*assetstypes.MsgAssetManage{
				{
					Creator: "vcg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdz",
					Name:    "This is Name for assets name",
					Policies: []assetstypes.AssetPolicy{
						assetstypes.ASSET_POLICY_ISSUABLE,
						assetstypes.ASSET_POLICY_EXCHANGEABLE,
					},
					State: assetstypes.ASSET_STATE_ACTIVE,
					Properties: assetstypes.Properties{
						Precision:  8,
						FeePercent: 100,
					},
					Issued:        6785,
					Burned:        64556,
					Withdrawn:     7658655656,
					InCirculation: 8454562,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := assetsdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetManageAsset(tt.args.accfilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetManageAsset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.msgAssetManage) {
				t.Errorf("Repository.GetManageAsset() = %v, msg %v", got, tt.msgAssetManage)
			}
		})
	}
}
