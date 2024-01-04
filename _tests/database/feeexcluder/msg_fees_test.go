package feeexcluder

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	"github.com/brianvoe/gofakeit/v6"

	d "github.com/forbole/bdjuno/v4/_tests/database"
)

func TestRepository_InsertToFees(t *testing.T) {
	type args struct {
		msg []*fe.Fees
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertToFees",
			args: args{
				msg: []*fe.Fees{
					{
						AmountFrom:  "1000000",
						Fee:         "0.001",
						RefReward:   "0.25",
						StakeReward: "0.5",
						MinAmount:   1000,
						NoRefReward: true,
						Creator:     d.TestAddressCreator,
						Id:          0,
					},
					{
						AmountFrom:  "50000000",
						Fee:         "0.05",
						RefReward:   "0.25",
						StakeReward: "0.5",
						MinAmount:   100000,
						NoRefReward: true,
						Creator:     d.TestAddressCreator,
						Id:          1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range tt.args.msg {
				if _, err := d.Datastore.FeeExcluder.InsertToFees(nil, msg); (err != nil) != tt.wantErr {
					t.Errorf("InsertToFees() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestRepository_GetAllFees(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAllFees",
			args: args{
				filter: filter.NewFilter(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := d.Datastore.FeeExcluder.GetAllFees(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllFees() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("size: %d", len(entity))
		})
	}
}

func TestRepository_UpdateFees(t *testing.T) {
	type args struct {
		msg  []*fe.Fees
		id   uint64
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] UpdateFees",
			args: args{
				msg: []*fe.Fees{
					{
						AmountFrom:  "1000000",
						Fee:         "0.001",
						RefReward:   "0.25",
						StakeReward: "0.5",
						MinAmount:   1000,
						NoRefReward: false,
						Creator:     d.TestAddressCreator,
						Id:          0,
					},
					{
						AmountFrom:  "50000000",
						Fee:         "0.05",
						RefReward:   "0.25",
						StakeReward: "0.5",
						MinAmount:   100000,
						NoRefReward: false,
						Creator:     d.TestAddressCreator,
						Id:          1,
					},
				},
				id:   1,
				hash: gofakeit.LetterN(64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range tt.args.msg {
				if err := d.Datastore.FeeExcluder.UpdateFees(nil, tt.args.id, msg); (err != nil) != tt.wantErr {
					t.Errorf("UpdateFees() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestRepository_DeleteFees(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] DeleteFees",
			args: args{
				id: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.FeeExcluder.DeleteFees(nil, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteFees() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
