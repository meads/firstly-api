package handler

import (
	"net/http"
	"time"

	"github.com/meads/firstly-api/internal/security"
)

func passClaimsMiddleware(r *http.Request, tokener *security.MockTokener) {
	tokenString := "mocktoken"
	r.Header.Add("Authorization", "Bearer mocktoken")

	userClaims := &security.UserClaims{Type: "access"}
	tokener.EXPECT().VerifyToken(tokenString).Return(userClaims, nil)
}

var defaultDate time.Time = time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)
