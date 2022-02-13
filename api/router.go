package api

import "github.com/gin-gonic/gin"

func newRouter(server *Server) *gin.Engine {
	router := gin.Default()

	// User routes
	router.POST("/users", server.createUser)
	router.POST("/users/:id", server.updateUser)
	router.GET("/users/:id", server.getUserById)
	router.GET("/users", server.listUsers)

	// Beat routes
	router.POST("/beats", server.createBeat)
	router.POST("/beats/:id", server.updateBeat)
	router.GET("/beats/:id", server.getBeat)
	router.GET("/beats", server.listBeatsById)
	router.GET("/users/:id/beats", server.listBeatsByCreatorId)

	// Like routes
	router.POST("/likes", server.createLike)
	router.GET("/likes/:uid/:bid", server.getLike)
	router.GET("/beats/:id/likes", server.listLikesByBeatID)
	router.GET("/users/:id/likes", server.listLikesByUserID)

	return router
}
