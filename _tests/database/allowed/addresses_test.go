package allowed

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	allowed "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"

	d "github.com/forbole/bdjuno/v4/_tests/database"
)

func TestRepository_InsertToAddresses(t *testing.T) {
	type args struct {
		msg []allowed.Addresses
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertToAddresses",
			args: args{
				msg: []allowed.Addresses{
					{
						Address: []string{d.TestAddressCreator},
						Creator: d.TestAddressCreator,
					},
					{
						Address: []string{
							"ovg1ajeju2zaajdnj5j857sk7sjjlve65k2287jrfv",
							"ovg1n37hxqzfzu44ezyhhe3nzcym7u2ycxtf2w5mgw",
							"ovg1wvuy80m54dl8qw63u3jnaqjc3y82gnlk36gkjj",
						},
						Creator: d.TestAddressCreator,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.Allowed.InsertToAddresses(tt.args.msg...); (err != nil) != tt.wantErr {
				t.Errorf("InsertToAddresses() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetAllAddresses(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAddresses",
			args: args{
				filter: filter.NewFilter(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := d.Datastore.Allowed.GetAllAddresses(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllAddresses() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("size: %d", len(entity))
		})
	}
}

func TestRepository_UpdateAddresses(t *testing.T) {
	type args struct {
		msg []allowed.Addresses
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] UpdateAddresses",
			args: args{
				msg: []allowed.Addresses{
					{
						Id:      1,
						Address: []string{"ovg1wvuy80m54dl8qw63u3jnaqjc3y82gnlk36gkjj"},
						Creator: d.TestAddressCreator,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.Allowed.UpdateAddresses(tt.args.msg...); (err != nil) != tt.wantErr {
				t.Errorf("UpdateAddresses() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_DeleteAddressesByAddress(t *testing.T) {
	type args struct {
		addresses []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] DeleteAddressesByAddress",
			args: args{
				addresses: []string{
					"ovg1ajeju2zaajdnj5j857sk7sjjlve65k2287jrfv",
					"ovg1n37hxqzfzu44ezyhhe3nzcym7u2ycxtf2w5mgw",
					"ovg1wvuy80m54dl8qw63u3jnaqjc3y82gnlk36gkjj",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.Allowed.DeleteAddressesByAddress(tt.args.addresses...); (err != nil) != tt.wantErr {
				t.Errorf("DeleteAddressesByAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_DeleteAddressesByID(t *testing.T) {
	type args struct {
		ids []uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] DeleteAddressesByAddress",
			args: args{
				ids: []uint64{1, 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.Allowed.DeleteAddressesByID(tt.args.ids...); (err != nil) != tt.wantErr {
				t.Errorf("DeleteAddressesByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
