package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/kelvinator07/golang-bank-microservices/db/sqlc"
)

const (
	SuccessStatusCode    string = "00"
	FailedStatusCode     string = "99"
	SuccessStatusMessage string = "Success"
	FailedStatusMessage  string = "Failed"
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

	router.POST("/api/v1/users", server.createUser)
	router.GET("/api/v1/users/:id", server.getOneUser)
	router.GET("/api/v1/users", server.getAllUsers)

	server.router = router
	return server
}

func (server *Server) Start(serverAddress string) error {
	return server.router.Run(serverAddress)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"status":  FailedStatusCode,
		"message": FailedStatusMessage,
		"error":   err.Error(),
	}
}

type HttpResponse struct {
	status  string `json:"status"`
	message string `json:"message"`
	data    any    `json:"data"`
}

func NewHttpResponse(status string, message string, data any) *HttpResponse {
	return &HttpResponse{status, message, data}
}

func validResponse(d any) gin.H {
	return gin.H{
		"status":  SuccessStatusCode,
		"data":    d,
		"message": SuccessStatusMessage,
	}
}
