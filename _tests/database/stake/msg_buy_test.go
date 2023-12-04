package core

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	stake "git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	"github.com/brianvoe/gofakeit/v6"

	d "github.com/forbole/bdjuno/v4/_tests/database"
	"github.com/forbole/bdjuno/v4/database/types"
)

func TestRepository_InsertMsgBuy(t *testing.T) {
	type args struct {
		msg  []stake.MsgBuyRequest
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertMsgBuy",
			args: args{
				msg: []stake.MsgBuyRequest{
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
			name: "[success] InsertMsgBuy (random hash)",
			args: args{
				msg: []stake.MsgBuyRequest{
					{
						Creator: d.TestAddressCreator,
						Amount:  "5000000000",
					},
				},
				hash: gofakeit.LetterN(64),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := d.Datastore.Stake.InsertMsgBuy(tt.args.hash, tt.args.msg...)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertMsgBuy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetAllMsgBuy(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAllMsgBuy",
			args: args{
				filter: filter.NewFilter(),
			},
		},
		{
			name: "[success] GetAllMsgBuy by address",
			args: args{
				filter: filter.NewFilter().SetArgument(types.FieldAddress, d.TestAddressCreator),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := d.Datastore.Stake.GetAllMsgBuy(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMsgBuy() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("size: %d", len(entity))
			for _, e := range entity {
				t.Logf("creator: %s", e.Creator)
				t.Logf("amount: %s", e.Amount)
			}
		})
	}
}
