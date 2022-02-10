package main

import (
	"database/sql"
	"log"

	"github.com/danglebary/beatstore-backend-go/api"
	db "github.com/danglebary/beatstore-backend-go/db/sqlc"
	"github.com/danglebary/beatstore-backend-go/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load config settings", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Failed to connect to db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Failed to start the server", err)
	}
}
