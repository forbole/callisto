package bank

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	d "git.ooo.ua/vipcoin/ovg-chain/x/domain"
	"github.com/brianvoe/gofakeit/v6"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"

	db "github.com/forbole/bdjuno/v4/_tests/database"
)

func TestRepository_InsertMsgMultiSend(t *testing.T) {
	type args struct {
		msg  []bank.MsgMultiSend
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertMsgMultiSend",
			args: args{
				msg: []bank.MsgMultiSend{
					{
						Inputs: []bank.Input{
							{
								Address: db.TestAddressCreator,
								Coins: sdk.NewCoins(
									sdk.NewCoin(d.DenomOVG, sdk.NewInt(5000_0000)),
									sdk.NewCoin(d.DenomSTOVG, sdk.NewInt(2_0000_0000)),
								),
							},
						},
						Outputs: []bank.Output{
							{
								Address: db.TestAddressCreator,
								Coins: sdk.NewCoins(
									sdk.NewCoin(d.DenomOVG, sdk.NewInt(2500_0000)),
									sdk.NewCoin(d.DenomSTOVG, sdk.NewInt(1_0000_0000)),
								),
							},
						},
					},
				},
				hash: gofakeit.LetterN(64),
			},
		},
	}
	/* EXAMPLE:
	INSERT INTO msg_multi_send (tx_hash, inputs, outputs) VALUES
	(
	    '326F78C19899ABBC0A472A28C00D23D3D809F3955114478F731C4F55CFD04216',
	    ARRAY[
	        ('ovg', '100000000')::COIN,
	        ('stovg', '100000000')::COIN
	    ],
	    ARRAY[
	        ('ovg', '100000000')::COIN,
	        ('stovg', '100000000')::COIN
	    ]
	);
	*/
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.Datastore.Bank.InsertMsgMultiSend(tt.args.hash, tt.args.msg...)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertMsgMultiSend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetAllMsgMultiSend(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAllMsgMultiSend",
			args: args{
				filter: filter.NewFilter(), // .SetArgument(types.FieldTxHash, "2683D5BAF21F8CE89613E5FF99DD8B0C8CCAF87BDC9EB1A9CD1456EE9F613EED"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := db.Datastore.Bank.GetAllMsgMultiSend(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMsgMultiSend() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, msg := range data {
				t.Logf("inputs: %v", msg.Inputs)
				t.Logf("outputs: %v", msg.Outputs)
			}
		})
	}
}
