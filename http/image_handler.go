package http

import (
	"errors"
	"log"
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

// type getImageRequest struct {
// 	ID int64 `json:"id" binding:"required"`
// }

// func (api *FirstlyServer) getImage(ctx *gin.Context) {
// var req getConfigRequest

// idString, ok := ctx.Params.Get("id")
// id, err := strconv.Atoi(idString)
// if err != nil || !ok {
// 	ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("id is required")))
// 	return
// }

// req.ID = int64(id)
// config, err := server.store.GetConfig(ctx, req.ID)
// if err != nil {
// 	if err == sql.ErrNoRows {
// 		ctx.JSON(http.StatusNotFound, nil)
// 		return
// 	}
// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 	return
// }
// ctx.JSON(http.StatusOK, config)
// 	return
// }

func (server *FirstlyServer) ListImagesHandler(store db.Store) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		images, err := store.List(ctx)
		if err != nil {
			log.Fatalf("Error calling ListImages: %s", err)
			return
		}

		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.JSON(http.StatusOK, images)
	}

	// limit := ctx.Query("limit")
	// if limit == "0" || limit == "" {
	// 	limit = "50"
	// }
	// offset := ctx.Query("offset")
	// if offset == "" {
	// 	offset = "0"
	// }
	// params := db.ListImages{}

	// limitInt, errLimitConvert := strconv.Atoi(limit)
	// offsetInt, errOffsetConvert := strconv.Atoi(offset)
	// if errLimitConvert == nil && errOffsetConvert == nil {
	// 	params.Limit = int32(limitInt)
	// 	params.Offset = int32(offsetInt)
	// }

	// configs, err := server.store.ListConfigs(ctx, params)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	// 	return
	// }

	// ctx.JSON(http.StatusOK, configs)
}

// type updateImageRequest struct {
// 	ID   int64  `json:"id" binding:"required"`
// 	Name string `json:"name" binding:"required"`
// }

// func (server *Server) updateConfig(ctx *gin.Context) {
// 	var req updateConfigRequest

// 	// validate the request
// 	if err := ctx.BindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	// query the record for update
// 	config, err := server.store.GetConfigForUpdate(ctx, req.ID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	// overlay the request values specified
// 	var configParams db.UpdateConfigParams
// 	configParams.ID = req.ID
// 	configParams.Name = req.Name

// 	// issue the update with the update params struct
// 	config, err = server.store.UpdateConfig(ctx, configParams)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, config)
// }
