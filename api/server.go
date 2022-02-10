package api

import (
	db "github.com/danglebary/beatstore-backend-go/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Serves all HTTP requests for our service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// Creates a new HTTP server instance and initializes routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// User routes
	router.POST("/users", server.createUser)
	router.POST("/users/:id", server.updateUser)
	router.GET("/users/:id", server.getUserById)
	router.GET("/users/:username", server.getUserByUsername)
	router.GET("/users", server.listUsers)

	// Beat routes
	router.POST("/users/:id/create-beat", server.createBeat)
	router.POST("/beats/:id", server.updateBeat)
	router.GET("/beats/:id", server.getBeat)
	router.GET("/beats", server.listBeatsById)
	router.GET("/users/:id/beats", server.listBeatsByCreatorId)

	// Like routes
	router.POST("/likes", server.createLike)
	router.GET("/likes/:uid/:bid", server.getLike)
	router.GET("/likes/:uid", server.listLikesByUserId)
	router.GET("/likes/:bid", server.listLikesByBeatId)

	server.router = router
	return server
}

// Runs HTTP server on specified address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
