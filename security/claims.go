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
	return token.Token.SignedString(key)
}

type UsernameClaims struct {
	Username string `json:"username"`
	*jwt.StandardClaims
}

func NewUsernameClaims() *UsernameClaims {
	return &UsernameClaims{
		StandardClaims: &jwt.StandardClaims{},
	}
}

type ClaimsValidator struct {
	signer    jwt.SigningMethod
	validator jwt.Claims
}

type Claimer interface {
	GetClaimToken() *ClaimToken
	GetFiveMinuteExpirationToken(username string) (string, time.Time, error)
	GetFromTokenString(tokenString string) (*ClaimToken, *UsernameClaims, error)
	ParseWithClaims(tokenString string, claims *UsernameClaims, keyFunc jwt.Keyfunc) (*ClaimToken, error)
}

func NewClaimsValidator() Claimer {
	return &ClaimsValidator{
		signer: jwt.SigningMethodHS256,
	}
}

func (c *ClaimsValidator) GetFromTokenString(tokenString string) (*ClaimToken, *UsernameClaims, error) {
	// Get the JWT string from the cookie
	usernameClaims := NewUsernameClaims()
	claimToken, err := c.ParseWithClaims(tokenString, usernameClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	return claimToken, usernameClaims, err
}

func (c *ClaimsValidator) GetClaimToken() *ClaimToken {
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

func (c *ClaimsValidator) ParseWithClaims(tokenString string, claims *UsernameClaims, keyFunc jwt.Keyfunc) (*ClaimToken, error) {
	t, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	return &ClaimToken{t}, err
}

func (c *ClaimsValidator) GetFiveMinuteExpirationToken(username string) (string, time.Time, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claimsValidator := &UsernameClaims{
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
