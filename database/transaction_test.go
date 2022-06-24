package database

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/forbole/juno/v2/database/postgresql"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func TestDb_GetTransaction(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				filter: filter.NewFilter(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Db{
				Database: &postgresql.Database{EncodingConfig: &params.EncodingConfig{}},
				Sqlx:     db,
			}

			db.EncodingConfig.Marshaler = codec.Marshaler

			_, err := db.GetTransaction(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Db.GetTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
