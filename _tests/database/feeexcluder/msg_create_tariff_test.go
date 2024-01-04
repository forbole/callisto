package feeexcluder

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	"github.com/brianvoe/gofakeit/v6"

	d "github.com/forbole/bdjuno/v4/_tests/database"
)

func TestRepository_InsertToMsgCreateTariffs(t *testing.T) {
	type args struct {
		msg  []fe.MsgCreateTariffs
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertToMsgCreateTariffs",
			args: args{
				msg: []fe.MsgCreateTariffs{
					{
						Creator: d.TestAddressCreator,
						Denom:   "ovg",
						Tariff: &fe.Tariff{
							Id:            0,
							Amount:        "1",
							Denom:         "stovg",
							MinRefBalance: "10000000000",
							Fees: []*fe.Fees{
								{
									AmountFrom:  "0",
									Fee:         "0.01",
									RefReward:   "0.25",
									StakeReward: "0.5",
									MinAmount:   1000,
									NoRefReward: false,
									Creator:     "ovg1dcftms3rgxvsa2pffedke7jz5np8k4lzp6pet9",
									Id:          0,
								},
								{
									AmountFrom:  "10000000000",
									Fee:         "0.009",
									RefReward:   "0.25",
									StakeReward: "0.5",
									MinAmount:   1000,
									NoRefReward: false,
									Creator:     "ovg1dcftms3rgxvsa2pffedke7jz5np8k4lzp6pet9",
									Id:          1,
								},
							},
						},
					},
					{
						Creator: d.TestAddressCreator,
						Denom:   "ovg",
						Tariff: &fe.Tariff{
							Id:            1,
							Amount:        "100000000000",
							Denom:         "stovg",
							MinRefBalance: "10000000000",
							Fees: []*fe.Fees{
								{
									AmountFrom:  "1000000000000",
									Fee:         "0.002",
									RefReward:   "0.25",
									StakeReward: "0.5",
									MinAmount:   1000,
									NoRefReward: false,
									Creator:     "ovg1dcftms3rgxvsa2pffedke7jz5np8k4lzp6pet9",
									Id:          2,
								},
								{
									AmountFrom:  "10000000000000",
									Fee:         "0.001",
									RefReward:   "0.25",
									StakeReward: "0.5",
									MinAmount:   1000,
									NoRefReward: false,
									Creator:     "ovg1dcftms3rgxvsa2pffedke7jz5np8k4lzp6pet9",
									Id:          3,
								},
							},
						},
					},
				},
				hash: gofakeit.LetterN(64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range tt.args.msg {
				if err := d.Datastore.FeeExcluder.InsertToMsgCreateTariffs(tt.args.hash, msg); (err != nil) != tt.wantErr {
					t.Errorf("InsertToMsgCreateTariffs() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestRepository_GetAllMsgCreateTariffs(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAllMsgCreateTariffs",
			args: args{
				filter: filter.NewFilter(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := d.Datastore.FeeExcluder.GetAllMsgCreateTariffs(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMsgCreateTariffs() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("size: %d", len(entity))
		})
	}
}
