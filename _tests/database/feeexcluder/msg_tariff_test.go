package feeexcluder

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	"github.com/brianvoe/gofakeit/v6"

	d "github.com/forbole/bdjuno/v4/_tests/database"
)

func TestRepository_InsertToTariff(t *testing.T) {
	type args struct {
		msg []*fe.Tariff
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertToTariff",
			args: args{
				msg: []*fe.Tariff{
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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range tt.args.msg {
				if _, err := d.Datastore.FeeExcluder.InsertToTariff(nil, msg); (err != nil) != tt.wantErr {
					t.Errorf("InsertToTariff() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestRepository_GetAllTariff(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAllTariff",
			args: args{
				filter: filter.NewFilter(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := d.Datastore.FeeExcluder.GetAllTariff(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllTariff() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("size: %d", len(entity))
		})
	}
}

func TestRepository_UpdateTariff(t *testing.T) {
	type args struct {
		msg  *fe.Tariff
		hash string
		id   uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] UpdateTariff",
			args: args{
				msg: &fe.Tariff{
					Id:            0,
					Amount:        "100000000",
					Denom:         "ovg",
					MinRefBalance: "100000000",
					Fees:          nil,
				},
				hash: gofakeit.LetterN(64),
				id:   1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.FeeExcluder.UpdateTariff(nil, tt.args.id, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("UpdateTariff() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_DeleteTariff(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] DeleteTariff 0",
			args: args{
				id: 0,
			},
		},
		{
			name: "[success] DeleteTariff 1",
			args: args{
				id: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.FeeExcluder.DeleteTariff(nil, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteTariff() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
