package accounts

import (
	"reflect"
	"testing"

	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	extratypes "git.ooo.ua/vipcoin/chain/x/types"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/cosmos/cosmos-sdk/simapp"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	accountsdb "github.com/forbole/bdjuno/v3/database/overgold/chain/accounts"
	"github.com/forbole/bdjuno/v3/database/types"
)

func TestRepository_SaveAddAffiliate(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		msg []*accountstypes.MsgAddAffiliate
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{[]*accountstypes.MsgAddAffiliate{
				{
					Creator:         "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
					AccountHash:     "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd1",
					AffiliationHash: "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
					Affiliation:     accountstypes.AFFILIATION_KIND_REFERRAL,
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
				},
				{
					Creator:         "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv1g",
					AccountHash:     "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd3",
					AffiliationHash: "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd4",
					Affiliation:     accountstypes.AFFILIATION_KIND_REFERRAL,
					Extras: []*extratypes.Extra{
						{
							Kind: extratypes.EXTRA_KIND_EMAIL,
							Data: "rete@gmail",
						},
						{
							Kind: extratypes.EXTRA_KIND_PHONE,
							Data: "+380685456",
						},
					},
				},
			},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := accountsdb.NewRepository(db, codec.Marshaler)

			if err := r.SaveAddAffiliate(tt.args.msg[0], ""); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveAddAffiliate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetAddAffiliate(t *testing.T) {
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
		want    []*accountstypes.MsgAddAffiliate
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				accountFilter: filter.NewFilter().SetArgument(types.FieldCreator, "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv1g"),
			},
			want: []*accountstypes.MsgAddAffiliate{
				{
					Creator:         "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv1g",
					AccountHash:     "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd3",
					AffiliationHash: "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd4",
					Affiliation:     accountstypes.AFFILIATION_KIND_REFERRAL,
					Extras: []*extratypes.Extra{
						{
							Kind: extratypes.EXTRA_KIND_EMAIL,
							Data: "rete@gmail",
						},
						{
							Kind: extratypes.EXTRA_KIND_PHONE,
							Data: "+380685456",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := accountsdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetAddAffiliate(tt.args.accountFilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetAddAffiliate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetAddAffiliate() = %v, want %v", got, tt.want)
			}
		})
	}
}
