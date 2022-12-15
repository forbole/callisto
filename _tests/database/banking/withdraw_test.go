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

	"github.com/forbole/bdjuno/v2/database/types"
	bankingdb "github.com/forbole/bdjuno/v2/database/vipcoin/chain/banking"
)

func TestRepository_SaveWithdraws(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		msg []*bankingtypes.Withdraw
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				msg: []*bankingtypes.Withdraw{
					{
						BaseTransfer: bankingtypes.BaseTransfer{
							Id:     4,
							Asset:  "asset",
							Amount: 4000,
							Kind:   bankingtypes.TRANSFER_KIND_WITHDRAW,
							Extras: []*extratypes.Extra{
								{
									Kind: extratypes.EXTRA_KIND_COMMENT,
									Data: "Withdraw test",
								},
							},
							Timestamp: 0,
							TxHash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
						},
						Wallet: "wwwwljs7p2p9ae3en8knr3d3ke8srsfcj2zjvewwww",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			if err := r.SaveWithdraws(tt.args.msg...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveWithdraws() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetWithdraws(t *testing.T) {
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
		msg     []*bankingtypes.Withdraw
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				bfilter: filter.NewFilter().SetArgument(types.FieldID, "4"),
			},
			msg: []*bankingtypes.Withdraw{
				{
					BaseTransfer: bankingtypes.BaseTransfer{
						Id:     4,
						Asset:  "asset",
						Amount: 4000,
						Kind:   bankingtypes.TRANSFER_KIND_WITHDRAW,
						Extras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_COMMENT,
								Data: "Withdraw test",
							},
						},
						Timestamp: 10800,
						TxHash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
					},
					Wallet: "wwwwljs7p2p9ae3en8knr3d3ke8srsfcj2zjvewwww",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetWithdraws(tt.args.bfilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetWithdraws() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.msg) {
				t.Errorf("Repository.GetWithdraws() = %v, msg %v", got, tt.msg)
			}
		})
	}
}

func TestRepository_UpdateWithdraws(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		withdraws []*bankingtypes.Withdraw
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				withdraws: []*bankingtypes.Withdraw{
					{
						BaseTransfer: bankingtypes.BaseTransfer{
							Id:     4,
							Asset:  "asset",
							Amount: 4000,
							Kind:   bankingtypes.TRANSFER_KIND_WITHDRAW,
							Extras: []*extratypes.Extra{
								{
									Kind: extratypes.EXTRA_KIND_COMMENT,
									Data: "Withdraw test",
								},
							},
							Timestamp: 10800,
							TxHash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
						},
						Wallet: "wwwwljs7p2p9ae3en8knr3d3ke8srsfcj2zjvewwww",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid",
			args: args{
				withdraws: []*bankingtypes.Withdraw{
					{
						BaseTransfer: bankingtypes.BaseTransfer{
							Id:     4,
							Asset:  "asset",
							Amount: 4000,
							Kind:   bankingtypes.TRANSFER_KIND_WITHDRAW,
							Extras: []*extratypes.Extra{
								{
									Kind: extratypes.EXTRA_KIND_COMMENT,
									Data: "Withdraw test (update)",
								},
							},
							Timestamp: 10800,
							TxHash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
						},
						Wallet: "ovg12ndz92smw9pz34m7kfr5vqq6qscee7nv83test",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			// update wallet and extras.data
			if err := r.UpdateWithdraws(tt.args.withdraws...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateWithdraws() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
