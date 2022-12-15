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

	"github.com/forbole/bdjuno/v2/database/types"
	accountsdb "github.com/forbole/bdjuno/v2/database/vipcoin/chain/accounts"
)

func TestRepository_SaveCreateAccount(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		msg []*accountstypes.MsgCreateAccount
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				[]*accountstypes.MsgCreateAccount{
					{
						Creator:   "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
						Hash:      "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd1",
						Address:   "ovg1hwxejcutmafuedd8trjqumfdkst2498pggx45q",
						PublicKey: "4133425431324570546d614730316858526a302f6e7a6437726d3663526751755367626a694244566f4a486e",
						Kinds:     []accountstypes.AccountKind{accountstypes.ACCOUNT_KIND_SYSTEM},
						State:     accountstypes.ACCOUNT_STATE_ACTIVE,
						Extras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_WALLET,
								Data: "walletHolderName3",
							},
							{
								Kind: extratypes.EXTRA_KIND_COMMENT,
								Data: "coment 123",
							},
						},
					},
					{
						Creator:   "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv1g",
						Hash:      "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd2",
						Address:   "ovg1hwxejcutmafuedd8trjqumfdkst2498pggx44q",
						PublicKey: "4133425431324570546d614730316858526a302f6e7a6437726d3663526751755367626a694244566f4a487e",
						Kinds:     []accountstypes.AccountKind{accountstypes.ACCOUNT_KIND_SYSTEM},
						State:     accountstypes.ACCOUNT_STATE_ACTIVE,
						Extras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_WALLET,
								Data: "walletHolderName4",
							},
							{
								Kind: extratypes.EXTRA_KIND_COMMENT,
								Data: "coment 125",
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

			if err := r.SaveCreateAccount(tt.args.msg[0], ""); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveCreateAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetCreateAccount(t *testing.T) {
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
		want    []*accountstypes.MsgCreateAccount
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				accfilter: filter.NewFilter().SetArgument(types.FieldAddress, "ovg1hwxejcutmafuedd8trjqumfdkst2498pggx45q"),
			},
			want: []*accountstypes.MsgCreateAccount{
				{
					Creator:   "ovg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv0g",
					Hash:      "a935ea2c467d7f666ea2a67870564f2efb902c05f0a2bb4b6202832aedd26cd1",
					Address:   "ovg1hwxejcutmafuedd8trjqumfdkst2498pggx45q",
					PublicKey: "4133425431324570546d614730316858526a302f6e7a6437726d3663526751755367626a694244566f4a486e",
					Kinds:     []accountstypes.AccountKind{accountstypes.ACCOUNT_KIND_SYSTEM},
					State:     accountstypes.ACCOUNT_STATE_ACTIVE,
					Extras: []*extratypes.Extra{
						{
							Kind: extratypes.EXTRA_KIND_WALLET,
							Data: "walletHolderName3",
						},
						{
							Kind: extratypes.EXTRA_KIND_COMMENT,
							Data: "coment 123",
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

			got, err := r.GetCreateAccount(tt.args.accfilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetCreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetCreateAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}
