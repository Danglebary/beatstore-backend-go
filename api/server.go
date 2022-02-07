package api

import (
	db "github.com/danglebary/beatstore-backend-go/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Serves all HTTP requests for our service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// Creates a new HTTP server instance and initializes routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.GET("/", playgroundHandler())
	router.POST("/", graphqlHandler())

	server.router = router
	return server
}

// Runs HTTP server on specified address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
