package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/meads/firstly-api/internal/service"
)

type NoteHandler struct {
	svc service.NoteServicer
}

func NewNoteHandler(svc service.NoteServicer) *NoteHandler {
	return &NoteHandler{svc: svc}
}

func (h *NoteHandler) CreateNote(ctx *gin.Context) {
	var req CreateNoteRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	note, err := h.svc.CreateNote(ctx, req.UserID, req.Title, req.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, note)
}

func (h *NoteHandler) UpdateNote(ctx *gin.Context) {
	var req UpdateNoteRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	note, err := h.svc.UpdateNote(ctx, req.ID, req.UserID, req.Title, req.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, note)
}

func (h *NoteHandler) ListNotes(ctx *gin.Context) {
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

	notes, err := h.svc.ListNotes(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, ListNotesResponse{
		Notes: notes,
	})
}

func (h *NoteHandler) DeleteNote(ctx *gin.Context) {
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

	err = h.svc.DeleteNote(ctx, noteID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("delete note failed: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, DeleteNoteResponse{
		Message: "Resource successfully deleted",
	})
}
