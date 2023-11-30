package bank

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	d "git.ooo.ua/vipcoin/ovg-chain/x/domain"
	"github.com/brianvoe/gofakeit/v6"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"

	db "github.com/forbole/bdjuno/v4/_tests/database"
	"github.com/forbole/bdjuno/v4/database/types"
)

func TestRepository_InsertMsgSend(t *testing.T) {
	type args struct {
		msg  []bank.MsgSend
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertMsgSend",
			args: args{
				msg: []bank.MsgSend{
					{
						FromAddress: db.TestAddressCreator,
						ToAddress:   db.TestAddressCreator,
						Amount: sdk.NewCoins(
							sdk.NewCoin(d.DenomOVG, sdk.NewInt(1_00000000)),
							sdk.NewCoin(d.DenomSTOVG, sdk.NewInt(50_00000000)),
						),
					},
				},
				hash: gofakeit.LetterN(64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.Datastore.Bank.InsertMsgSend(tt.args.hash, tt.args.msg...)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertMsgSend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetAllMsgSend(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAllMsgSend by to_address",
			args: args{
				filter: filter.NewFilter().SetArgument(types.FieldToAddress, db.TestAddressCreator),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := db.Datastore.Bank.GetAllMsgSend(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMsgSend() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, msg := range data {
				t.Logf("%s: %s", msg.ToAddress, msg.Amount.String())
			}
		})
	}
}
