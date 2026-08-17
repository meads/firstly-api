package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/meads/firstly-api/internal/domain"
)

type UserClaims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Type     string `json:"type"`
	*jwt.RegisteredClaims
}

func NewUserClaims(id int64, username string, usage string, duration time.Duration) (*UserClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("Error generating token id: %w", err)
	}
	return &UserClaims{
		ID:       id,
		Username: username,
		Type:     usage,
		RegisteredClaims: &jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}

type Tokener interface {
	GenerateToken(id int64, username string, usage string, duration time.Duration) (string, *UserClaims, error)
	VerifyToken(tokenString string) (*UserClaims, error)
}

type TokenManager struct {
	secretKey string
}

func NewTokenManager(secretKey string) Tokener {
	return &TokenManager{
		secretKey: secretKey,
	}
}

func (c *TokenManager) GenerateToken(id int64, username string, usage string, duration time.Duration) (string, *UserClaims, error) {
	claims, err := NewUserClaims(id, username, usage, duration)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(c.secretKey))
	if err != nil {
		return "", nil, fmt.Errorf("error signing token: %w", err)
	}

	return tokenStr, claims, nil
}

func (c *TokenManager) VerifyToken(tokenStr string) (*UserClaims, error) {
	// Parse directly into the custom claims struct
	claims := &UserClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			// invalid token signing method
			return nil, domain.ErrInvalidToken
		}
		return []byte(c.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrInvalidToken
		}
		return nil, domain.ErrInvalidToken
	}

	if token.Valid {
		return claims, nil
	}

	return nil, domain.ErrInvalidToken
}
