package main

import (
	"context"
	"log"

	"github.com/asstronom/EVO_tech_test/pkg/db"
	"github.com/asstronom/EVO_tech_test/pkg/migratedb"
	"github.com/asstronom/EVO_tech_test/pkg/server"
)

func main() {
	var err error
	dburl := "postgres://user:mypassword@localhost:5432/transactions"
	migrateurl := "postgres://user:mypassword@localhost:5432/transactions?sslmode=disable"
	err = migratedb.MigrateUp(migrateurl)
	if err != nil {
		log.Fatalln(err)
	}
	trdb, err := db.Open(context.Background(), dburl)
	if err != nil {
		log.Fatalln(err)
	}
	srv, err := server.NewServer(trdb)
	if err != nil {
		log.Fatalln("error starting server", err)
	}
	srv.Run(":8080")
}
