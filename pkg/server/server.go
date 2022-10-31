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

func (srv *Server) initEndpoints(router *gin.Engine) {
	router.GET("/transactions/:id", srv.getTransactionByID)
	router.GET("/transactions", srv.getTransactions)
}

func NewServer(db *db.TransactionDB) (*Server, error) {
	router := gin.Default()
	srv := Server{router: router, db: db}
	srv.initEndpoints(router)
	return &srv, nil
}
