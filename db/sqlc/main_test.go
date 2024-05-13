package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var queries *Queries

func TestMain(m *testing.M) {
	db, err := sql.Open("mysql", "root:secret@tcp(127.0.0.1:3306)/movie_hub?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	queries = New(db)
	os.Exit(m.Run())
}
