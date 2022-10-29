package security

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type ClaimToken struct {
	*jwt.Token
}

func (token *ClaimToken) SignedString(key interface{}) (string, error) {
	return token.SignedString(key)
}

type ClaimsValidator struct {
	Username string `json:"username"`
	*jwt.StandardClaims
}

type Claims struct {
	token     ClaimToken
	signer    jwt.SigningMethod
	validator jwt.Claims
}

type Claimer interface {
	// NewWithClaims(method MethodSigner, claims ClaimValidator) *ClaimToken
	GetFiveMinuteExpirationToken(username string) (string, time.Time, error)
	GetClaimToken() *ClaimToken
	ParseWithClaims(tokenString string, claims ClaimsValidator, keyFunc jwt.Keyfunc) (*ClaimToken, error)
}

func NewClaims() Claimer {
	return &Claims{
		signer: jwt.SigningMethodHS256,
	}
}

func (c *Claims) GetClaimToken() *ClaimToken {
	return &ClaimToken{
		&jwt.Token{
			Header: map[string]interface{}{
				"typ": "JWT",
				"alg": c.signer.Alg(),
			},
			Claims: c.validator,
			Method: c.signer,
		},
	}
}

func (c *Claims) ParseWithClaims(tokenString string, claims ClaimsValidator, keyFunc jwt.Keyfunc) (*ClaimToken, error) {
	t, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)
	return &ClaimToken{t}, err
}

func (c *Claims) GetFiveMinuteExpirationToken(username string) (string, time.Time, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claimsValidator := &ClaimsValidator{
		Username: username,
		StandardClaims: &jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	c.validator = claimsValidator

	// Declare the token with the algorithm used for signing, and the claims
	claimToken := c.GetClaimToken()

	// Create the JWT string
	tokenString, err := claimToken.SignedString([]byte(os.Getenv("SECRET")))
	return tokenString, expirationTime, err
}
