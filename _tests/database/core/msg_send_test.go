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

func TestRepository_InsertMsgSend(t *testing.T) {
	type args struct {
		msg  []*core.MsgSend
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
				msg: []*core.MsgSend{
					{
						Creator: d.TestAddressCreator,
						From:    d.TestAddressCreator,
						To:      d.TestAddress,
						Amount:  "100000000",
						Denom:   domain.DenomOVG,
					},
				},
				hash: gofakeit.LetterN(64),
			},
		},
		{
			name: "[success] InsertMsgSend (random address)",
			args: args{
				msg: []*core.MsgSend{
					{
						Creator: d.TestAddress,
						From:    d.TestAddress,
						To:      gofakeit.Regex("^ovg[a-z0-9]{39}"),
						Amount:  "50000000000",
						Denom:   domain.DenomOVG,
					},
				},
				hash: gofakeit.LetterN(64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := d.Datastore.Core.InsertMsgSend(tt.args.hash, tt.args.msg...)
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
			name: "[success] GetAllMsgSend",
			args: args{
				filter: filter.NewFilter(),
			},
		},
		{
			name: "[success] GetAllMsgSend by address from",
			args: args{
				filter: filter.NewFilter().SetArgument(types.FieldAddressFrom, d.TestAddressCreator),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := d.Datastore.Core.GetAllMsgSend(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMsgSend() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("size: %d", len(entity))
			for _, e := range entity {
				t.Logf("from: %s", e.From)
			}
		})
	}
}
