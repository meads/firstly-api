package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/db"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *FirstlyServer) LoginHandler(store db.Store) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var req loginRequest
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
