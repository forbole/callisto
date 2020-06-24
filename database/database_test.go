package database_test

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	jconfig "github.com/desmos-labs/juno/config"
	"github.com/forbole/bdjuno/database"
	"github.com/stretchr/testify/suite"

	_ "github.com/proullon/ramsql/driver"
)

type DbTestSuite struct {
	suite.Suite

	database database.BigDipperDb
}

func (suite *DbTestSuite) SetupTest() {
	// Create the codec
	codec := simapp.MakeCodec()

	// Build the database
	config := jconfig.Config{
		DatabaseConfig: jconfig.DatabaseConfig{
			Type: "psql",
			Config: &jconfig.PostgreSQLConfig{
				Name:     "bdjuno",
				Host:     "localhost",
				Port:     5433,
				User:     "bdjuno",
				Password: "password",
			},
		},
	}

	db, err := database.Builder(config, codec)
	suite.Require().NoError(err)

	bigDipperDb, ok := (*db).(database.BigDipperDb)
	suite.Require().True(ok)

	// Delete the public schema
	_, err = bigDipperDb.Sql.Exec(`DROP SCHEMA public CASCADE;`)
	suite.Require().NoError(err)

	// Re-create the schema
	_, err = bigDipperDb.Sql.Exec(`CREATE SCHEMA public;`)
	suite.Require().NoError(err)

	dirPath := "../schema"
	dir, err := ioutil.ReadDir(dirPath)
	for _, fileInfo := range dir {
		file, err := ioutil.ReadFile(filepath.Join(dirPath, fileInfo.Name()))
		suite.Require().NoError(err)

		commentsRegExp := regexp.MustCompile(`/\*.*\*/`)
		requests := strings.Split(string(file), ";")
		for _, request := range requests {
			_, err := bigDipperDb.Sql.Exec(commentsRegExp.ReplaceAllString(request, ""))
			suite.Require().NoError(err)
		}
	}

	suite.database = bigDipperDb
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}
