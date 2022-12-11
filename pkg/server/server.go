package server

import (
	"fmt"

	"github.com/asstronom/EVO_tech_test/pkg/domain"
	"github.com/asstronom/EVO_tech_test/pkg/parse"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router  *gin.Engine
	service domain.TransactionService
}

// run server
func (srv *Server) Run(port string) error {
	err := parse.ValidatePort(port)
	if err != nil {
		return fmt.Errorf("bad port number: %w", err)
	}
	err = srv.router.Run(port)
	if err != nil {
		return fmt.Errorf("error running router: %w", err)
	}
	return nil
}

// add endpoints to router
func (srv *Server) initEndpoints() {
	srv.router.GET("/transactions/:id", srv.transactionByID)
	srv.router.GET("/transactions", srv.transactions)
	srv.router.POST("/upload", srv.uploadCSV)
}

// create server
func NewServer(service domain.TransactionService) (*Server, error) {
	router := gin.Default()
	srv := Server{router: router, service: service}
	srv.initEndpoints()
	return &srv, nil
}
