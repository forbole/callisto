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
	"github.com/cosmos/cosmos-sdk/simapp"
	accountsdb "github.com/forbole/bdjuno/v2/database/vipcoin/chain/accounts"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func TestRepository_SaveRegisterUser(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		msg []*accountstypes.MsgRegisterUser
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				msg: []*accountstypes.MsgRegisterUser{
					{
						Creator:         "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
						Address:         "vcg1hwxejcutmafuedd8trjqumfdkst2498pggx45q",
						Hash:            "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd1",
						PublicKey:       "4133425431324570546d614730316858526a302f6e7a6437726d3663526751755367626a694244566f4a486e",
						HolderWallet:    "vcg15hwngy8dys5l8kwqdyyuulhde39x6c6ad3wh0g",
						RefRewardWallet: "vcg12ndz92smw9pz34m7kfr5vqq6qscee7nv83vset",
						HolderWalletExtras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_WALLET,
								Data: "walletHolderName1",
							},
						},
						RefRewardWalletExtras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_WALLET,
								Data: "walletRefRewardName1",
							},
						},
						ReferrerHash: "",
					},
					{
						Creator:         "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv1g",
						Address:         "vcg1hwxejcutmafuedd8trjqumfdkst2498pggx42q",
						Hash:            "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd3",
						PublicKey:       "4133425431324570546d614730316858526a302f6e7a6437726d3663526751755367626a694244566f4a485e",
						HolderWallet:    "vcg15hwngy8dys5l8kwqdyyuulhde39x6c6ad3wh1g",
						RefRewardWallet: "vcg12ndz92smw9pz34m7kfr5vqq6qscee7nv85vset",
						HolderWalletExtras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_WALLET,
								Data: "walletHolderName2",
							},
						},
						RefRewardWalletExtras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_WALLET,
								Data: "walletRefRewardName2",
							},
						},
						ReferrerHash: "",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := accountsdb.NewRepository(db, codec.Marshaler)

			if err := r.SaveRegisterUser(tt.args.msg...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveRegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetRegisterUser(t *testing.T) {
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
		want    []*accountstypes.MsgRegisterUser
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				accfilter: filter.NewFilter().SetArgument("address", "vcg1hwxejcutmafuedd8trjqumfdkst2498pggx45q"),
			},
			want: []*accountstypes.MsgRegisterUser{
				{
					Creator:         "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
					Address:         "vcg1hwxejcutmafuedd8trjqumfdkst2498pggx45q",
					Hash:            "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd1",
					PublicKey:       "4133425431324570546d614730316858526a302f6e7a6437726d3663526751755367626a694244566f4a486e",
					HolderWallet:    "vcg15hwngy8dys5l8kwqdyyuulhde39x6c6ad3wh0g",
					RefRewardWallet: "vcg12ndz92smw9pz34m7kfr5vqq6qscee7nv83vset",
					HolderWalletExtras: []*extratypes.Extra{
						{
							Kind: extratypes.EXTRA_KIND_WALLET,
							Data: "walletHolderName1",
						},
					},
					RefRewardWalletExtras: []*extratypes.Extra{
						{
							Kind: extratypes.EXTRA_KIND_WALLET,
							Data: "walletRefRewardName1",
						},
					},
					ReferrerHash: "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := accountsdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetRegisterUser(tt.args.accfilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetRegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetRegisterUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
