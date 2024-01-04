package feeexcluder

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"

	d "github.com/forbole/bdjuno/v4/_tests/database"
)

func TestRepository_InsertToAddress(t *testing.T) {
	type args struct {
		msg []fe.Address
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertToAddress",
			args: args{
				msg: []fe.Address{
					{
						Id:      1,
						Address: d.TestAddressCreator,
						Creator: d.TestAddressCreator,
					},
					{
						Id:      2,
						Address: "ovg1wvuy80m54dl8qw63u3jnaqjc3y82gnlk36gkjj",
						Creator: d.TestAddressCreator,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range tt.args.msg {
				if _, err := d.Datastore.FeeExcluder.InsertToAddress(nil, msg); (err != nil) != tt.wantErr {
					t.Errorf("InsertToAddress() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestRepository_GetAllAddress(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAddress",
			args: args{
				filter: filter.NewFilter(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := d.Datastore.FeeExcluder.GetAllAddress(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllAddress() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("size: %d", len(entity))
		})
	}
}

func TestRepository_UpdateAddress(t *testing.T) {
	type args struct {
		msg []fe.Address
		id  uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] UpdateAddress",
			args: args{
				msg: []fe.Address{
					{
						Id:      3,
						Address: "ovg1wvuy80m54dl8qw63u3jnaqjc3y82gnlk36gkjj",
						Creator: d.TestAddressCreator,
					},
				},
				id: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range tt.args.msg {
				if err := d.Datastore.FeeExcluder.UpdateAddress(nil, tt.args.id, msg); (err != nil) != tt.wantErr {
					t.Errorf("UpdateAddress() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestRepository_DeleteAddress(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] DeleteAddress",
			args: args{
				id: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.FeeExcluder.DeleteAddress(nil, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
