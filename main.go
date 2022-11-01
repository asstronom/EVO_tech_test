package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/asstronom/EVO_tech_test/pkg/db"
	"github.com/asstronom/EVO_tech_test/pkg/migratedb"
	"github.com/asstronom/EVO_tech_test/pkg/server"
)

func main() {
	var err error
	dbuser := os.Getenv("POSTGRES_USER")
	dbpass := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	dbhost := os.Getenv("DATABASE_HOST")
	dbport := os.Getenv("DATABASE_PORT")
	dburl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbuser, dbpass, dbhost, dbport, dbname)
	migrateurl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbuser, dbpass, dbhost, dbport, dbname)
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
