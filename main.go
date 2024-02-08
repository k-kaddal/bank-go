package main

import (
	"database/sql"
	"log"

	"github.com/k-kaddal/bank-go/api"
	db "github.com/k-kaddal/bank-go/database/sqlc"
	"github.com/k-kaddal/bank-go/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connection, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	store := db.NewStore(connection)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	
}