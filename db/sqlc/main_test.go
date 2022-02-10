package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/danglebary/beatstore-backend-go/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Failed to load config settings", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Failed to connect to db", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
