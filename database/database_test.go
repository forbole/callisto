package database_test

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	dbconfig "github.com/forbole/juno/v5/database/config"
	"github.com/forbole/juno/v5/logging"

	junodb "github.com/forbole/juno/v5/database"

	"github.com/forbole/bdjuno/v4/database"

	juno "github.com/forbole/juno/v5/types"

	tmversion "github.com/cometbft/cometbft/proto/tendermint/version"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/juno/v5/types/params"

	"github.com/stretchr/testify/suite"

	_ "github.com/proullon/ramsql/driver"
)

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}

type DbTestSuite struct {
	suite.Suite

	database *database.Db
}

func (suite *DbTestSuite) SetupTest() {
	// Create the codec
	codec := params.MakeTestEncodingConfig()

	// Build the database
	dbCfg := dbconfig.NewDatabaseConfig(
		"postgresql://user:password@localhost:5432/database-name?sslmode=disable&search_path=public",
		"postgresql://user:password@localhost:5432/provider-database-name?sslmode=disable&search_path=public",
		1,
		1,
		100000,
		1000,
		"",
		"",
		"",
		"",
	)
	db, err := database.Builder(junodb.NewContext(dbCfg, codec, logging.DefaultLogger()))
	suite.Require().NoError(err)

	bigDipperDb, ok := (db).(*database.Db)
	suite.Require().True(ok)

	// Delete the public schema
	_, err = bigDipperDb.SQL.Exec(`DROP SCHEMA public CASCADE;`)
	suite.Require().NoError(err)

	// Re-create the schema
	_, err = bigDipperDb.SQL.Exec(`CREATE SCHEMA public;`)
	suite.Require().NoError(err)

	dirPath := path.Join(".", "schema")
	dir, err := os.ReadDir(dirPath)
	suite.Require().NoError(err)

	for _, fileInfo := range dir {
		file, err := os.ReadFile(filepath.Join(dirPath, fileInfo.Name()))
		suite.Require().NoError(err)

		commentsRegExp := regexp.MustCompile(`/\*.*\*/`)
		requests := strings.Split(string(file), ";")
		for _, request := range requests {
			_, err := bigDipperDb.SQL.Exec(commentsRegExp.ReplaceAllString(request, ""))
			suite.Require().NoError(err)
		}
	}

	suite.database = bigDipperDb
}

// getBlock builds, stores and returns a block for the provided height
func (suite *DbTestSuite) getBlock(height int64) *juno.Block {

	addr, err := sdk.ConsAddressFromBech32("cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl")
	suite.Require().NoError(err)

	tmBlock := &tmctypes.ResultBlock{
		BlockID: tmtypes.BlockID{},
		Block: &tmtypes.Block{
			Header: tmtypes.Header{
				Version:            tmversion.Consensus{},
				ChainID:            "",
				Height:             height,
				Time:               time.Now(),
				LastBlockID:        tmtypes.BlockID{},
				LastCommitHash:     nil,
				DataHash:           nil,
				ValidatorsHash:     []byte("hash"),
				NextValidatorsHash: nil,
				ConsensusHash:      nil,
				AppHash:            nil,
				LastResultsHash:    nil,
				EvidenceHash:       nil,
				ProposerAddress:    tmtypes.Address(addr.Bytes()),
			},
			Data:     tmtypes.Data{},
			Evidence: tmtypes.EvidenceData{},
			LastCommit: &tmtypes.Commit{
				Height:     height - 1,
				Round:      1,
				BlockID:    tmtypes.BlockID{},
				Signatures: nil,
			},
		},
	}

	block := juno.NewBlockFromTmBlock(tmBlock, 10000)
	err = suite.database.SaveBlock(block)
	suite.Require().NoError(err)
	return block
}

// getAccount saves inside the database an account having the given address
func (suite *DbTestSuite) getAccount(addr string) sdk.AccAddress {
	delegator, err := sdk.AccAddressFromBech32(addr)
	suite.Require().NoError(err)

	_, err = suite.database.SQL.Exec(`INSERT INTO account (address) VALUES ($1) ON CONFLICT DO NOTHING`, delegator.String())
	suite.Require().NoError(err)

	return delegator
}
