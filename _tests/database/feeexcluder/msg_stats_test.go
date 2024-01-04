package feeexcluder

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	"github.com/cosmos/cosmos-sdk/types"

	d "github.com/forbole/bdjuno/v4/_tests/database"
)

func TestRepository_InsertToStats(t *testing.T) {
	type args struct {
		msg []fe.Stats
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertToStats",
			args: args{
				msg: []fe.Stats{
					{
						Index: "1",
						Date:  "2023-12-19",
						Stats: &fe.DailyStats{
							Id:            0,
							AmountWithFee: nil,
							AmountNoFee:   types.NewCoins(types.NewCoin("ovg", types.NewInt(1000000196105000))),
							Fee:           nil,
							CountWithFee:  0,
							CountNoFee:    10,
						},
					},
					{
						Index: "2",
						Date:  "2023-12-20",
						Stats: &fe.DailyStats{
							Id:            0,
							AmountWithFee: types.NewCoins(types.NewCoin("ovg", types.NewInt(18944001000))),
							AmountNoFee:   nil,
							Fee:           types.NewCoins(types.NewCoin("ovg", types.NewInt(174440010))),
							CountWithFee:  5,
							CountNoFee:    0,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, stats := range tt.args.msg {
				if _, err := d.Datastore.FeeExcluder.InsertToStats(nil, stats); (err != nil) != tt.wantErr {
					t.Errorf("InsertToStats() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestRepository_GetAllStats(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetStats",
			args: args{
				filter: filter.NewFilter(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := d.Datastore.FeeExcluder.GetAllStats(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllStats() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("size: %d", len(entity))
		})
	}
}

func TestRepository_UpdateStats(t *testing.T) {
	type args struct {
		msg []fe.Stats
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] UpdateStats",
			args: args{
				msg: []fe.Stats{
					{
						Index: "1",
						Date:  "2023-01-20",
						Stats: &fe.DailyStats{
							Id:            0,
							AmountWithFee: types.NewCoins(types.NewCoin("stovg", types.NewInt(18944001000))),
							AmountNoFee:   nil,
							Fee:           types.NewCoins(types.NewCoin("stovg", types.NewInt(174440010))),
							CountWithFee:  5,
							CountNoFee:    0,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range tt.args.msg {
				if err := d.Datastore.FeeExcluder.UpdateStats(nil, msg); (err != nil) != tt.wantErr {
					t.Errorf("UpdateStats() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestRepository_DeleteStats(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] DeleteStats",
			args: args{
				id: "2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.FeeExcluder.DeleteStats(nil, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteStats() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
