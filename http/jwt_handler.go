package http

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/db"
)

func claimsMiddleware(h gin.HandlerFunc) gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		authorizationHeader := ctx.Request.Header.Get("Authorization")
		if authorizationHeader == "" || !strings.Contains(authorizationHeader, "Bearer") {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("token header missing")))
			return
		}
		parts := strings.Split(authorizationHeader, " ")
		tokenString := ""
		if len(parts) == 2 {
			tokenString = parts[1]
		}
		userClaims, err := firstly.tokener.VerifyToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// _ = userClaims
		if userClaims.Type != "access" {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid token type")))
			return
		}

		h(ctx)
	})
}

type loginRequest struct {
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type loginResponse struct {
	SessionID             string    `json:"sessionId"`
	AccessToken           string    `json:"accessToken"`
	RefreshToken          string    `json:"refreshToken"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
	Username              string    `json:"username"`
	UserID                int64     `json:"userId"`
}

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"accessToken"`
	AccessTokenExpiresAt time.Time `json:"accessTokenExpiresAt"`
}

func loginHandler(ctx *gin.Context) {
	var req loginRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	dbUser, err := firstly.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = firstly.hasher.ComparePassword(dbUser.Password, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid username or password")))
		return
	}

	accessToken, accessClaims, err := firstly.tokener.GenerateToken(dbUser.ID, dbUser.Username, "access", 15*time.Minute)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	refreshToken, refreshClaims, err := firstly.tokener.GenerateToken(dbUser.ID, dbUser.Username, "refresh", 24*time.Hour)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	session, err := firstly.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshClaims.RegisteredClaims.ID,
		UserID:       dbUser.ID,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("error creating session")))
		return
	}
	// ctx.Writer.Header().Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	ctx.JSON(http.StatusOK, loginResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		Username:              dbUser.Username,
		UserID:                dbUser.ID,
	})
}

func logoutHandler(ctx *gin.Context) {
	idParam := ctx.Param("sessionid")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("missing session ID")))
		return
	}

	err := firstly.store.DeleteSession(ctx, idParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("error deleting session %w", err)))
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
}

func renewAccessTokenHandler(ctx *gin.Context) {
	var req renewAccessTokenRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshClaims, err := firstly.tokener.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("error verifying refresh token")))
		return
	}

	session, err := firstly.store.GetSession(ctx, refreshClaims.RegisteredClaims.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid session")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("error getting session")))
		return
	}

	if session.IsRevoked {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("session revoked")))
		return
	}

	// if session.Username != refreshClaims.Username {
	// 	ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid session")))
	// 	return
	// }

	accessToken, accessClaims, err := firstly.tokener.GenerateToken(
		refreshClaims.ID, refreshClaims.Username, "access", 15*time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("error creating token")))
		return
	}

	ctx.JSON(http.StatusOK, renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	})
}

func revokeSessionHandler(ctx *gin.Context) {
	idParam := ctx.Param("sessionid")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("missing session ID")))
		return
	}

	err := firstly.store.RevokeSession(ctx, idParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("error revoking session")))
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
}
