package http

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/db"
)

type createAccountRequest struct {
	Username string `json:"username" binding:"required"`
	Phrase   string `json:"phrase" binding:"required"`
}

func createAccountHandler(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tmpAccount, err := firstly.store.GetAccountByUsername(ctx, req.Username)
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
	param.Salt = firstly.hasher.GenerateSalt()

	phrase, err := firstly.hasher.GeneratePasswordHash([]byte(req.Phrase), param.Salt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("account creation failed")))
		return
	}
	param.Phrase = phrase

	account, err := firstly.store.CreateAccount(ctx, param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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
}

type loginAccountRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func deleteAccountHandler(ctx *gin.Context) {
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

	err = firstly.store.DeleteAccount(ctx, id)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func listAccountsHandler(ctx *gin.Context) {
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
		ctx.AbortWithError(http.StatusBadRequest, errors.New("error parsing limit as int"))
		return
	}

	j, err := strconv.ParseInt(offset, 10, 32)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("error parsing offset as int"))
		return
	}

	images, err := firstly.store.ListAccounts(ctx, db.ListAccountsParams{Limit: int32(i), Offset: int32(j)})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.JSON(http.StatusOK, images)
}

type updateAccountRequest struct {
	ID       int64  `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Phrase   string `json:"phrase" binding:"required"`
}

func updateAccountHandler(ctx *gin.Context) {
	// validate the update request
	var req updateAccountRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// verify the account exists
	accountExists, err := firstly.store.AccountExists(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if !accountExists {
		ctx.JSON(http.StatusNotFound, errorResponse(errors.New("error updating account")))
		return
	}

	// get the account for the request id
	account, err := firstly.store.GetAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// update the phrase for the account using the current salt for the account
	newPhrase, err := firstly.hasher.GeneratePasswordHash([]byte(req.Phrase), account.Salt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// initialize the parameters for the update using the known account id and new phrase
	// TODO: add more combinations to the hmac.
	var updateParams db.UpdateAccountParams
	updateParams.ID = account.ID
	updateParams.Phrase = newPhrase

	err = firstly.store.UpdateAccount(ctx, updateParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	account.Phrase = []byte(req.Phrase)

	ctx.JSON(http.StatusOK, struct{ username string }{username: account.Username})
}
