package core

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	stake "git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	"github.com/brianvoe/gofakeit/v6"

	d "github.com/forbole/bdjuno/v4/_tests/database"
	"github.com/forbole/bdjuno/v4/database/types"
)

func TestRepository_InsertMsgSell(t *testing.T) {
	type args struct {
		msg  []stake.MsgSellRequest
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertMsgSell",
			args: args{
				msg: []stake.MsgSellRequest{
					{
						Creator: d.TestAddressCreator,
						Amount:  "100000000",
					},
				},
				hash: gofakeit.LetterN(64),
			},
			wantErr: false,
		},
		{
			name: "[success] InsertMsgSell (random hash)",
			args: args{
				msg: []stake.MsgSellRequest{
					{
						Creator: d.TestAddressCreator,
						Amount:  "50000000000",
					},
				},
				hash: gofakeit.LetterN(64),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := d.Datastore.Stake.InsertMsgSell(tt.args.hash, tt.args.msg...)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertMsgSell() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetAllMsgSell(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAllMsgSell",
			args: args{
				filter: filter.NewFilter(),
			},
		},
		{
			name: "[success] GetAllMsgSell by address",
			args: args{
				filter: filter.NewFilter().SetArgument(types.FieldAddress, d.TestAddressCreator),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := d.Datastore.Stake.GetAllMsgSell(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMsgSell() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("size: %d", len(entity))
			for _, e := range entity {
				t.Logf("creator: %s", e.Creator)
				t.Logf("amount: %s", e.Amount)
			}
		})
	}
}
