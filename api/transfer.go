package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/kelvinator07/golang-bank-microservices/db/sqlc"
	"github.com/kelvinator07/golang-bank-microservices/token"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	CurrencyCode  string `json:"currency_code" binding:"required,currencyCode"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.CurrencyCode)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.UserID != authPayload.UserID {
		err := errors.New("from account doesnt belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validAccountBalance(ctx, req.FromAccountID, req.Amount)
	if !valid {
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccountID, req.CurrencyCode)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, validResponse(result))
}

func (server Server) validAccount(ctx *gin.Context, accountID int64, currencyCode string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("account with id %v doesnt exist", accountID)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		// unexpected error, then return internal server
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.CurrencyCode != currencyCode {
		err := fmt.Errorf("account [%d] currency mistmatch: %s vs %s", account.ID, account.CurrencyCode, currencyCode)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}

func (server Server) validAccountBalance(ctx *gin.Context, accountID int64, amount int64) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		// unexpected error, then return internal server
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if amount > account.Balance {
		err := fmt.Errorf("account [%d] doesn't have enough balance: %v", account.ID, account.Balance)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
