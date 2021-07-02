package postgres

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

func ConnectDB() {
	db, err := sqlx.Connect("postgres", "user=postgres password=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	fmt.Println(db.Ping())
}
