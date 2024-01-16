package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/kelvinator07/golang-bank-microservices/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Add logging middleware
	// router.Use(RequestLogger())
	// router.Use(ResponseLogger())

	router.POST("/api/v1/accounts", server.createAccount)
	router.GET("/api/v1/accounts/:id", server.getAccount)
	router.GET("/api/v1/accounts", server.getAllAccounts)

	server.router = router
	return server
}

func (server *Server) Start(serverAddress string) error {
	return server.router.Run(serverAddress)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
