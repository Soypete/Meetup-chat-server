package postgres

import (
	"fmt"
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
	// sanity check
	fmt.Println(db.Ping())
	return PG{
		Client: db,
	}
}
