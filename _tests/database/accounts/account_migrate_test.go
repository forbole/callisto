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

	"github.com/forbole/bdjuno/v2/database/types"
	accountsdb "github.com/forbole/bdjuno/v2/database/vipcoin/chain/accounts"
)

func TestRepository_SaveAccountMigrate(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		msg []*accountstypes.MsgAccountMigrate
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				[]*accountstypes.MsgAccountMigrate{
					{
						Creator:   "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
						Address:   "vcg1hwxejcutmafuedd8trjqumfdkst2498pggx45q",
						Hash:      "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd1",
						PublicKey: "4133425431324570546d614730316858526a302f6e7a6437726d3663526751755367626a694244566f4a486e",
					},
					{
						Creator:   "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv1g",
						Address:   "vcg1hwxejcutmafuedd8trjqumfdkst2498pggx46q",
						Hash:      "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
						PublicKey: "4133425431324570546d614730316858526a302f6e7a6437726d3663526751755367626a694244566f4a484e",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := accountsdb.NewRepository(db, codec.Marshaler)

			if err := r.SaveAccountMigrate(tt.args.msg...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveAccountMigrate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetAccountMigrate(t *testing.T) {
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
		want    []*accountstypes.MsgAccountMigrate
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				accountFilter: filter.NewFilter().SetArgument(types.FieldAddress, "vcg1hwxejcutmafuedd8trjqumfdkst2498pggx45q"),
			},
			want: []*accountstypes.MsgAccountMigrate{
				{
					Creator:   "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
					Address:   "vcg1hwxejcutmafuedd8trjqumfdkst2498pggx45q",
					Hash:      "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd1",
					PublicKey: "4133425431324570546d614730316858526a302f6e7a6437726d3663526751755367626a694244566f4a486e",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := accountsdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetAccountMigrate(tt.args.accountFilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetAccountMigrate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetAccountMigrate() = %v, want %v", got, tt.want)
			}
		})
	}
}
