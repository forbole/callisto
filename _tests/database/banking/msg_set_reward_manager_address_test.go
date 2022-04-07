package banking

import (
	"reflect"
	"testing"

	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/cosmos/cosmos-sdk/simapp"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v2/database/types"
	bankingdb "github.com/forbole/bdjuno/v2/database/vipcoin/chain/banking"
)

func TestRepository_SaveMsgSetRewardManagerAddress(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		address []*bankingtypes.MsgSetRewardManagerAddress
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] save message into db",
			args: args{
				address: []*bankingtypes.MsgSetRewardManagerAddress{
					{
						Creator: "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv2g",
						Address: "vcg1k2gx4u0hwk87ja3wyakne8cl5gytnz0uc27xm4",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			if err := r.SaveMsgSetRewardMgrAddress(tt.args.address...); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveMsgSetRewardMgrAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetSaveMsgSetRewardManagerAddress(t *testing.T) {
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
		want    []*bankingtypes.MsgSetRewardManagerAddress
		wantErr bool
	}{
		{
			name: "[success] save message into db",
			args: args{
				filter: filter.NewFilter().SetArgument(types.FieldCreator, "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv2g"),
			},
			want: []*bankingtypes.MsgSetRewardManagerAddress{
				{
					Creator: "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv2g",
					Address: "vcg1k2gx4u0hwk87ja3wyakne8cl5gytnz0uc27xm4",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetMsgSetRewardMgrAddress(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetMsgSetRewardMgrAddress) error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetMsgSetRewardMgrAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
