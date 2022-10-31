package server

import (
	"context"
	"log"
	"testing"

	"github.com/asstronom/EVO_tech_test/pkg/db"
)

func TestServer(t *testing.T) {
	dburl := "postgres://user:mypassword@localhost:5432/transactions"
	trdb, err := db.Open(context.Background(), dburl)
	if err != nil {
		log.Fatalln(err)
	}
	srv, err := NewServer(trdb)
	if err != nil {
		t.Errorf("error creating server: %s", err)
	}
	err = srv.Run(":8080")
	if err != nil {
		t.Errorf("error starting server: %s", err)
	}
}
