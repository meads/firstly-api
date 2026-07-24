package http

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
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
		_, err := firstly.claimer.VerifyToken(tokenString, []byte(os.Getenv("SECRET")))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		h(ctx)
	})
}

// Create a struct that models the structure of a user, both in the request body, and in the DB
type signInRequest struct {
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
}

func signinHandler(ctx *gin.Context) {
	var req signInRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := firstly.store.GetAccountByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	valid, err := firstly.hasher.IsValidPassword(account.Password, account.Salt, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !valid {
		ctx.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString, err := firstly.claimer.GenerateToken(account.Username)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Writer.Header().Add("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	ctx.Status(http.StatusOK)
}

func refreshHandler(ctx *gin.Context) {
	// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route
	// c, err := ctx.Request.Cookie("token")
	// if err != nil {
	// 	if err == http.ErrNoCookie {
	// 		ctx.Writer.WriteHeader(http.StatusUnauthorized)
	// 		return
	// 	}
	// 	ctx.Writer.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// claimToken, usernameClaims, err := firstly.claimer.GetFromTokenString(c.Value)
	// if !claimToken.Valid {
	// 	ctx.Writer.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	// if err != nil {
	// 	// if err == jwt.ErrSignatureInvalid {
	// 	// 	ctx.Writer.WriteHeader(http.StatusUnauthorized)
	// 	// 	return
	// 	// }
	// 	ctx.Writer.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// (END) The code uptil this point is the same as the first part of the `Welcome` route

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	// if time.Until(time.Unix(usernameClaims.ExpiresAt, 0)) > 30*time.Second {
	// 	ctx.Writer.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// Now, create a new token for the current use, with a renewed expiration time
	// tokenString, expirationTime, err := firstly.claimer.GenerateToken()
	// if err != nil {
	// 	ctx.Writer.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// Set the new token as the users `session_token` cookie
	// http.SetCookie(ctx.Writer, &http.Cookie{
	// 	Name:    "session_token",
	// 	Value:   tokenString,
	// 	Expires: expirationTime,
	// })
}
