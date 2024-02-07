package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:password@localhost:5432/bank_db?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	connection, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	testQueries = New(connection)

	os.Exit(m.Run())
}