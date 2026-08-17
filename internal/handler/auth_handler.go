package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meads/firstly-api/internal/domain"
)

type AuthServicer interface {
	Register(ctx context.Context, username, password string) (*domain.RegisterResult, error)
	Login(ctx context.Context, username, password string) (*domain.LoginResult, error)
	Logout(ctx context.Context, sessionID string) error
	RenewAccessToken(ctx context.Context, refreshToken string) (*domain.RenewAccessTokenResult, error)
	RevokeSession(ctx context.Context, sessionID string) error
}

type AuthHandler struct {
	authService AuthServicer
}

func NewAuthHandler(authService AuthServicer) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	registerResult, err := h.authService.Register(ctx, req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("Registration failed")))
		return
	}

	ctx.JSON(http.StatusOK, MapToRegisterResponse(registerResult))
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	loginResult, err := h.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			ctx.JSON(http.StatusUnauthorized, errorResponse(domain.ErrInvalidCredentials))

		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("internal server error")))
		}
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
	err := h.authService.Logout(ctx, idParam)
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

	renewAccessTokenResult, err := h.authService.RenewAccessToken(ctx, req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("error renewing access token: %w", err)))
		return
	}

	ctx.JSON(http.StatusOK, MapToRenewAccessTokenResponse(renewAccessTokenResult))
}

func (h *AuthHandler) RevokeSession(ctx *gin.Context) {
	idParam := ctx.Param("sessionid")
	err := h.authService.RevokeSession(ctx, idParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("error revoking session: %w", err))
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
}
