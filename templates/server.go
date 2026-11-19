package main

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	r *gin.Engine
}

func NewServer() *Server {
	r := gin.Default()

	return &Server{
		r,
	}
}

func (s *Server) start() error {
	return s.r.Run()
}
