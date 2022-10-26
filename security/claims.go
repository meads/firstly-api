package security

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type MethodSigner interface {
	jwt.SigningMethod
}

type ClaimValidator interface {
	jwt.Claims
}

type ClaimToken struct {
	t *jwt.Token
}

type ClaimsValidator struct {
	Username string `json:"username"`
	*jwt.StandardClaims
}

type Claimer interface {
	NewWithClaims(method MethodSigner, claims ClaimValidator) *ClaimToken
	SignedString(key interface{}) (string, error)
	ParseWithClaims(tokenString string, claims ClaimsValidator, keyFunc jwt.Keyfunc) (*ClaimToken, error)
}

func NewWithClaims(method MethodSigner, claims ClaimValidator) *ClaimToken {
	return &ClaimToken{
		t: &jwt.Token{
			Header: map[string]interface{}{
				"typ": "JWT",
				"alg": method.Alg(),
			},
			Claims: claims,
			Method: method,
		},
	}
}

func (token *ClaimToken) SignedString(key interface{}) (string, error) {
	return token.t.SignedString(key)
}

func ParseWithClaims(tokenString string, claims ClaimsValidator, keyFunc jwt.Keyfunc) (*ClaimToken, error) {
	t, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)
	return &ClaimToken{t: t}, err
}

func GetFiveMinuteExpirationToken(username string) (string, time.Time, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &ClaimsValidator{
		Username: username,
		StandardClaims: &jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	claimToken := NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := claimToken.SignedString([]byte(os.Getenv("SECRET")))
	return tokenString, expirationTime, err
}

func (validator *ClaimsValidator) Valid() error {
	return nil
}
