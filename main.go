package main

import (
	"database/sql"
	"log"

	"github.com/debugroach/movie-hub-serve/api"
	"github.com/debugroach/movie-hub-serve/config"
	"github.com/debugroach/movie-hub-serve/db"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	conn, err := sql.Open("mysql", config.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	q := db.New(conn)
	server := api.NewServer(q)

	err = server.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
