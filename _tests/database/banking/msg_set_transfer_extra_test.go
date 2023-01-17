package banking

import (
	"reflect"
	"testing"

	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	extratypes "git.ooo.ua/vipcoin/chain/x/types"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/cosmos/cosmos-sdk/simapp"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	bankingdb "github.com/forbole/bdjuno/v2/database/overgold/chain/banking"
	"github.com/forbole/bdjuno/v2/database/types"
)

func TestRepository_SaveMsgSetTransferExtra(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		extras []*bankingtypes.MsgSetTransferExtra
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				extras: []*bankingtypes.MsgSetTransferExtra{
					{
						Creator: "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
						Id:      1,
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
						Creator: "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv2g",
						Id:      2,
						Extras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_EMAIL,
								Data: "retg@mail",
							},
							{
								Kind: extratypes.EXTRA_KIND_PHONE,
								Data: "+380685549",
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

			if err := r.SaveMsgSetTransferExtra(tt.args.extras[0], ""); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveMsgSetTransferExtra() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetMsgSetTransferExtra(t *testing.T) {
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
		want    []*bankingtypes.MsgSetTransferExtra
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				filter: filter.NewFilter().SetArgument(types.FieldID, 1),
			},
			want: []*bankingtypes.MsgSetTransferExtra{
				{
					Creator: "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
					Id:      1,
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

			got, err := r.GetMsgSetTransferExtra(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetMsgSetTransferExtra() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetMsgSetTransferExtra() = %v, want %v", got, tt.want)
			}
		})
	}
}
