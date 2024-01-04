package feeexcluder

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	"github.com/cosmos/cosmos-sdk/types"

	d "github.com/forbole/bdjuno/v4/_tests/database"
)

func TestRepository_InsertToGenesisState(t *testing.T) {
	type args struct {
		msg []fe.GenesisState
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertToGenesisState",
			args: args{
				msg: []fe.GenesisState{
					{
						Params: fe.Params{},
						AddressList: []fe.Address{
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
						AddressCount: 2,
						DailyStatsList: []fe.DailyStats{
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
						DailyStatsCount: 2,
						StatsList: []fe.Stats{
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
						TariffsList: []fe.Tariffs{
							{
								Denom: "ovg",
								Tariffs: []*fe.Tariff{
									{
										Id:            0,
										Amount:        "0",
										Denom:         "stovg",
										MinRefBalance: "10000000000",
										Fees: []*fe.Fees{
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
								Creator: d.TestAddressCreator,
							},
							{
								Denom: "ovg",
								Tariffs: []*fe.Tariff{
									{
										Id:            1,
										Amount:        "100000000000",
										Denom:         "stovg",
										MinRefBalance: "10000000000",
										Fees: []*fe.Fees{
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
										},
									},
								},
								Creator: d.TestAddressCreator,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range tt.args.msg {
				if err := d.Datastore.FeeExcluder.InsertToGenesisState(msg); (err != nil) != tt.wantErr {
					t.Errorf("InsertToGenesisState() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestRepository_GetAllGenesisState(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAllGenesisState",
			args: args{
				filter: filter.NewFilter(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := d.Datastore.FeeExcluder.GetAllGenesisState(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllGenesisState() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("size: %d", len(entity))
		})
	}
}

func TestRepository_DeleteGenesisState(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] DeleteGenesisState",
			args: args{
				id: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.FeeExcluder.DeleteGenesisState(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteGenesisState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
