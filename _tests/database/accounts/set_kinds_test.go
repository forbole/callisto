/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package accounts

import (
	"reflect"
	"testing"

	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/cosmos/cosmos-sdk/simapp"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	accountsdb "github.com/forbole/bdjuno/v3/database/overgold/chain/accounts"
	"github.com/forbole/bdjuno/v3/database/types"
)

func TestRepository_SaveKinds(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		msg []*accountstypes.MsgSetKinds
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				msg: []*accountstypes.MsgSetKinds{
					{
						Creator: "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
						Hash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd1",
						Kinds:   []accountstypes.AccountKind{accountstypes.ACCOUNT_KIND_SYSTEM},
					},
					{
						Creator: "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv1g",
						Hash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
						Kinds:   []accountstypes.AccountKind{accountstypes.ACCOUNT_KIND_SYSTEM},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := accountsdb.NewRepository(db, codec.Marshaler)

			if err := r.SaveKinds(tt.args.msg[0], ""); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveKinds() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetKinds(t *testing.T) {
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
		want    []*accountstypes.MsgSetKinds
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				accountFilter: filter.NewFilter().SetArgument(types.FieldCreator, "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g"),
			},
			want: []*accountstypes.MsgSetKinds{
				{
					Creator: "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
					Hash:    "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd1",
					Kinds:   []accountstypes.AccountKind{accountstypes.ACCOUNT_KIND_SYSTEM},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := accountsdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetKinds(tt.args.accountFilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetKinds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetKinds() = %v, want %v", got, tt.want)
			}
		})
	}
}
