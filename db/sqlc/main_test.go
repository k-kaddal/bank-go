package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/k-kaddal/bank-go/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	
	testDB, err = sql.Open(config.DBDriver, "postgresql://root:password@localhost:5432/bank_db?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}