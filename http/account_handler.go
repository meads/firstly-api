package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/db"
)

type accountRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *FirstlyServer) LoginAccountHandler(store db.Store) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var req accountRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		// query the account based on hmac etc.
		// return token

		// image, err := store.Login(ctx, req.Data)
		// if err != nil {
		// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		// 	return
		// }

		ctx.JSON(http.StatusOK, req)
	}
}

func (server *FirstlyServer) CreateAccountHandler(store db.Store) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var req accountRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		// image, err := store.Login(ctx, req.Data)
		// if err != nil {
		// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		// 	return
		// }

		ctx.JSON(http.StatusOK, req)
	}
}
