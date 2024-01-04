package feeexcluder

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	"github.com/cosmos/cosmos-sdk/types"

	d "github.com/forbole/bdjuno/v4/_tests/database"
)

func TestRepository_InsertToDailyStats(t *testing.T) {
	type args struct {
		msg []fe.DailyStats
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertToDailyStats",
			args: args{
				msg: []fe.DailyStats{
					{
						Id:            0,
						AmountWithFee: types.NewCoins(types.NewCoin("ovg", types.NewInt(1000))),
						AmountNoFee:   types.NewCoins(types.NewCoin("ovg", types.NewInt(0))),
						Fee:           types.NewCoins(types.NewCoin("ovg", types.NewInt(15))),
						CountWithFee:  1,
						CountNoFee:    0,
					},
					{
						Id:            0,
						AmountWithFee: types.NewCoins(types.NewCoin("ovg", types.NewInt(100000000))),
						AmountNoFee:   types.NewCoins(types.NewCoin("ovg", types.NewInt(20000))),
						Fee:           types.NewCoins(types.NewCoin("ovg", types.NewInt(1500))),
						CountWithFee:  10,
						CountNoFee:    2,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range tt.args.msg {
				if _, err := d.Datastore.FeeExcluder.InsertToDailyStats(nil, msg); (err != nil) != tt.wantErr {
					t.Errorf("InsertToDailyStats() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestRepository_GetAllDailyStats(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetDailyStats",
			args: args{
				filter: filter.NewFilter(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := d.Datastore.FeeExcluder.GetAllDailyStats(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllDailyStats() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("size: %d", len(entity))
		})
	}
}

func TestRepository_UpdateDailyStats(t *testing.T) {
	type args struct {
		msg []fe.DailyStats
		id  uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] UpdateDailyStats",
			args: args{
				msg: []fe.DailyStats{
					{
						Id:            1,
						AmountWithFee: types.NewCoins(types.NewCoin("ovg", types.NewInt(100000000))),
						AmountNoFee:   types.NewCoins(types.NewCoin("ovg", types.NewInt(20000))),
						Fee:           types.NewCoins(types.NewCoin("ovg", types.NewInt(1500))),
						CountWithFee:  100,
						CountNoFee:    20,
					},
				},
				id: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range tt.args.msg {
				if err := d.Datastore.FeeExcluder.UpdateDailyStats(nil, tt.args.id, msg); (err != nil) != tt.wantErr {
					t.Errorf("UpdateDailyStats() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestRepository_DeleteDailyStats(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] DeleteDailyStats",
			args: args{
				id: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.FeeExcluder.DeleteDailyStats(nil, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteDailyStats() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
