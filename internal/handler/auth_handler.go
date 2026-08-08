package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/meads/firstly-api/internal/service"
)

type AuthHandler struct {
	svc service.AuthServicer
}

func NewAuthHandler(svc service.AuthServicer) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	registerResult, err := h.svc.Register(ctx, req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("Registration failed: %w", err)))
		return
	}

	ctx.JSON(http.StatusOK, registerResult)
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	loginResult, err := h.svc.Login(ctx, req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("login failed: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, &LoginResponse{
		SessionID:             loginResult.SessionID,
		AccessToken:           loginResult.AccessToken,
		RefreshToken:          loginResult.RefreshToken,
		AccessTokenExpiresAt:  loginResult.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: loginResult.RefreshTokenExpiresAt,
		Username:              loginResult.Username,
		UserID:                loginResult.UserID,
	})
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	idParam := ctx.Param("sessionid")
	err := h.svc.Logout(ctx, idParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("logout failed %w", err)))
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) RenewAccessToken(ctx *gin.Context) {
	var req RenewAccessTokenRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	renewAccessTokenResult, err := h.svc.RenewAccessToken(ctx, req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("error renewing access token: %w", err)))
		return
	}

	ctx.JSON(http.StatusOK, renewAccessTokenResult)
}

func (h *AuthHandler) RevokeSession(ctx *gin.Context) {
	idParam := ctx.Param("sessionid")
	err := h.svc.RevokeSession(ctx, idParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("error revoking session: %w", err))
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
}
