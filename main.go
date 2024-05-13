package main

import (
	"database/sql"
	"log"

	"github.com/debugroach/video-hub-serve/api"
	"github.com/debugroach/video-hub-serve/config"
	db "github.com/debugroach/video-hub-serve/db/sqlc"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	conn, err := sql.Open("mysql", config.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
