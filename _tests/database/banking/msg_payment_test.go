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

func TestRepository_SaveMsgPayments(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		payments []*bankingtypes.MsgPayment
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				payments: []*bankingtypes.MsgPayment{
					{
						Creator:    "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
						WalletFrom: "vcg10xq8z7refmx3eqkv2ym76jkw72xkd5949k30gj",
						WalletTo:   "vcg1fvz25fze4c7cwg2xvv6p0trlf34w9pu5m4vwk5",
						Asset:      assets.AssetVCG,
						Amount:     1000,
						Extras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_EMAIL,
								Data: "retg@mail",
							},
							{
								Kind: extratypes.EXTRA_KIND_PHONE,
								Data: "+380685548",
							},
						},
					},
					{
						Creator:    "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv2g",
						WalletFrom: "vcg10xq8z7refmx3eqkv2ym76jkw72xkd5949k33gj",
						WalletTo:   "vcg1fvz25fze4c7cwg2xvv6p0trlf34w9pu5m4vwk4",
						Asset:      assets.AssetVCG,
						Amount:     1001,
						Extras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_EMAIL,
								Data: "retg@mail",
							},
							{
								Kind: extratypes.EXTRA_KIND_PHONE,
								Data: "+380685548",
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

			if err := r.SaveMsgPayments(tt.args.payments...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SavePayments() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetMsgPayments(t *testing.T) {
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
		want    []*bankingtypes.MsgPayment
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				filter: filter.NewFilter().SetArgument(types.FieldCreator, "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g"),
			},
			want: []*bankingtypes.MsgPayment{
				{
					Creator:    "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
					WalletFrom: "vcg10xq8z7refmx3eqkv2ym76jkw72xkd5949k30gj",
					WalletTo:   "vcg1fvz25fze4c7cwg2xvv6p0trlf34w9pu5m4vwk5",
					Asset:      assets.AssetVCG,
					Amount:     1000,
					Extras: []*extratypes.Extra{
						{
							Kind: extratypes.EXTRA_KIND_EMAIL,
							Data: "retg@mail",
						},
						{
							Kind: extratypes.EXTRA_KIND_PHONE,
							Data: "+380685548",
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

			got, err := r.GetMsgPayments(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetPayments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetPayments() = %v, want %v", got, tt.want)
			}
		})
	}
}
