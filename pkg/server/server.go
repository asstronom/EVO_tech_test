package server

import (
	"fmt"

	"github.com/asstronom/EVO_tech_test/pkg/parse"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
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

func initEndpoints(router *gin.Engine) {
	router.GET("/transactions/:id", getTransactionByID)
	router.GET("/transactions", getTransactions)
}

func NewServer() (*Server, error) {
	router := gin.Default()
	initEndpoints(router)
	return &Server{router: router}, nil
}
