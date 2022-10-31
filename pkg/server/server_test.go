package server

import "testing"

func TestServer(t *testing.T) {
	srv, err := NewServer()
	if err != nil {
		t.Errorf("error creating server: %s", err)
	}
	err = srv.Run(":8080")
	if err != nil {
		t.Errorf("error starting server: %s", err)
	}
}
