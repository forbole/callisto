package allowed

import (
	"testing"

	allowed "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
	"github.com/brianvoe/gofakeit/v6"

	d "github.com/forbole/bdjuno/v4/_tests/database"
)

func TestRepository_InsertToDeleteByID(t *testing.T) {
	type args struct {
		msg  []*allowed.MsgDeleteByID
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertToDeleteByID",
			args: args{
				msg: []*allowed.MsgDeleteByID{
					{
						Creator: d.TestAddressCreator,
						Id:      1,
					},
				},
				hash: gofakeit.LetterN(64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.Allowed.InsertToDeleteByID(tt.args.hash, tt.args.msg...); (err != nil) != tt.wantErr {
				t.Errorf("InsertToDeleteByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
