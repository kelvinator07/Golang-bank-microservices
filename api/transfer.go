package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/kelvinator07/golang-bank-microservices/db/sqlc"
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

	if !server.validAccount(ctx, req.FromAccountID, req.CurrencyCode) {
		return
	}

	if !server.validAccountBalance(ctx, req.FromAccountID, req.Amount) {
		return
	}

	if !server.validAccount(ctx, req.ToAccountID, req.CurrencyCode) {
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

	ctx.JSON(http.StatusOK, result)
	// TODO: Generics
	// ctx.JSON(http.StatusOK, *NewHttpResponse{"status", "message", "data"})
}

func (server Server) validAccount(ctx *gin.Context, accountID int64, currencyCode string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		// unexpected error, then return internal server
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.CurrencyCode != currencyCode {
		err := fmt.Errorf("account [%d] currency mistmatch: %s vs %s", account.ID, account.CurrencyCode, currencyCode)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}

func (server Server) validAccountBalance(ctx *gin.Context, accountID int64, amount int64) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		// unexpected error, then return internal server
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if amount > account.Balance {
		err := fmt.Errorf("account [%d] doesn't have enough balance: %v", account.ID, account.Balance)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
