package feeexcluder

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	"github.com/brianvoe/gofakeit/v6"

	d "github.com/forbole/bdjuno/v4/_tests/database"
)

func TestRepository_InsertToMsgCreateAddress(t *testing.T) {
	type args struct {
		msg  []fe.MsgCreateAddress
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertToMsgCreateAddress",
			args: args{
				msg: []fe.MsgCreateAddress{
					{
						Address: d.TestAddressCreator,
						Creator: d.TestAddressCreator,
					},
					{
						Address: "ovg1wvuy80m54dl8qw63u3jnaqjc3y82gnlk36gkjj",
						Creator: d.TestAddressCreator,
					},
				},
				hash: gofakeit.LetterN(64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range tt.args.msg {
				if err := d.Datastore.FeeExcluder.InsertToMsgCreateAddress(tt.args.hash, msg); (err != nil) != tt.wantErr {
					t.Errorf("InsertToMsgCreateAddress() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestRepository_GetAllMsgCreateAddress(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAllMsgCreateAddress",
			args: args{
				filter: filter.NewFilter(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := d.Datastore.FeeExcluder.GetAllMsgCreateAddress(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMsgCreateAddress() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("size: %d", len(entity))
		})
	}
}

func TestRepository_UpdateMsgCreateAddress(t *testing.T) {
	type args struct {
		msg  []fe.MsgCreateAddress
		id   uint64
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] UpdateMsgCreateAddress",
			args: args{
				msg: []fe.MsgCreateAddress{
					{
						Address: d.TestAddressCreator,
						Creator: d.TestAddressCreator,
					},
				},
				id:   2,
				hash: gofakeit.LetterN(64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range tt.args.msg {
				if err := d.Datastore.FeeExcluder.UpdateMsgCreateAddress(tt.args.hash, tt.args.id, msg); (err != nil) != tt.wantErr {
					t.Errorf("UpdateMsgCreateAddress() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestRepository_DeleteMsgCreateAddress(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] DeleteMsgCreateAddress",
			args: args{
				id: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.FeeExcluder.DeleteMsgCreateAddress(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteMsgCreateAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
