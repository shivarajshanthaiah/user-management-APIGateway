package server

import "github.com/gin-gonic/gin"

type Server struct {
	R *gin.Engine
}

func (s *Server) StartServer(port string) {
	s.R.Run(":" + port)
}

func NewServer() *Server {
	engine := gin.Default()
	return &Server{
		R: engine,
	}
}
