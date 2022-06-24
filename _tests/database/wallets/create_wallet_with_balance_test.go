package wallets

import (
	"reflect"
	"testing"

	extratypes "git.ooo.ua/vipcoin/chain/x/types"
	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	typesdb "github.com/forbole/bdjuno/v2/database/types"
	walletsdb "github.com/forbole/bdjuno/v2/database/vipcoin/chain/wallets"
)

func TestRepository_SaveCreateWalletWithBalance(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		msg []*walletstypes.MsgCreateWalletWithBalance
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				msg: []*walletstypes.MsgCreateWalletWithBalance{
					{
						Creator:        "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
						Address:        "vcg1hwxejcutmafuedd8trjqumfdkst2498pggx45q",
						AccountAddress: "vcg1hwxejcutmafuedd8trjqumfdkst2498pggx45q",
						Kind:           walletstypes.WALLET_KIND_HOLDER,
						State:          walletstypes.WALLET_STATE_ACTIVE,
						Extras:         []*extratypes.Extra{},
						Default:        true,
						Balance:        types.Coins{},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := walletsdb.NewRepository(db, codec.Marshaler)

			if err := r.SaveCreateWalletWithBalance(tt.args.msg, ""); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveCreateWalletWithBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetCreateWalletWithBalance(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		walletFilter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		msg     []*walletstypes.MsgCreateWalletWithBalance
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				walletFilter: filter.NewFilter().SetArgument(typesdb.FieldCreator, "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g"),
			},
			msg: []*walletstypes.MsgCreateWalletWithBalance{
				{
					Creator:        "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
					Address:        "vcg1hwxejcutmafuedd8trjqumfdkst2498pggx45q",
					AccountAddress: "vcg1hwxejcutmafuedd8trjqumfdkst2498pggx45q",
					Kind:           walletstypes.WALLET_KIND_HOLDER,
					State:          walletstypes.WALLET_STATE_ACTIVE,
					Extras:         []*extratypes.Extra{},
					Default:        true,
					Balance:        types.Coins{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := walletsdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetCreateWalletWithBalance(tt.args.walletFilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetCreateWalletWithBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.msg) {
				t.Errorf("Repository.GetCreateWalletWithBalance() = %v, msg %v", got, tt.msg)
			}
		})
	}
}
