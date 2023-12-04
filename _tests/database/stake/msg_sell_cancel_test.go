package core

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	"git.ooo.ua/vipcoin/ovg-chain/x/domain"
	stake "git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	"github.com/brianvoe/gofakeit/v6"
	sdk "github.com/cosmos/cosmos-sdk/types"

	d "github.com/forbole/bdjuno/v4/_tests/database"
	"github.com/forbole/bdjuno/v4/database/types"
)

func TestRepository_InsertMsgSellCancel(t *testing.T) {
	type args struct {
		msg  []stake.MsgMsgCancelSell
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertMsgSellCancel",
			args: args{
				msg: []stake.MsgMsgCancelSell{
					{
						Creator: d.TestAddressCreator,
						Amount:  sdk.NewCoin(domain.DenomSTOVG, sdk.NewInt(100000000)),
					},
				},
				hash: gofakeit.LetterN(64),
			},
			wantErr: false,
		},
		{
			name: "[success] InsertMsgSellCancel (random hash)",
			args: args{
				msg: []stake.MsgMsgCancelSell{
					{
						Creator: d.TestAddressCreator,
						Amount:  sdk.NewCoin(domain.DenomSTOVG, sdk.NewInt(5000000000)),
					},
				},
				hash: gofakeit.LetterN(64),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := d.Datastore.Stake.InsertMsgSellCancel(tt.args.hash, tt.args.msg...)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertMsgSellCancel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetAllMsgSellCancel(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAllMsgSellCancel",
			args: args{
				filter: filter.NewFilter(),
			},
		},
		{
			name: "[success] GetAllMsgSellCancel by address",
			args: args{
				filter: filter.NewFilter().SetArgument(types.FieldAddress, d.TestAddressCreator),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := d.Datastore.Stake.GetAllMsgSellCancel(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMsgSellCancel() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("size: %d", len(entity))
			for _, e := range entity {
				t.Logf("creator: %s", e.Creator)
				t.Logf("amount: %s", e.Amount.Amount.String())
			}
		})
	}
}
