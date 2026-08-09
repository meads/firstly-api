package domain

import (
	"context"
	"time"
)

type Session struct {
	ID           string    // `json:"id"`
	UserID       int64     // `json:"userId"`
	RefreshToken string    // `json:"refreshToken"`
	IsRevoked    bool      // `json:"isRevoked"`
	ExpiresAt    time.Time // `json:"expiresAt"`
	CreatedAt    time.Time // `json:"createdAt"`
}

type CreateSessionParams struct {
	ID           string    // `json:"id"`
	UserID       int64     // `json:"userId"`
	RefreshToken string    // `json:"refreshToken"`
	IsRevoked    bool      // `json:"isRevoked"`
	ExpiresAt    time.Time // `json:"expiresAt"`
}

type SessionRepository interface {
	CreateSession(ctx context.Context, arg CreateSessionParams) (*Session, error)
	GetSession(ctx context.Context, id string) (*Session, error)
	DeleteSession(ctx context.Context, id string) error
	RevokeSession(ctx context.Context, id string) error
	// RevokeUserSessions(ctx context.Context, userID int64) error
}
