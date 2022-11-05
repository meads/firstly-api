package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Create a struct that models the structure of a user, both in the request body, and in the DB
type signInRequest struct {
	Phrase   string `json:"phrase" binding:"required"`
	Username string `json:"username" binding:"required"`
}

func claimsMiddleware(h gin.HandlerFunc) gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		// We can obtain the session token from the requests cookies, which come with every request
		c, err := ctx.Request.Cookie("token")
		if err != nil {
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Get the JWT string from the cookie
		claimToken, _, err := firstly.claimer.GetFromTokenString(c.Value)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.Writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if !claimToken.Valid {
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		h(ctx)
	})
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

	valid, err := firstly.hasher.IsValidPassword(account.Phrase, account.Salt, req.Phrase)
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

	tokenString, expirationTime, err := firstly.claimer.GetFiveMinuteExpirationToken(account.Username)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	ctx.Status(http.StatusOK)
}

func welcomeHandler(ctx *gin.Context) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := ctx.Request.Cookie("token")
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, usernameClaims, err := firstly.claimer.GetFromTokenString(c.Value)

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if !token.Valid {
		ctx.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Return the welcome message to the user, along with their username
	ctx.Writer.Write([]byte(fmt.Sprintf("Welcome %s!", usernameClaims.Username)))
}

func refreshHandler(ctx *gin.Context) {
	// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route
	c, err := ctx.Request.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	claimToken, usernameClaims, err := firstly.claimer.GetFromTokenString(c.Value)
	if !claimToken.Valid {
		ctx.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	// (END) The code uptil this point is the same as the first part of the `Welcome` route

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(usernameClaims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	tokenString, expirationTime, err := firstly.claimer.GetFiveMinuteExpirationToken(usernameClaims.Username)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `session_token` cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "session_token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
