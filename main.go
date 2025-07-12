package main

import (
	"database/sql"
	"fmt"
	"log"
	"simle_bank/api"
	db "simle_bank/db/sqlc"
	"simle_bank/util"

	_ "github.com/lib/pq"
)

func init() {
	fmt.Println("🔍 Registered SQL drivers:", sql.Drivers())
}

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
