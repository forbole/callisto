package wallets

import (
	"reflect"
	"testing"

	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/cosmos/cosmos-sdk/simapp"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	walletsdb "github.com/forbole/bdjuno/v2/database/overgold/chain/wallets"
	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

func TestRepository_SaveMsgSetDefaultWallet(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		msg []*walletstypes.MsgSetDefaultWallet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				msg: []*walletstypes.MsgSetDefaultWallet{
					{
						Creator: "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
						Address: "ovg1hwxejcutmafuedd8trjqumfdkst2498pggx45q",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := walletsdb.NewRepository(db, codec.Marshaler)

			if err := r.SaveDefaultWallets(tt.args.msg[0], ""); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveKinds() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetMsgSetDefaultWallet(t *testing.T) {
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
		want    []*walletstypes.MsgSetDefaultWallet
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				filter: filter.NewFilter().SetArgument(dbtypes.FieldAddress, "ovg1hwxejcutmafuedd8trjqumfdkst2498pggx45q"),
			},
			want: []*walletstypes.MsgSetDefaultWallet{
				{
					Creator: "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
					Address: "ovg1hwxejcutmafuedd8trjqumfdkst2498pggx45q",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := walletsdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetDefaultWallets(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetDefaultWallets() error = %v\nwantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetDefaultWallets() = %v\nwant %v", got, tt.want)
			}
		})
	}
}
