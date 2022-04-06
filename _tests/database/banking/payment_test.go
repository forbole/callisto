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

	bankingdb "github.com/forbole/bdjuno/v2/database/vipcoin/chain/banking"
)

func TestRepository_SavePayments(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		msg []*bankingtypes.Payment
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				msg: []*bankingtypes.Payment{
					{
						BaseTransfer: bankingtypes.BaseTransfer{
							Id:        1,
							Asset:     "asset",
							Amount:    1000,
							Kind:      bankingtypes.TRANSFER_KIND_PAYMENT,
							Extras:    []*extratypes.Extra{},
							Timestamp: 0,
							TxHash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
						},
						WalletFrom: "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
						WalletTo:   "aaa5lws9pab9ae3en8d0r3d3ke8srsfcj2zjvefzzz",
						Fee:        5,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			if err := r.SavePayments(tt.args.msg...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SavePayments() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetPayments(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		bfilter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		msg     []*bankingtypes.Payment
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				bfilter: filter.NewFilter().SetArgument("id", "1"),
			},
			msg: []*bankingtypes.Payment{
				{
					BaseTransfer: bankingtypes.BaseTransfer{
						Id:        1,
						Asset:     "asset",
						Amount:    1000,
						Kind:      bankingtypes.TRANSFER_KIND_PAYMENT,
						Extras:    []*extratypes.Extra{},
						Timestamp: 10800,
						TxHash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
					},
					WalletFrom: "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
					WalletTo:   "aaa5lws9pab9ae3en8d0r3d3ke8srsfcj2zjvefzzz",
					Fee:        5,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetPayments(tt.args.bfilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetPayments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.msg) {
				t.Errorf("Repository.GetPayments() = %v, msg %v", got, tt.msg)
			}
		})
	}
}

func TestRepository_UpdatePayments(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		payments []*bankingtypes.Payment
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				payments: []*bankingtypes.Payment{
					{
						BaseTransfer: bankingtypes.BaseTransfer{
							Id:        1,
							Asset:     "asset",
							Amount:    1000,
							Kind:      bankingtypes.TRANSFER_KIND_PAYMENT,
							Extras:    []*extratypes.Extra{},
							Timestamp: 10800,
							TxHash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
						},
						WalletFrom: "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
						WalletTo:   "aaa5lws9pab9ae3en8d0r3d3ke8srsfcj2zjvefzzz",
						Fee:        5,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid",
			args: args{
				payments: []*bankingtypes.Payment{
					{
						BaseTransfer: bankingtypes.BaseTransfer{
							Id:        1,
							Asset:     "new asset",
							Amount:    1000,
							Kind:      bankingtypes.TRANSFER_KIND_PAYMENT,
							Extras:    []*extratypes.Extra{},
							Timestamp: 10800,
							TxHash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
						},
						WalletFrom: "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
						WalletTo:   "aaa5lws9pab9ae3en8d0r3d3ke8srsfcj2zjvefzzz",
						Fee:        10,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			// update fee and asset
			if err := r.UpdatePayments(tt.args.payments...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdatePayments() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
