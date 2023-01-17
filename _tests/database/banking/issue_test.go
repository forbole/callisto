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

func TestRepository_SaveIssues(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		msg []*bankingtypes.Issue
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				msg: []*bankingtypes.Issue{
					{
						BaseTransfer: bankingtypes.BaseTransfer{
							Id:     3,
							Asset:  "asset",
							Amount: 1000,
							Kind:   bankingtypes.TRANSFER_KIND_ISSUE,
							Extras: []*extratypes.Extra{
								{
									Kind: extratypes.EXTRA_KIND_COMMENT,
									Data: "Issue test",
								},
							},
							Timestamp: 0,
							TxHash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
						},
						Wallet: "wiwiwis7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			if err := r.SaveIssues(tt.args.msg...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveIssues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetIssues(t *testing.T) {
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
		msg     []*bankingtypes.Issue
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				bfilter: filter.NewFilter().SetArgument(types.FieldID, "3"),
			},
			msg: []*bankingtypes.Issue{
				{
					BaseTransfer: bankingtypes.BaseTransfer{
						Id:     3,
						Asset:  "asset",
						Amount: 1000,
						Kind:   bankingtypes.TRANSFER_KIND_ISSUE,
						Extras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_COMMENT,
								Data: "Issue test",
							},
						},
						Timestamp: 10800,
						TxHash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
					},
					Wallet: "wiwiwis7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetIssues(tt.args.bfilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetIssues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.msg) {
				t.Errorf("Repository.GetIssues() = %v, msg %v", got, tt.msg)
			}
		})
	}
}

func TestRepository_UpdateIssues(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		issues []*bankingtypes.Issue
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				issues: []*bankingtypes.Issue{
					{
						BaseTransfer: bankingtypes.BaseTransfer{
							Id:     3,
							Asset:  "asset",
							Amount: 1000,
							Kind:   bankingtypes.TRANSFER_KIND_ISSUE,
							Extras: []*extratypes.Extra{
								{
									Kind: extratypes.EXTRA_KIND_COMMENT,
									Data: "Issue test",
								},
							},
							Timestamp: 10800,
							TxHash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
						},
						Wallet: "wiwiwis7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid",
			args: args{
				issues: []*bankingtypes.Issue{
					{
						BaseTransfer: bankingtypes.BaseTransfer{
							Id:     3,
							Asset:  "asset",
							Amount: 10000,
							Kind:   bankingtypes.TRANSFER_KIND_ISSUE,
							Extras: []*extratypes.Extra{
								{
									Kind: extratypes.EXTRA_KIND_COMMENT,
									Data: "Issue test (update)",
								},
							},
							Timestamp: 10800,
							TxHash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
						},
						Wallet: "wiwiwis7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			// update amount and extras.data
			if err := r.UpdateIssues(tt.args.issues...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateIssues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
