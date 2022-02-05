package main

import (
	"database/sql"
	"log"

	"github.com/danglebary/beatstore-backend-go/api"
	db "github.com/danglebary/beatstore-backend-go/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/beatstore?sslmode=disable"
	serverAddress = "0.0.0.0:1337"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Failed to connect to db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Failed to start the server", err)
	}
}
