package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/kelvinator07/golang-bank-microservices/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currencyCode", validCurrency)
	}

	// Add logging middleware
	// router.Use(RequestLogger())
	// router.Use(ResponseLogger())

	router.POST("/api/v1/accounts", server.createAccount)
	router.GET("/api/v1/accounts/:id", server.getAccount)
	router.GET("/api/v1/accounts", server.getAllAccounts)

	router.POST("/api/v1/transfers", server.createTransfer)

	server.router = router
	return server
}

func (server *Server) Start(serverAddress string) error {
	return server.router.Run(serverAddress)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

type HttpResponse[T any] struct {
	status  string `json:"status"`
	message string `json:"message"`
	data    T      `json:"data"`
}

func NewHttpResponse[T any](status string, message string, d T) *HttpResponse[any] {
	return &HttpResponse[any]{status, message, d}
}

func validResponse[T any](t T) gin.H {
	fmt.Println("T ", t)
	val := &HttpResponse[T]{"00", "Success", t}
	fmt.Println("T val ", val)
	return gin.H{"result": val}
}
