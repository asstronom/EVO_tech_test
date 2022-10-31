package server

import "testing"

func TestServer(t *testing.T) {
	srv, err := NewServer()
	if err != nil {
		t.Errorf("error creating server: %s", err)
	}
	err = srv.Run(":")
	if err != nil {
		t.Errorf("error starting server: %s", err)
	}
}
