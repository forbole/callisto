package wallets

import (
	"testing"

	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/cosmos/cosmos-sdk/simapp"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v3/database/types"

	walletsdb "github.com/forbole/bdjuno/v3/database/overgold/chain/wallets"
)

// TestRepository_SaveSetCreateUserWalletPrice - test for SaveSetCreateUserWalletPrice method
func TestRepository_SaveSetCreateUserWalletPrice(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		msg *walletstypes.MsgSetCreateUserWalletPrice
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] save MsgSetCreateUserWalletPrice",
			args: args{
				msg: &walletstypes.MsgSetCreateUserWalletPrice{
					Creator: "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
					Amount:  100000,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := walletsdb.NewRepository(db, codec.Marshaler)

			if err = r.SaveSetCreateUserWalletPrice(tt.args.msg, ""); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveSetCreateUserWalletPrice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestRepository_GetSetCreateUserWalletPrice - test for GetSetCreateUserWalletPrice method
func TestRepository_GetSetCreateUserWalletPrice(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		msg     *walletstypes.MsgSetCreateUserWalletPrice
		wantErr bool
	}{
		{
			name: "[success] get GetSetCreateUserWalletPrice",
			args: args{
				filter: filter.NewFilter().SetSort(types.FieldID, filter.DirectionDescending),
			},
			msg: &walletstypes.MsgSetCreateUserWalletPrice{
				Creator: "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
				Amount:  1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := walletsdb.NewRepository(db, codec.Marshaler)

			if err = r.SaveSetCreateUserWalletPrice(tt.msg, "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g"); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveSetCreateUserWalletPrice() error = %v, wantErr %v", err, tt.wantErr)
			}

			if _, err := r.GetSetCreateUserWalletPrice(tt.args.filter); (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetSetCreateUserWalletPrice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestRepository_GetAllSetCreateUserWalletPrice - test for GetAllSetCreateUserWalletPrice method
func TestRepository_GetAllSetCreateUserWalletPricey(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		msg     *walletstypes.MsgSetCreateUserWalletPrice
		wantErr bool
	}{
		{
			name: "[success] get GetAllSetCreateUserWalletPrice",
			args: args{
				filter: filter.NewFilter().SetSort(types.FieldID, filter.DirectionDescending),
			},
			msg: &walletstypes.MsgSetCreateUserWalletPrice{
				Creator: "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
				Amount:  1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := walletsdb.NewRepository(db, codec.Marshaler)

			if err = r.SaveSetCreateUserWalletPrice(tt.msg, "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g"); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveSetCreateUserWalletPrice() error = %v, wantErr %v", err, tt.wantErr)
			}

			if _, err := r.GetAllSetCreateUserWalletPrice(tt.args.filter); (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetAllSetCreateUserWalletPrice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
