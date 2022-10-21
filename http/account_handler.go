package http

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/db"
	"github.com/meads/firstly-api/security"
)

type createAccountRequest struct {
	Username string `json:"username" binding:"required"`
	Phrase   string `json:"phrase" binding:"required"`
}

func (server *FirstlyServer) CreateAccountHandler(store db.Store, hasher security.Hasher) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var req createAccountRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		tmpAccount, err := store.GetAccountByUsername(ctx, req.Username)
		if err != nil && !errors.Is(sql.ErrNoRows, err) {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if tmpAccount.ID > 0 {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("please choose another username")))
			return
		}

		var param db.CreateAccountParams
		param.Username = req.Username
		param.Salt = hasher.GenerateSalt()

		phrase, err := hasher.GeneratePasswordHash([]byte(req.Phrase), param.Salt)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("account creation failed")))
			return
		}
		param.Phrase = phrase

		account, err := store.CreateAccount(ctx, param)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// Declare the expiration time of the token
		// here, we have kept it as 5 minutes
		expirationTime := time.Now().Add(5 * time.Minute)

		// Create the JWT claims, which includes the username and expiry time
		claims := &Claims{
			Username: account.Username,
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
	}
}

type loginAccountRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *FirstlyServer) DeleteAccountHandler(store db.Store) func(*gin.Context) {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		if idParam == "" {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("id parameter is required"))
			return
		}
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("id parameter must be a valid integer"))
			return
		}

		err = store.DeleteAccount(ctx, id)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, nil)
	}
}

func (server *FirstlyServer) ListAccountsHandler(store db.Store) func(ctx *gin.Context) {
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
	return func(ctx *gin.Context) {
		limit, offset := getLimitAndOffset(ctx)
		i, err := strconv.ParseInt(limit, 10, 32)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("error parsing limit as int"))
			return
		}

		j, err := strconv.ParseInt(offset, 10, 32)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("error parsing offset as int"))
			return
		}

		images, err := store.ListAccounts(ctx, db.ListAccountsParams{Limit: int32(i), Offset: int32(j)})
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.JSON(http.StatusOK, images)
	}
}

type updateAccountRequest struct {
	ID       int64  `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Phrase   string `json:"phrase" binding:"required"`
}

func (server *FirstlyServer) UpdateAccountHandler(store db.Store, hasher security.Hasher) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		// validate the update request
		var req updateAccountRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		// verify the account exists
		accountExists, err := store.AccountExists(ctx, req.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		if !accountExists {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("error updating account")))
			return
		}

		// get the account for the request id
		account, err := store.GetAccount(ctx, req.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// update the phrase for the account using the current salt for the account
		newPhrase, err := hasher.GeneratePasswordHash([]byte(req.Phrase), account.Salt)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// initialize the parameters for the update using the known account id and new phrase
		// TODO: add more combinations to the hmac.
		var updateParams db.UpdateAccountParams
		updateParams.ID = account.ID
		updateParams.Phrase = newPhrase

		err = store.UpdateAccount(ctx, updateParams)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		account.Phrase = []byte(req.Phrase)

		ctx.JSON(http.StatusOK, struct{ username string }{username: account.Username})
	}
}
