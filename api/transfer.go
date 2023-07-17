package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/alrasyidin/simplebank-go/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createTransferParams struct {
	FromAccountID int64  `json:"from_account_id" binding:"required"`
	ToAccountID   int64  `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var param createTransferParams

	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fmt.Println("param", param)

	if !server.validCurrency(ctx, param.FromAccountID, param.Currency) {
		return
	}
	if !server.validCurrency(ctx, param.ToAccountID, param.Currency) {
		return
	}

	data := db.TransferTxParam{
		Amount:        param.Amount,
		FromAccountID: param.FromAccountID,
		ToAccountID:   param.ToAccountID,
	}

	result, err := server.store.TransferTx(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validCurrency(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	fmt.Println("account", account)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
