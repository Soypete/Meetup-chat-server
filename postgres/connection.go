package postgres

import (
	"log"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type PG struct {
	Client *sqlx.DB
}

func ConnectDB() PG {
	db, err := sqlx.Connect("postgres", "user=postgres password=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	return PG{
		Client: db,
	}
}
