package http

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/db"
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
}

func registerUserHandler(ctx *gin.Context) {
	var req registerUserRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check if there is an existing user with that username
	userExists, err := firstly.store.UserExists(ctx, req.Username)
	// if there is an error and it isn't because of no db rows then error
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// if a user exists already with that name, fail
	if userExists {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("please choose another username")))
		return
	}

	// hash the request password
	var param db.CreateUserParams
	param.Username = req.Username
	param.Password, err = firstly.hasher.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// create the user with the hashed password and username
	dbUser, err := firstly.store.CreateUser(ctx, param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, accessClaims, err := firstly.tokener.GenerateToken(dbUser.ID, dbUser.Username, "access", 5*time.Minute)
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
		Username:     dbUser.Username,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("error creating session")))
		return
	}
	// ctx.Writer.Header().Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	ctx.JSON(http.StatusOK, registerUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		Username:              dbUser.Username,
	})
	// // generate a 5 minute token for the new user
	// tokenString, userClaims, err := firstly.tokener.GenerateToken(dbUser.ID, dbUser.Username, "access", 5*time.Minute)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	// 	return
	// }

	// // do something with userClaims
	// _ = userClaims

	// // ctx.SetCookieData(&http.Cookie{
	// // 	Name:     "token",
	// // 	Value:    tokenString,
	// // 	Expires:  time.Now().Add(time.Minute * 5),
	// // 	Path:     "/",
	// // 	Domain:   "localhost",
	// // 	SameSite: http.SameSiteNoneMode,
	// // 	Secure:   true,
	// // 	HttpOnly: true,
	// // })

	// ctx.Writer.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	// ctx.JSON(http.StatusOK, nil)
}

func deleteUserHandler(ctx *gin.Context) {
	idParam := ctx.Param("id")
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

	ctx.JSON(http.StatusOK, nil)
}

func listUsersHandler(ctx *gin.Context) {
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
	limit, offset := getLimitAndOffset(ctx)
	i, err := strconv.ParseInt(limit, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error parsing limit as int",
		})

		return
	}

	j, err := strconv.ParseInt(offset, 10, 32)
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

	ctx.JSON(http.StatusOK, users)
}

type updateUserRequest struct {
	ID              int64  `json:"id" binding:"required"`
	Username        string `json:"username" binding:"required"`
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}

func updateUserHandler(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// hash the current password
	currentPasswordHash, err := firstly.hasher.HashPassword(req.CurrentPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// get the current user
	user, err := firstly.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// compare the current passowrd hashed against the stored password
	if currentPasswordHash != user.Password {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("current password invalid")))
		return
	}

	// hash the new password
	newPasswordHash, err := firstly.hasher.HashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// update the password for the user with the new password hash and current id
	var updateParams db.UpdateUserParams
	updateParams.ID = user.ID
	updateParams.Password = newPasswordHash

	err = firstly.store.UpdateUser(ctx, updateParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
