package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/db/sqlc"
)

type CreateNoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	UserID  int64  `json:"userId" binding:"required"`
}

type CreateNoteResponse struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int64  `json:"userId"`
}

func createNoteHandler(ctx *gin.Context) {
	var req CreateNoteRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userExists, err := firstly.store.UserExists(ctx, req.UserID)
	// if there is an error and it isn't because of no db rows then error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// if the userExists was false and there isn't any error, fail because the user is not found
	if !userExists {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("user not found")))
		return
	}

	dbNote, err := firstly.store.CreateNote(ctx, db.CreateNoteParams{
		Title:   req.Title,
		Content: req.Content,
		UserID:  req.UserID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, CreateNoteResponse{
		ID:      dbNote.ID,
		Title:   dbNote.Title,
		Content: dbNote.Content,
		UserID:  dbNote.UserID,
		// CreatedAt: dbNote.CreatedAt,
	})
}

type ListNotesResponse struct {
	Notes []db.Note `json:"notes"`
}

func listNotesHandler(ctx *gin.Context) {
	userIDParam := ctx.Param("userid")
	if userIDParam == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("missing user id param")))
		return
	}

	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid user id param")))
		return
	}

	userExists, err := firstly.store.UserExists(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !userExists {
		ctx.JSON(http.StatusNotFound, errorResponse(errors.New("user not found")))
		return
	}

	notes, err := firstly.store.ListNotesByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ListNotesResponse{
		Notes: notes,
	})
}

type DeleteNoteResponse struct {
	Message string `json:"message"`
}

func deleteNoteHandler(ctx *gin.Context) {
	userIDParam := ctx.Param("userid")
	if userIDParam == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("missing user id param")))
		return
	}

	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid user id param")))
		return
	}

	noteIDParam := ctx.Param("noteid")
	if noteIDParam == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("missing note id param")))
		return
	}

	noteID, err := strconv.ParseInt(noteIDParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid note id param")))
		return
	}

	userExists, err := firstly.store.UserExists(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !userExists {
		ctx.JSON(http.StatusNotFound, errorResponse(errors.New("user not found")))
		return
	}

	err = firstly.store.DeleteNote(ctx, db.DeleteNoteParams{
		ID:     noteID,
		UserID: userID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, DeleteNoteResponse{
		Message: "Resource successfully deleted",
	})
}

type UpdateNoteRequest struct {
	ID      int64  `json:"id" binding:"required"`
	UserID  int64  `json:"userId" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateNoteResponse struct {
	ID      int64  `json:"id"`
	UserID  int64  `json:"userId"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func updateNoteHandler(ctx *gin.Context) {
	var req UpdateNoteRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userExists, err := firstly.store.UserExists(ctx, req.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !userExists {
		ctx.JSON(http.StatusNotFound, errorResponse(errors.New("user not found")))
		return
	}

	err = firstly.store.UpdateNote(ctx, db.UpdateNoteParams{
		Title:   req.Title,
		Content: req.Content,
		ID:      req.ID,
		UserID:  req.UserID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, UpdateNoteResponse{
		ID:      req.ID,
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	})
}
