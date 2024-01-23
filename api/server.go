package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/kelvinator07/golang-bank-microservices/db/sqlc"
	"github.com/kelvinator07/golang-bank-microservices/token"
	"github.com/kelvinator07/golang-bank-microservices/util"
)

const (
	SuccessStatusCode    string = "00"
	FailedStatusCode     string = "99"
	SuccessStatusMessage string = "Success"
	FailedStatusMessage  string = "Failed"
)

type Server struct {
	config     util.Env
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Env, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	// server := &Server{store: store}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currencyCode", validCurrency)
	}

	// Add logging middleware
	// router.Use(RequestLogger())
	// router.Use(ResponseLogger())

	// server.router = router
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	router.POST("/api/v1/users", server.createUser)
	router.POST("/api/v1/users/login", server.loginUser)
	authRoutes.GET("/api/v1/users/:id", server.getOneUser)
	authRoutes.GET("/api/v1/users", server.getAllUsers2)
	// authRoutes.GET("/api/v1/users", server.getAllUsers2)

	authRoutes.POST("/api/v1/accounts", server.createAccount)
	authRoutes.GET("/api/v1/accounts/:id", server.getAccount)
	authRoutes.GET("/api/v1/accounts", server.getAllAccounts)

	authRoutes.POST("/api/v1/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(serverAddress string) error {
	return server.router.Run(serverAddress)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"status_code": FailedStatusCode,
		"message":     FailedStatusMessage,
		"error":       err.Error(),
	}
}

type HttpResponse struct {
	StatusCode string `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data"` // Use generics, also for tests
}

func NewHttpResponse(statusCode string, message string, data any) *HttpResponse {
	return &HttpResponse{statusCode, message, data}
}

func validResponse(d any) gin.H {
	return gin.H{
		"status_code": SuccessStatusCode,
		"data":        d,
		"message":     SuccessStatusMessage,
	}
}

type HttpResponseG[T any] struct {
	StatusCode string `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data"` // Use generics, also for tests
}

func NewHttpResponseG[T any](statusCode string, message string, data T) *HttpResponseG[T] {
	return &HttpResponseG[T]{statusCode, message, data}
}

type ErrorResponse struct {
	StatusCode string `json:"status_code"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}

func NewErrorResponse(statusCode string, message string, e error) *ErrorResponse {
	return &ErrorResponse{statusCode, message, e.Error()}
}
