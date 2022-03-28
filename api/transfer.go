package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/token"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency"  binding:"required,currency"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}
	fromAccount, valid := s.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}
	_, valid = s.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if payload == nil || payload.Username != fromAccount.Owner {
		ctx.JSON(http.StatusUnauthorized, authAccountNotBelongUserError)
		return
	}
	//insert db account
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	result, err := s.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (s *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, responseError(err))
			return account, false
		}
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return account, false
	}
	if account.Currency != currency {
		err = fmt.Errorf("account[%d] currency mismatch %s vs %s",
			accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return account, false
	}
	return account, true
}
