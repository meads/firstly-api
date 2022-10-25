package security

import "github.com/dgrijalva/jwt-go"

type Signage interface {
	jwt.SigningMethod
}

type ClaimValidator interface {
	jwt.Claims
}

type ClaimToken struct {
	*jwt.Token
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Keyfunc func(token *ClaimToken) (interface{}, error)

type Claimer interface {
	NewWithClaims(method Signage, claims ClaimValidator) *ClaimToken
	SignedString(key interface{}) (string, error)
	ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*ClaimToken, error)
}
