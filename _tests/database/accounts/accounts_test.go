/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package accounts

import (
	"reflect"
	"testing"

	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	extratypes "git.ooo.ua/vipcoin/chain/x/types"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/brianvoe/gofakeit/v6"
	anytype "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v2/database/types"
	accountsdb "github.com/forbole/bdjuno/v2/database/vipcoin/chain/accounts"
)

func TestRepository_SaveAccounts(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		accounts []*accountstypes.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				accounts: []*accountstypes.Account{
					{
						Hash:    "b8b6cb7629d68b3ecf9ce200f631ffc72232bc798a7db755307332a40add5e37",
						Address: "ovg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdy",
						PublicKey: &anytype.Any{
							TypeUrl: "/cosmos.crypto.secp256k1.PubKey",
							Value:   []uint8{10, 33, 2, 32, 174, 170, 220, 129, 199, 203, 202, 84, 205, 96, 6, 247, 144, 46, 61, 225, 73, 220, 82, 19, 53, 39, 205, 55, 45, 114, 65, 148, 77, 198, 60},
						},
						Kinds:      []accountstypes.AccountKind{accountstypes.ACCOUNT_KIND_SYSTEM},
						State:      accountstypes.ACCOUNT_STATE_ACTIVE,
						Extras:     []*extratypes.Extra{},
						Affiliates: []*accountstypes.Affiliate{},
						Wallets:    []string{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid",
			args: args{
				accounts: []*accountstypes.Account{
					{
						Hash:    gofakeit.Regex("^0x[a-fA-F0-9]{40}$"),
						Address: gofakeit.Regex("^0x[a-fA-F0-9]{40}$"),
						PublicKey: &anytype.Any{
							TypeUrl: "/cosmos.crypto.secp256k1.PubKey",
							Value:   []uint8{10, 33, 2, 32, 174, 170, 220, 129, 199, 203, 202, 84, 205, 96, 6, 247, 144, 46, 61, 225, 73, 220, 82, 19, 53, 39, 205, 55, 45, 114, 65, 148, 77, 198, 60},
						},
						Kinds: []accountstypes.AccountKind{accountstypes.ACCOUNT_KIND_GENERAL, accountstypes.ACCOUNT_KIND_MERCHANT},
						State: accountstypes.ACCOUNT_STATE_ACTIVE,
						Extras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_COMMENT,
								Data: "some text",
							},
						},
						Affiliates: []*accountstypes.Affiliate{
							{
								Address:     gofakeit.Regex("^0x[a-fA-F0-9]{40}$"),
								Affiliation: accountstypes.AFFILIATION_KIND_REFERRAL,
								Extras: []*extratypes.Extra{
									{
										Kind: extratypes.EXTRA_KIND_EMAIL,
										Data: gofakeit.Email(),
									},
									{
										Kind: extratypes.EXTRA_KIND_PHONE,
										Data: gofakeit.Phone(),
									},
								},
							},
							{
								Address:     gofakeit.Regex("^0x[a-fA-F0-9]{40}$"),
								Affiliation: accountstypes.AFFILIATION_KIND_REFERRAL,
								Extras: []*extratypes.Extra{
									{
										Kind: extratypes.EXTRA_KIND_EMAIL,
										Data: gofakeit.Email(),
									},
									{
										Kind: extratypes.EXTRA_KIND_PHONE,
										Data: gofakeit.Phone(),
									},
								},
							},
						},
						Wallets: []string{gofakeit.Regex("^0x[a-fA-F0-9]{40}$")},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := accountsdb.NewRepository(db, codec.Marshaler)

			if err := r.SaveAccounts(tt.args.accounts...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveAccounts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_UpdateAccounts(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		accounts []*accountstypes.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				accounts: []*accountstypes.Account{
					{
						Hash:    "b8b6cb7629d68b3ecf9ce200f631ffc72232bc798a7db755307332a40add5e37",
						Address: "ovg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdz",
						PublicKey: &anytype.Any{
							TypeUrl: "/cosmos.crypto.secp256k1.PubKey",
							Value:   []uint8{10, 33, 2, 32, 174, 170, 220, 129, 199, 203, 202, 84, 205, 96, 6, 247, 144, 46, 61, 225, 73, 220, 82, 19, 53, 39, 205, 55, 45, 114, 65, 148, 77, 198, 60},
						},
						Kinds: []accountstypes.AccountKind{accountstypes.ACCOUNT_KIND_GENERAL, accountstypes.ACCOUNT_KIND_MERCHANT},
						State: accountstypes.ACCOUNT_STATE_ACTIVE,
						Extras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_COMMENT,
								Data: "some text",
							},
						},
						Affiliates: []*accountstypes.Affiliate{
							{
								Address:     gofakeit.Regex("^0x[a-fA-F0-9]{40}$"),
								Affiliation: accountstypes.AFFILIATION_KIND_REFERRAL,
								Extras: []*extratypes.Extra{
									{
										Kind: extratypes.EXTRA_KIND_EMAIL,
										Data: gofakeit.Email(),
									},
									{
										Kind: extratypes.EXTRA_KIND_PHONE,
										Data: gofakeit.Phone(),
									},
								},
							},
							{
								Address:     gofakeit.Regex("^0x[a-fA-F0-9]{40}$"),
								Affiliation: accountstypes.AFFILIATION_KIND_REFERRAL,
								Extras: []*extratypes.Extra{
									{
										Kind: extratypes.EXTRA_KIND_EMAIL,
										Data: gofakeit.Email(),
									},
									{
										Kind: extratypes.EXTRA_KIND_PHONE,
										Data: gofakeit.Phone(),
									},
								},
							},
						},
						Wallets: []string{gofakeit.Regex("^0x[a-fA-F0-9]{40}$")},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid",
			args: args{
				accounts: []*accountstypes.Account{
					{
						Hash:    "b8b6cb7629d68b3ecf9ce200f631ffc72232bc798a7db755307332a40add5e37",
						Address: "ovg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdz",
						PublicKey: &anytype.Any{
							TypeUrl: "/cosmos.crypto.secp256k1.PubKey",
							Value:   []uint8{10, 33, 2, 32, 174, 170, 220, 129, 199, 203, 202, 84, 205, 96, 6, 247, 144, 46, 61, 225, 73, 220, 82, 19, 53, 39, 205, 55, 45, 114, 65, 148, 77, 198, 60},
						},
						Kinds:      []accountstypes.AccountKind{accountstypes.ACCOUNT_KIND_SYSTEM},
						State:      accountstypes.ACCOUNT_STATE_ACTIVE,
						Extras:     []*extratypes.Extra{},
						Affiliates: []*accountstypes.Affiliate{},
						Wallets:    []string{},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := accountsdb.NewRepository(db, codec.Marshaler)

			if err := r.UpdateAccounts(tt.args.accounts...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateAccounts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetAccounts(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		accountFilter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		want    []*accountstypes.Account
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				accountFilter: filter.NewFilter().SetArgument(types.FieldAddress, "ovg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdy"),
			},
			want: []*accountstypes.Account{
				{
					Hash:    "b8b6cb7629d68b3ecf9ce200f631ffc72232bc798a7db755307332a40add5e37",
					Address: "ovg1qq995wzw6zgqjm8g2twsykl7xsj0apxmtuxfdy",
					PublicKey: &anytype.Any{
						TypeUrl: "/cosmos.crypto.secp256k1.PubKey",
						Value:   []uint8{10, 33, 2, 32, 174, 170, 220, 129, 199, 203, 202, 84, 205, 96, 6, 247, 144, 46, 61, 225, 73, 220, 82, 19, 53, 39, 205, 55, 45, 114, 65, 148, 77, 198, 60},
					},
					Kinds:      []accountstypes.AccountKind{accountstypes.ACCOUNT_KIND_SYSTEM},
					State:      accountstypes.ACCOUNT_STATE_ACTIVE,
					Extras:     []*extratypes.Extra{},
					Affiliates: []*accountstypes.Affiliate{},
					Wallets:    []string{},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := accountsdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetAccounts(tt.args.accountFilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetAccounts() = %v, want %v", got, tt.want)
			}
		})
	}
}
