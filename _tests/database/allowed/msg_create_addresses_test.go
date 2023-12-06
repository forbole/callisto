package allowed

import (
	"testing"

	allowed "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
	"github.com/brianvoe/gofakeit/v6"

	d "github.com/forbole/bdjuno/v4/_tests/database"
)

func TestRepository_InsertToCreateAddresses(t *testing.T) {
	type args struct {
		msg  []*allowed.MsgCreateAddresses
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertToCreateAddresses single",
			args: args{
				msg: []*allowed.MsgCreateAddresses{
					{
						Address: []string{"ovg18p9heyy3m4dsq7fj86p6v9yzzx8a64f86eec7u"},
						Creator: d.TestAddressCreator,
					},
				},
				hash: gofakeit.LetterN(64),
			},
		},
		{
			name: "[success] InsertToCreateAddresses multiple",
			args: args{
				msg: []*allowed.MsgCreateAddresses{
					{
						Address: []string{
							"ovg1n37hxqzfzu44ezyhhe3nzcym7u2ycxtf2w5mgw",
							"ovg1wvuy80m54dl8qw63u3jnaqjc3y82gnlk36gkjj",
						},
						Creator: d.TestAddressCreator,
					},
					{
						Address: []string{
							"ovg1ajeju2zaajdnj5j857sk7sjjlve65k2287jrfv",
						},
						Creator: d.TestAddressCreator,
					},
				},
				hash: gofakeit.LetterN(64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := d.Datastore.Allowed.InsertToCreateAddresses(tt.args.hash, tt.args.msg...)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertToCreateAddresses() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
