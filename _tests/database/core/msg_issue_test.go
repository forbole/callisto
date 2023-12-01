package core

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	core "git.ooo.ua/vipcoin/ovg-chain/x/core/types"
	"git.ooo.ua/vipcoin/ovg-chain/x/domain"
	"github.com/brianvoe/gofakeit/v6"

	d "github.com/forbole/bdjuno/v4/_tests/database"
	"github.com/forbole/bdjuno/v4/database/types"
)

func TestRepository_InsertMsgIssue(t *testing.T) {
	type args struct {
		msg  []*core.MsgIssue
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertMsgIssue",
			args: args{
				msg: []*core.MsgIssue{
					{
						Creator: d.TestAddressCreator,
						Amount:  "100000000",
						Denom:   domain.DenomOVG,
						Address: d.TestAddress,
					},
				},
				hash: gofakeit.LetterN(64),
			},
		},
		{
			name: "[success] InsertMsgIssue (random address)",
			args: args{
				msg: []*core.MsgIssue{
					{
						Creator: d.TestAddressCreator,
						Amount:  "50000000000",
						Denom:   domain.DenomOVG,
						Address: gofakeit.Regex("^ovg[a-z0-9]{39}"),
					},
				},
				hash: gofakeit.LetterN(64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := d.Datastore.Core.InsertMsgIssue(tt.args.hash, tt.args.msg...)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertMsgIssue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetAllMsgIssue(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAllMsgIssue",
			args: args{
				filter: filter.NewFilter(),
			},
		},
		{
			name: "[success] GetAllMsgIssue by address",
			args: args{
				filter: filter.NewFilter().SetArgument(types.FieldAddress, d.TestAddress),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := d.Datastore.Core.GetAllMsgIssue(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMsgIssue() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("size: %d", len(entity))
			for _, e := range entity {
				t.Logf("address: %s", e.Address)
			}
		})
	}
}
