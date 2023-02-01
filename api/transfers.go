package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/prachurjya15/simple-bank/db/sqlc"
	"github.com/prachurjya15/simple-bank/token"
)

type CreateTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required"`
	ToAccountID   int64  `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) CreateTransfer(ctx *gin.Context) {
	var req CreateTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fromAccount, valid := s.validateCurrencyMatch(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}
	_, valid = s.validateCurrencyMatch(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}
	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)
	if authPayload.Username != fromAccount.Owner {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("can send from logged in user account only")))
		return
	}
	arg := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	account, err := s.store.Query.CreateTransfer(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, account)
}

func (s *Server) validateCurrencyMatch(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {
	account, err := s.store.Query.GetAccountById(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}
	if account.Currency != currency {
		err = fmt.Errorf("input currency of the request: [%s] and the account currency: [%s] with account id: [%d] doesnt match", currency, account.Currency, accountId)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}
	return account, true
}
