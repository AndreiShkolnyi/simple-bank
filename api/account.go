package api

import (
	"database/sql"
	"fmt"
	"net/http"
	db "simle_bank/db/sqlc"
	"simle_bank/token"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,oneof=USD EUR RUB"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}
	account, err := server.db.CreateAccount(ctx, arg)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
				case "foreign_key_violation", "unique_violation":
					ctx.JSON(http.StatusForbidden, errorResponse(err))
					return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}


func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.db.GetAccount(ctx, req.ID); 
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
		}
	
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := fmt.Errorf("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type ListAccountRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req ListAccountRequest

	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	accounts, err := server.db.ListAccounts(ctx, arg )
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}