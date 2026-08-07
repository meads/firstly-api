package http

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/internal/db/sqlc"
)

type registerUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerUserResponse struct {
	SessionID             string    `json:"sessionId"`
	AccessToken           string    `json:"accessToken"`
	RefreshToken          string    `json:"refreshToken"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
	Username              string    `json:"username"`
	UserID                int64     `json:"userId"`
}

func registerUserHandler(ctx *gin.Context) {
	var req registerUserRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	usernameExists, err := firstly.store.UsernameExists(ctx, req.Username)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if usernameExists {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("please choose another username")))
		return
	}

	var param db.CreateUserParams
	param.Username = req.Username
	param.Password, err = firstly.hasher.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	dbUser, err := firstly.store.CreateUser(ctx, param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, accessClaims, err := firstly.tokener.GenerateToken(dbUser.ID, dbUser.Username, "access", 15*time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshClaims, err := firstly.tokener.GenerateToken(dbUser.ID, dbUser.Username, "refresh", 24*time.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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

	ctx.JSON(http.StatusOK, registerUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		Username:              dbUser.Username,
		UserID:                dbUser.ID,
	})

	// ctx.SetCookieData(&http.Cookie{
	// 	Name:     "token",
	// 	Value:    tokenString,
	// 	Expires:  time.Now().Add(time.Minute * 5),
	// 	Path:     "/",
	// 	Domain:   "localhost",
	// 	SameSite: http.SameSiteNoneMode,
	// 	Secure:   true,
	// 	HttpOnly: true,
	// })
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}

func deleteUserHandler(ctx *gin.Context) {
	idParam := ctx.Param("userid")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "id parameter is required",
		})

		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "id parameter must be a valid integer",
		})

		return
	}

	err = firstly.store.DeleteUser(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, DeleteUserResponse{
		Message: "Resource successfully deleted",
	})
}

type ListUsersResponse struct {
	Users []db.User `json:"users"`
}

func listUsersHandler(ctx *gin.Context) {
	limit := ctx.Query("limit")
	if limit == "0" || limit == "" {
		limit = "50"
	}

	offset := ctx.Query("offset")
	if offset == "" {
		offset = "0"
	}

	i, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error parsing limit as int",
		})

		return
	}

	j, err := strconv.ParseInt(offset, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error parsing offset as int",
		})

		return
	}

	users, err := firstly.store.ListUsers(ctx, db.ListUsersParams{Limit: int32(i), Offset: int32(j)})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})

		return
	}

	ctx.JSON(http.StatusOK, ListUsersResponse{Users: users})
}

type PatchUserRequest struct {
	ID              int64  `json:"id" binding:"required"`
	Username        string `json:"username" binding:"required"`
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}

type PatchUserResponse struct {
	Message string `json:"message"`
}

func patchUserHandler(ctx *gin.Context) {
	var req PatchUserRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := firstly.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = firstly.hasher.ComparePassword(user.Password, req.CurrentPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("current password invalid")))
		return
	}

	newPasswordHash, err := firstly.hasher.HashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var updateParams db.UpdateUserPasswordParams
	updateParams.ID = user.ID
	updateParams.Password = newPasswordHash

	err = firstly.store.UpdateUserPassword(ctx, updateParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, PatchUserResponse{
		Message: "Resource successfully patched",
	})
}
