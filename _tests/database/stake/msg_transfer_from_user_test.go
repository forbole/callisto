package bank

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	"git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	"github.com/brianvoe/gofakeit/v6"

	db "github.com/forbole/bdjuno/v4/_tests/database"
)

func TestRepository_InsertMsgTransferFromUser(t *testing.T) {
	type args struct {
		msg  []types.MsgTransferFromUser
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertMsgTransferFromUser",
			args: args{
				msg: []types.MsgTransferFromUser{
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
			err := db.Datastore.Stake.InsertMsgTransferFromUser(tt.args.hash, tt.args.msg...)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertMsgTransferFromUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetAllMsgTransferFromUser(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAllMsgTransferFromUser",
			args: args{
				filter: filter.NewFilter().SetLimit(2),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := db.Datastore.Stake.GetAllMsgTransferFromUser(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMsgTransferFromUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			for idx, msg := range data {
				t.Logf("index: %d, message: %+v", idx, msg)
			}
		})
	}
}
