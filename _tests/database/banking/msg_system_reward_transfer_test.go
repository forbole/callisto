package banking

import (
	"reflect"
	"testing"

	assets "git.ooo.ua/vipcoin/chain/x/assets/types"
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	extratypes "git.ooo.ua/vipcoin/chain/x/types"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/cosmos/cosmos-sdk/simapp"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v2/database/types"
	bankingdb "github.com/forbole/bdjuno/v2/database/vipcoin/chain/banking"
)

func TestRepository_SaveMsgSystemRewardTransfers(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		transfers []*bankingtypes.MsgSystemRewardTransfer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				transfers: []*bankingtypes.MsgSystemRewardTransfer{
					{
						Creator:    "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv2g",
						WalletFrom: "vcg1k2gx4u0hwk87ja3wyakne8cl5gytnz0uc27xm4",
						WalletTo:   "vcg10z0f2qdu483pn04cpfa9ux9njk8jplk9y8hfs5",
						Asset:      assets.AssetVCG,
						Amount:     1000,
						Extras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_COMMENT,
								Data: "MsgSystemRewardTransfer",
							},
						},
					},
					{
						Creator:    "vcg2ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv3g",
						WalletFrom: "vcg1k2gx4u0hwk87ja3wyakne8cl5gytnz0uc27xm1",
						WalletTo:   "vcg10z0f2qdu483pn04cpfa9ux9njk8jplk9y8hfs2",
						Asset:      assets.AssetVCG,
						Amount:     2000,
						Extras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_COMMENT,
								Data: "MsgSystemRewardTransfer",
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
			r := bankingdb.NewRepository(db, codec.Marshaler)

			if err := r.SaveMsgSystemRewardTransfers(tt.args.transfers[0], ""); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveMsgSystemRewardTransfers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetMsgSystemRewardTransfers(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		want    []*bankingtypes.MsgSystemRewardTransfer
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				filter: filter.NewFilter().SetArgument(types.FieldCreator, "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv2g"),
			},
			want: []*bankingtypes.MsgSystemRewardTransfer{
				{
					Creator:    "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv2g",
					WalletFrom: "vcg1k2gx4u0hwk87ja3wyakne8cl5gytnz0uc27xm4",
					WalletTo:   "vcg10z0f2qdu483pn04cpfa9ux9njk8jplk9y8hfs5",
					Asset:      assets.AssetVCG,
					Amount:     1000,
					Extras: []*extratypes.Extra{
						{
							Kind: extratypes.EXTRA_KIND_COMMENT,
							Data: "MsgSystemRewardTransfer",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetMsgSystemRewardTransfers(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetMsgSystemRewardTransfers() error = %v\nwantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetMsgSystemRewardTransfers() = %v\nwant %v", got, tt.want)
			}
		})
	}
}
