package server

import (
	"fmt"

	"github.com/asstronom/EVO_tech_test/pkg/db"
	"github.com/asstronom/EVO_tech_test/pkg/parse"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	db     *db.TransactionDB
}

func (s *Server) Run(port string) error {
	err := parse.ValidatePort(port)
	if err != nil {
		return fmt.Errorf("bad port number: %w", err)
	}
	err = s.router.Run(port)
	if err != nil {
		return fmt.Errorf("error running router: %w", err)
	}
	return nil
}

func (srv *Server) initEndpoints() {
	srv.router.GET("/transactions/:id", srv.transactionByID)
	srv.router.GET("/transactions", srv.transactions)
	srv.router.POST("/upload", srv.uploadCSV)
}

func NewServer(db *db.TransactionDB) (*Server, error) {
	router := gin.Default()
	srv := Server{router: router, db: db}
	srv.initEndpoints()
	return &srv, nil
}
