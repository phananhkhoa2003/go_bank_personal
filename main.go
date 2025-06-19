package main

import (
	"context"
	"log"
	"simple_bank/api"
	db "simple_bank/db/sqlc"
	"simple_bank/util"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	connPool, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(connPool)

	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

}
