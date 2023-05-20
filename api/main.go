package main

import (
	"database/sql"
	"log"

	"github.com/ot07/next-bazaar/api"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/util"

	_ "github.com/lib/pq"
	_ "github.com/ot07/next-bazaar/docs"
)

// @title Next Bazaar API
// @version 0.0.1
// @BasePath /api/v1
func main() {
	config, err := util.LoadConfig()
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
		log.Fatal("cannot start to server:", err)
	}
}
