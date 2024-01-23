package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/kelvinator07/golang-bank-microservices/db/sqlc"
	"github.com/kelvinator07/golang-bank-microservices/token"
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

	ctx.JSON(http.StatusOK, NewHttpResponseG("00", "Success", res))
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// if errors.Is(err, db.ErrRecordNotFound) {
		// 	err = fmt.Errorf("user with email %s doesn't exist", req.Email)
		// 	ctx.JSON(http.StatusNotFound, errorResponse(err))
		// 	return
		// }
		err = fmt.Errorf("user with email %s doesn't exist", req.Email)
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	err = util.ComparePasswords(req.Password, user.HashedPassword)
	if err != nil {
		err = errors.New("email and password doesn't match")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	token, err := server.tokenMaker.CreateToken(user.ID, user.AccountName, user.Email, time.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := loginUserResponse{
		AccessToken: token,
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
			err = fmt.Errorf("user with id %v doesn't exist", req.ID)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if user.ID != authPayload.UserID {
		err := errors.New("account doesnt belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
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
	PageID   int32 `form:"page_id"`
	PageSize int32 `form:"page_size"`
}

func (server *Server) getAllUsers(ctx *gin.Context) {
	var req getAllUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// Set default
	if req.PageID <= 0 {
		req.PageID = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 10
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

type getAllUsersRequest2 struct {
	Limit  int32 `form:"limit"`
	Cursor int64 `form:"cursor"`
}

type getAllUsersResponse2 struct {
	Users  []createUserResponse `form:"users"`
	Cursor int64                `form:"cursor"`
}

func (server *Server) getAllUsers2(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	cursor, _ := strconv.Atoi(ctx.DefaultQuery("cursor", "0"))

	// GET /payments?limit=10
	// GET /payments?limit=10&cursor=last_id_from_previous_fetch

	// SELECT * FROM users WHERE id > $1 ORDER BY id LIMIT $2;

	arg := db.GetAllUsersParams{
		ID:    int64(cursor),
		Limit: int32(limit),
	}

	users, err := server.store.GetAllUsers(ctx, arg)
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

	var nextCursor int64
	if len(usersList) > 0 {
		nextCursor = usersList[len(usersList)-1].UserID // lastID
	}

	ctx.JSON(http.StatusOK, validResponse(getAllUsersResponse2{
		Users:  usersList,
		Cursor: nextCursor,
	}))
}
