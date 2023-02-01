package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/prachurjya15/simple-bank/api"
	db "github.com/prachurjya15/simple-bank/db/sqlc"
	"github.com/prachurjya15/simple-bank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cant Read Startup Config")
	}
	dbDriver := config.DbDriver
	dbSource := config.DbSource
	address := config.Address
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("error connecting to DB: %#v \n", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatalf("Error Starting Server %s", err)
	}
	err = server.StartServer(address)

	if err != nil {
		log.Fatalf("Error Starting Server %s", err)
	}
}
