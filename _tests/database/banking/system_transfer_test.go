package banking

import (
	"reflect"
	"testing"

	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	extratypes "git.ooo.ua/vipcoin/chain/x/types"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/stdlib"

	bankingdb "github.com/forbole/bdjuno/v2/database/overgold/chain/banking"
	"github.com/forbole/bdjuno/v2/database/types"
)

func TestRepository_SaveSystemTransfers(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		msg []*bankingtypes.SystemTransfer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				msg: []*bankingtypes.SystemTransfer{
					{
						BaseTransfer: bankingtypes.BaseTransfer{
							Id:        2,
							Asset:     "asset",
							Amount:    2000,
							Kind:      bankingtypes.TRANSFER_KIND_SYSTEM,
							Extras:    []*extratypes.Extra{},
							Timestamp: 0,
							TxHash:    "b835ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
						},
						WalletFrom: "ssssljs7p2p9ae3en8knr3d3ke8srsfcj2zjvessss",
						WalletTo:   "sssslws9pab9ae3en8d0r3d3ke8srsfcj2zjvessss",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			if err := r.SaveSystemTransfers(tt.args.msg...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveSystemTransfers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetSystemTransfers(t *testing.T) {
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
		msg     []*bankingtypes.SystemTransfer
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				bfilter: filter.NewFilter().SetArgument(types.FieldID, "2"),
			},
			msg: []*bankingtypes.SystemTransfer{
				{
					BaseTransfer: bankingtypes.BaseTransfer{
						Id:        2,
						Asset:     "asset",
						Amount:    2000,
						Kind:      bankingtypes.TRANSFER_KIND_SYSTEM,
						Extras:    []*extratypes.Extra{},
						Timestamp: 10800,
						TxHash:    "b835ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
					},
					WalletFrom: "ssssljs7p2p9ae3en8knr3d3ke8srsfcj2zjvessss",
					WalletTo:   "sssslws9pab9ae3en8d0r3d3ke8srsfcj2zjvessss",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetSystemTransfers(tt.args.bfilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetSystemTransfers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.msg) {
				t.Errorf("Repository.GetSystemTransfers() = %v, msg %v", got, tt.msg)
			}
		})
	}
}

func TestRepository_UpdateSystemTransfers(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		transfers []*bankingtypes.SystemTransfer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				transfers: []*bankingtypes.SystemTransfer{
					{
						BaseTransfer: bankingtypes.BaseTransfer{
							Id:        2,
							Asset:     "asset",
							Amount:    2000,
							Kind:      bankingtypes.TRANSFER_KIND_SYSTEM,
							Extras:    []*extratypes.Extra{},
							Timestamp: 10800,
							TxHash:    "b835ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
						},
						WalletFrom: "ssssljs7p2p9ae3en8knr3d3ke8srsfcj2zjvessss",
						WalletTo:   "sssslws9pab9ae3en8d0r3d3ke8srsfcj2zjvessss",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid",
			args: args{
				transfers: []*bankingtypes.SystemTransfer{
					{
						BaseTransfer: bankingtypes.BaseTransfer{
							Id:        2,
							Asset:     "new asset",
							Amount:    22200,
							Kind:      bankingtypes.TRANSFER_KIND_SYSTEM,
							Extras:    []*extratypes.Extra{},
							Timestamp: 10800,
							TxHash:    "b835ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
						},
						WalletFrom: "ssssljs7p2p9ae3en8knr3d3ke8srsfcj2zjvessss",
						WalletTo:   "sssslws9pab9ae3en8d0r3d3ke8srsfcj2zjvessss",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			// update asset and amount
			if err := r.UpdateSystemTransfers(tt.args.transfers...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateSystemTransfers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
