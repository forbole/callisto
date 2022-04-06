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

func TestRepository_SaveBaseTransfers(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		transfers []*bankingtypes.BaseTransfer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				transfers: []*bankingtypes.BaseTransfer{
					{
						Id:     7,
						Asset:  assets.AssetVCG,
						Amount: 4005,
						Kind:   bankingtypes.TRANSFER_KIND_ISSUE,
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
						Timestamp: 1649174827,
						TxHash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd7",
					},
					{
						Id:     8,
						Asset:  assets.AssetVCG,
						Amount: 4006,
						Kind:   bankingtypes.TRANSFER_KIND_DEFERRED,
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
						Timestamp: 1649174828,
						TxHash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd8",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			if err := r.SaveBaseTransfers(tt.args.transfers...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveBaseTransfers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_UpdateBaseTransfers(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		transfers []*bankingtypes.BaseTransfer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				transfers: []*bankingtypes.BaseTransfer{
					{
						Id:     7,
						Asset:  assets.AssetVCG,
						Amount: 4005,
						Kind:   bankingtypes.TRANSFER_KIND_ISSUE,
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
						Timestamp: 1649174827,
						TxHash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd7",
					},
					{
						Id:     8,
						Asset:  assets.AssetVCG,
						Amount: 4006,
						Kind:   bankingtypes.TRANSFER_KIND_DEFERRED,
						Extras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_EMAIL,
								Data: "retg@mail",
							},
							{
								Kind: extratypes.EXTRA_KIND_PHONE,
								Data: "+380685540",
							},
						},
						Timestamp: 1649174828,
						TxHash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd8",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			if err := r.UpdateBaseTransfers(tt.args.transfers...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateBaseTransfers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetBaseTransfers(t *testing.T) {
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
		want    []*bankingtypes.BaseTransfer
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				filter: filter.NewFilter().SetArgument(types.FieldTxHash, "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd7"),
			},
			want: []*bankingtypes.BaseTransfer{
				{
					Id:     7,
					Asset:  assets.AssetVCG,
					Amount: 4005,
					Kind:   bankingtypes.TRANSFER_KIND_ISSUE,
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
					Timestamp: 1649174827,
					TxHash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd7",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetBaseTransfers(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetBaseTransfers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetBaseTransfers() = %v, want %v", got, tt.want)
			}
		})
	}
}
