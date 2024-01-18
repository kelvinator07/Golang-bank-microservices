package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/kelvinator07/golang-bank-microservices/db/sqlc"
	"github.com/kelvinator07/golang-bank-microservices/util"
	"github.com/lib/pq"
)

type createUserRequest struct {
	AccountName string `json:"account_name" binding:"required,min=5"`
	Password    string `json:"password" binding:"required,min=6"`
	Address     string `json:"address" binding:"required,min=5"`
	Gender      string `json:"gender" binding:"required,oneof=MALE FEMALE"`
	PhoneNumber int64  `json:"phone_number" binding:"required"` // TODO Validate PhoneNumber Properly e164=string with +234
	Email       string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	UserID      int64  `json:"user_id"`
	AccountName string `json:"account_name"`
	Address     string `json:"address"`
	Gender      string `json:"gender"`
	PhoneNumber int64  `json:"phone_number"`
	Email       string `json:"email"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get user ID from request
	arg := db.CreateUserParams{
		AccountName:    req.AccountName,
		HashedPassword: hashedPassword,
		Address:        req.Address,
		Gender:         req.Gender,
		PhoneNumber:    req.PhoneNumber,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := createUserResponse{
		UserID:      user.ID,
		AccountName: user.AccountName,
		Address:     user.Address,
		Gender:      user.Gender,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
	}

	ctx.JSON(http.StatusOK, validResponse(res))
}

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getOneUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := createUserResponse{
		UserID:      user.ID,
		AccountName: user.AccountName,
		Address:     user.Address,
		Gender:      user.Gender,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
	}

	ctx.JSON(http.StatusOK, validResponse(res))
}

type getAllUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) getAllUsers(ctx *gin.Context) {
	var req getAllUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var usersList []createUserResponse
	for _, user := range users {
		usersList = append(usersList, createUserResponse{
			UserID:      user.ID,
			AccountName: user.AccountName,
			Address:     user.Address,
			Gender:      user.Gender,
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
		})
	}

	ctx.JSON(http.StatusOK, validResponse(usersList))
}
