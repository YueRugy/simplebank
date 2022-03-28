package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/token"
)

type createAccountRequest struct {
	Currency string `json:"currency"  binding:"required,currency"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	//insert db account
	arg := db.CreateAccountParams{
		Owner:    payload.Username,
		Balance:  0,
		Currency: req.Currency,
	}
	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}
	account, err := s.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, responseError(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if account.Owner != payload.Username {
		ctx.JSON(http.StatusUnauthorized,
			responseError(authAccountNotBelongUserError))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageSize   int32 `form:"page_size" binding:"required,min=1"`
	PageNumber int32 `form:"page_number" binding:"required,min=5,max=20"`
}

func (s *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	args := db.ListAccountsParams{
		Owner:  payload.Username,
		Limit:  req.PageNumber,
		Offset: (req.PageSize - 1) * req.PageNumber,
	}
	accounts, err := s.store.ListAccounts(ctx, args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, responseError(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
