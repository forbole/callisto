package wallets

import (
	"reflect"
	"testing"

	extratypes "git.ooo.ua/vipcoin/chain/x/types"
	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"

	"git.ooo.ua/vipcoin/lib/filter"

	walletsdb "github.com/forbole/bdjuno/v2/database/vipcoin/chain/wallets"
	_ "github.com/jackc/pgx/stdlib"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/jmoiron/sqlx"
)

func TestRepository_SaveCreateWallet(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		msg []*walletstypes.MsgCreateWallet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				msg: []*walletstypes.MsgCreateWallet{
					{
						Creator:        "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
						Address:        "vcg1hwxejcutmafuedd8trjqumfdkst2498pggx45q",
						AccountAddress: "vcg1hwxejcutmafuedd8trjqumfdkst2498pggx45q",
						Kind:           walletstypes.WALLET_KIND_HOLDER,
						State:          walletstypes.WALLET_STATE_ACTIVE,
						Extras:         []*extratypes.Extra{},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := walletsdb.NewRepository(db, codec.Marshaler)

			if err := r.SaveCreateWallet(tt.args.msg...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveCreateWallet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetCreateWallet(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		accfilter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		msg     []*walletstypes.MsgCreateWallet
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				accfilter: filter.NewFilter().SetArgument("creator", "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g"),
			},
			msg: []*walletstypes.MsgCreateWallet{
				{
					Creator:        "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
					Address:        "vcg1hwxejcutmafuedd8trjqumfdkst2498pggx45q",
					AccountAddress: "vcg1hwxejcutmafuedd8trjqumfdkst2498pggx45q",
					Kind:           walletstypes.WALLET_KIND_HOLDER,
					State:          walletstypes.WALLET_STATE_ACTIVE,
					Extras:         []*extratypes.Extra{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := walletsdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetCreateWallet(tt.args.accfilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetCreateWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.msg) {
				t.Errorf("Repository.GetCreateWallet() = %v, msg %v", got, tt.msg)
			}
		})
	}
}
