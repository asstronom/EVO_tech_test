package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func (s *Server) Run(port string) error {
	err := s.router.Run(port)
	if err != nil {
		return fmt.Errorf("error running router: %w", err)
	}
	return nil
}

func initEndpoints(router *gin.Engine) {
	router.GET("/transaction/:id", getTransactionByID)
}

func NewServer() (*Server, error) {
	router := gin.Default()
	initEndpoints(router)
	return &Server{router: router}, nil
}