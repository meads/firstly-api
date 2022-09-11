package http

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/db"
)

type createImageRequest struct {
	Data string `json:"data" binding:"required"`
}

func (server *FirstlyServer) CreateImageHandler(store db.Store) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var req createImageRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		image, err := store.Create(ctx, req.Data)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, image)
	}
}

func (server *FirstlyServer) DeleteImageHandler(store db.Store) func(*gin.Context) {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		if idParam == "" {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("id parameter is required"))
			return
		}
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("id parameter must be a valid integer"))
			return
		}

		err = store.Delete(ctx, id)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, nil)
	}
}

func (server *FirstlyServer) ListImagesHandler(store db.Store) func(ctx *gin.Context) {
	getLimitAndOffset := func(ctx *gin.Context) (string, string) {
		limit := ctx.Query("limit")
		if limit == "0" || limit == "" {
			limit = "50"
		}
		offset := ctx.Query("offset")
		if offset == "" {
			offset = "0"
		}

		return limit, offset
	}
	return func(ctx *gin.Context) {
		limit, offset := getLimitAndOffset(ctx)
		i, err := strconv.ParseInt(limit, 10, 32)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("error parsing limit as int"))
			return
		}

		j, err := strconv.ParseInt(offset, 10, 32)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("error parsing offset as int"))
			return
		}

		images, err := store.List(ctx, db.ListParams{Limit: int32(i), Offset: int32(j)})
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.JSON(http.StatusOK, images)
	}
}

type updateImageRequest struct {
	ID   int64  `json:"id" binding:"required"`
	Memo string `json:"memo" binding:"required"`
}

func (server *FirstlyServer) UpdateImageHandler(store db.Store) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req updateImageRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		image, err := store.Get(ctx, req.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var updateParams db.UpdateParams
		updateParams.ID = req.ID
		updateParams.Memo = req.Memo

		err = store.Update(ctx, updateParams)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		image.Memo = req.Memo
		ctx.JSON(http.StatusOK, image)
	}
}
