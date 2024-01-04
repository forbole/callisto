package bank

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	"git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	"github.com/brianvoe/gofakeit/v6"

	db "github.com/forbole/bdjuno/v4/_tests/database"
)

func TestRepository_InsertMsgTransferToUser(t *testing.T) {
	type args struct {
		msg  []types.MsgTransferToUser
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertMsgTransferToUser",
			args: args{
				msg: []types.MsgTransferToUser{
					{
						Creator: db.TestAddressCreator,
						Amount:  "100000000", // 1 OVG
						Address: gofakeit.Regex(`ovg1[a-z0-9]{38}$`),
					},
				},
				hash: gofakeit.LetterN(64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.Datastore.Stake.InsertMsgTransferToUser(tt.args.hash, tt.args.msg...)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertMsgTransferToUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetAllMsgTransferToUser(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAllMsgTransferToUser",
			args: args{
				filter: filter.NewFilter().SetLimit(2),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := db.Datastore.Stake.GetAllMsgTransferToUser(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMsgTransferToUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			for idx, msg := range data {
				t.Logf("index: %d, message: %+v", idx, msg)
			}
		})
	}
}
