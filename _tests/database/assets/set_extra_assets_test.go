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

func TestRepository_SaveSetExtraAsset(t *testing.T) {
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
		msgAssetSetExtra []*assetstypes.MsgAssetSetExtra
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] save message set extra assets into db",
			args: args{
				msgAssetSetExtra: []*assetstypes.MsgAssetSetExtra{
					{
						Creator: "vcg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdz",
						Name:    "This is Name for assets name",
						Extras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_WALLET,
								Data: "walletHolderName1",
							},
							{
								Kind: extratypes.EXTRA_KIND_WALLET,
								Data: "walletHolderName2",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := assetsdb.NewRepository(db, codec.Marshaler)

			if err := r.SaveExtraAsset(tt.args.msgAssetSetExtra[0], ""); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveSetExtraAssets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetSetExtraAsset(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name             string
		args             args
		msgAssetSetExtra []*assetstypes.MsgAssetSetExtra
		wantErr          bool
	}{
		{
			name: "[success] get set extra assets from db",
			args: args{
				filter: filter.NewFilter().SetArgument("name", "This is Name for assets name"),
			},
			msgAssetSetExtra: []*assetstypes.MsgAssetSetExtra{
				{
					Creator: "vcg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdz",
					Name:    "This is Name for assets name",
					Extras: []*extratypes.Extra{
						{
							Kind: extratypes.EXTRA_KIND_WALLET,
							Data: "walletHolderName1",
						},
						{
							Kind: extratypes.EXTRA_KIND_WALLET,
							Data: "walletHolderName2",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := assetsdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetExtraAsset(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetExtraAsset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.msgAssetSetExtra) {
				t.Errorf("Repository.GetExtraAsset() = %v, msg %v", got, tt.msgAssetSetExtra)
			}
		})
	}
}
