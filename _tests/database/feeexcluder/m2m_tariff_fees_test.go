package feeexcluder

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/brianvoe/gofakeit/v6"

	d "github.com/forbole/bdjuno/v4/_tests/database"
	"github.com/forbole/bdjuno/v4/database/types"
)

// NOTE: add entity's in other tables before testing (tariff, fees)

func TestRepository_InsertToM2MTariffFees(t *testing.T) {
	type args struct {
		msg  []types.FeeExcluderM2MTariffFees
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertToM2MTariffFees",
			args: args{
				msg: []types.FeeExcluderM2MTariffFees{
					{
						TariffID: 0,
						FeesID:   0,
					},
					{
						TariffID: 0,
						FeesID:   1,
					},
				},
				hash: gofakeit.LetterN(64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.FeeExcluder.InsertToM2MTariffFees(nil, tt.args.msg...); (err != nil) != tt.wantErr {
				t.Errorf("InsertToM2MTariffFees() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetAllM2MTariffFees(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAllM2MTariffFees",
			args: args{
				filter: filter.NewFilter(),
			},
		},
		{
			name: "[success] GetAllM2MTariffFees by id",
			args: args{
				filter: filter.NewFilter().SetArgument(types.FieldTariffID, 0),
			},
		},
		{
			name: "[success] GetAllM2MTariffFees by id",
			args: args{
				filter: filter.NewFilter().SetArgument(types.FieldFeesID, 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := d.Datastore.FeeExcluder.GetAllM2MTariffFees(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllM2MTariffFees() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("size: %d", len(entity))
		})
	}
}

func TestRepository_DeleteM2MTariffFees(t *testing.T) {
	type args struct {
		msg types.FeeExcluderM2MTariffFees
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] DeleteM2MTariffFeesByTariff",
			args: args{
				msg: types.FeeExcluderM2MTariffFees{
					TariffID: 0,
					FeesID:   5,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.FeeExcluder.DeleteM2MTariffFeesByTariff(nil, tt.args.msg.TariffID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteM2MTariffFeesByTariff() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
