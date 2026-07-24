package security

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// type ClaimToken struct {
// 	*jwt.Token
// }

// func (token *ClaimToken) SignedString(key interface{}) (string, error) {
// 	return token.Token.SignedString(key)
// }

type CustomClaims struct {
	Username string `json:"username"`
	*jwt.RegisteredClaims
}

// func NewUsernameClaims() *UsernameClaims {
// 	return &UsernameClaims{
// 		RegisteredClaims: &jwt.RegisteredClaims{},
// 	}
// }

// type ClaimsValidator struct {
// 	signer    jwt.SigningMethod
// 	validator jwt.Claims
// }

type ClaimsValidator struct {
}

type Claimer interface {
	// GetClaimToken() *ClaimToken
	GenerateToken(username string) (string, error)
	VerifyToken(tokenString string, secretKey []byte) (*CustomClaims, error)
	// GetFiveMinuteExpirationToken(username string) (string, time.Time, error)
	// GetFromTokenString(tokenString string) (*ClaimToken, *UsernameClaims, error)
	// ParseWithClaims(tokenString string, claims *UsernameClaims, keyFunc jwt.Keyfunc) (*ClaimToken, error)
}

// func NewClaimsValidator() Claimer {
// 	return &ClaimsValidator{
// 		signer: jwt.SigningMethodHS256,
// 	}
// }

func NewClaimsValidator() Claimer {
	return &ClaimsValidator{}
}

// func (c *ClaimsValidator) GetFromTokenString(tokenString string) (*ClaimToken, *UsernameClaims, error) {
// 	// Get the JWT string from the cookie
// 	usernameClaims := NewUsernameClaims()
// 	claimToken, err := c.ParseWithClaims(tokenString, usernameClaims, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(os.Getenv("SECRET")), nil
// 	})
// 	return claimToken, usernameClaims, err
// }

// func (c *ClaimsValidator) GetClaimToken() *ClaimToken {
// 	return &ClaimToken{
// 		&jwt.Token{
// 			Header: map[string]interface{}{
// 				"typ": "JWT",
// 				"alg": c.signer.Alg(),
// 			},
// 			Claims: c.validator,
// 			Method: c.signer,
// 		},
// 	}
// }

// func (c *ClaimsValidator) ParseWithClaims(tokenString string, claims *UsernameClaims, keyFunc jwt.Keyfunc) (*ClaimToken, error) {
// 	t, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
// 	return &ClaimToken{t}, err
// }

func (c *ClaimsValidator) VerifyToken(tokenString string, secretKey []byte) (*CustomClaims, error) {
	var claims CustomClaims

	// Parse directly into the custom claims struct
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired at %v", claims.ExpiresAt)
		}
		return nil, err
	}

	if token.Valid {
		return &claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (c *ClaimsValidator) GenerateToken(username string) (string, error) {
	// Create claims instance
	claims := &CustomClaims{
		Username: username,
		RegisteredClaims: &jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET")))
}

// func (c *ClaimsValidator) GetFiveMinuteExpirationToken(username string) (string, time.Time, error) {
// 	expirationTime := time.Now().Add(5 * time.Minute)
// 	// Create the JWT claims, which includes the username and expiry time
// 	claimsValidator := &UsernameClaims{
// 		Username: username,
// 		RegisteredClaims: &jwt.RegisteredClaims{
// 			// In JWT, the expiry time is expressed as unix milliseconds
// 			ExpiresAt: jwt.NewNumericDate(expirationTime),
// 			IssuedAt:  jwt.NewNumericDate(time.Now()),
// 		},
// 	}

// 	c.validator = claimsValidator

// 	// Declare the token with the algorithm used for signing, and the claims
// 	claimToken := c.GetClaimToken()

// 	// Create the JWT string
// 	tokenString, err := claimToken.SignedString([]byte(os.Getenv("SECRET")))
// 	return tokenString, expirationTime, err
// }
