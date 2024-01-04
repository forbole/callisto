package feeexcluder

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"

	d "github.com/forbole/bdjuno/v4/_tests/database"
)

func TestRepository_InsertToTariffs(t *testing.T) {
	type args struct {
		msg []fe.Tariffs
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertToTariffs",
			args: args{
				msg: []fe.Tariffs{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range tt.args.msg {
				if _, err := d.Datastore.FeeExcluder.InsertToTariffs(nil, msg); (err != nil) != tt.wantErr {
					t.Errorf("InsertToTariffs() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestRepository_GetAllTariffs(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAllTariffs",
			args: args{
				filter: filter.NewFilter(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := d.Datastore.FeeExcluder.GetAllTariffs(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllTariffs() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("size: %d", len(entity))
		})
	}
}

func TestRepository_GetAllTariffsDB(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAllTariffsDB",
			args: args{
				filter: filter.NewFilter(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := d.Datastore.FeeExcluder.GetAllTariffsDB(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllTariffsDB() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("size: %d", len(entity))
		})
	}
}

func TestRepository_UpdateTariffs(t *testing.T) {
	type args struct {
		msg fe.Tariffs
		id  uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] UpdateTariffs stovg -> ovg",
			args: args{
				msg: fe.Tariffs{
					Denom: "ovg",
					Tariffs: []*fe.Tariff{
						{
							Id:            1,
							Amount:        "100000000000",
							Denom:         "ovg",
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
						{
							Id:            0,
							Amount:        "0",
							Denom:         "ovg",
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
				id: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.FeeExcluder.UpdateTariffs(nil, tt.args.id, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("UpdateTariffs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_DeleteTariffs(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] DeleteTariffs 1",
			args: args{
				id: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.FeeExcluder.DeleteTariffs(nil, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteTariffs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
