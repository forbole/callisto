package allowed

import (
	"testing"

	allowed "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
	"github.com/brianvoe/gofakeit/v6"

	d "github.com/forbole/bdjuno/v4/_tests/database"
)

func TestRepository_InsertToUpdateAddresses(t *testing.T) {
	type args struct {
		msg  []*allowed.MsgUpdateAddresses
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertToUpdateByID",
			args: args{
				msg: []*allowed.MsgUpdateAddresses{
					{
						Creator: d.TestAddressCreator,
						Id:      1,
						Address: []string{"ovg18p9heyy3m4dsq7fj86p6v9yzzx8a64f86eec7u"},
					},
				},
				hash: gofakeit.LetterN(64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := d.Datastore.Allowed.InsertToUpdateAddresses(tt.args.hash, tt.args.msg...)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertToUpdateAddresses() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
