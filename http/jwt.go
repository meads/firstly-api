package http

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/db"
	"github.com/meads/firstly-api/security"
)

var jwtKey = []byte(os.Getenv("SECRET"))

// Create a struct that models the structure of a user, both in the request body, and in the DB
type signInRequest struct {
	Phrase   string `json:"phrase" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (server *FirstlyServer) SigninHandler(store db.Store, hasher security.Hasher) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var req signInRequest

		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		account, err := store.GetAccountByUsername(ctx, req.Username)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		valid, err := hasher.IsValidPassword(account.Phrase, account.Salt, req.Phrase)
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

		// Declare the expiration time of the token
		// here, we have kept it as 5 minutes
		expirationTime := time.Now().Add(5 * time.Minute)

		// Create the JWT claims, which includes the username and expiry time
		claims := &Claims{
			Username: req.Username,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}

		// Declare the token with the algorithm used for signing, and the claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Create the JWT string
		tokenString, err := token.SignedString(jwtKey)
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
}

func (server *FirstlyServer) WelcomeHandler(store db.Store) func(*gin.Context) {
	return func(ctx *gin.Context) {
		// We can obtain the session token from the requests cookies, which come with every request
		c, err := ctx.Request.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				ctx.Writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			return
		}

		// Get the JWT string from the cookie
		tknStr := c.Value

		// Initialize a new instance of `Claims`
		claims := &Claims{}

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.Writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Finally, return the welcome message to the user, along with their
		// username given in the token
		ctx.Writer.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
	}
}

func (server *FirstlyServer) RefreshHandler(store db.Store) func(*gin.Context) {
	return func(ctx *gin.Context) {
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
		tknStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if !tkn.Valid {
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
		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			return
		}

		// Now, create a new token for the current use, with a renewed expiration time
		expirationTime := time.Now().Add(5 * time.Minute)
		claims.ExpiresAt = expirationTime.Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
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
}
