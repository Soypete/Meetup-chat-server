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

	return PG{
		Client: db,
	}
}
